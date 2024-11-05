---
title: Service Installation
weight: 1
description: >
  This section contains the installation instructions for the Jumpstarter distributed service.
---

When building a distributed environment with Jumpstarter, you will need to install the Jumpstarter distributed service. This service is responsible for managing the devices, and the communication between the devices and and clients.

You will need an OpenShift or Kubernetes deployment, and the right kubeconfig file with admin
 credentials (at least the first install may need assistance from your cluster administrator
 for the purpose of installing the service CRDs).

Alternatively you can setup a local Kubernetes cluster with [kind](https://kind.sigs.k8s.io/) (kubernetes in docker) following the instructions below.

{{% alert title="Note" color="info" %}}
The direct helm install will auto-generate random router and controller secrets, but if you use ArgoCD make sure to set these values to unique values.
{{% /alert %}}


{{< tabpane text=true right=false >}}
    {{% tab header="**Methods**:" disabled=true /%}}
    {{< tab header="**OpenShift**" >}}

      Please note that the global.baseDomain is used to create the host names for the services,
      with the provided example the services will be available at grpc.jumpstarter.example.com
      and router.jumpstarter.example.com.
      <br/><br/>
      To install using helm:

      <br/>
      {{< highlight bash  >}}
  helm upgrade jumpstarter --install oci://quay.io/jumpstarter-dev/helm/jumpstarter \
              --create-namespace --namespace jumpstarter-lab \
              --set global.baseDomain=jumpstarter.example.com \
              --set global.metrics.enabled=true \
              --set jumpstarter-controller.grpc.mode=route \
              --version=0.1.0
      {{< / highlight >}}
    {{< /tab >}}

    {{< tab header="**Kubernetes**" >}}
        {{< highlight bash  >}}
helm upgrade jumpstarter --install oci://quay.io/jumpstarter-dev/helm/jumpstarter \
            --create-namespace --namespace jumpstarter-lab \
            --set global.baseDomain=devel.jumpstarter.dev \
            --set global.metrics.enabled=true # disable if metrics not available \
            --set jumpstarter-controller.grpc.mode=ingress \
            --version=0.1.0
        {{< / highlight >}}
    {{< /tab >}}

    {{< tab header="**Kind**" >}}
    Kind is a tool for running local Kubernetes clusters using Podman or Docker container “nodes”.
    <br/><br/>
    Begin by figuring out the LAN ip address that it's accessible for your docker/podman host, and do:
    {{< highlight bash >}}
export IP="LAN accessible address to your docker/podman instance"
    {{< / highlight >}}
    <br/><br/>
    Then you can continue with:
    {{< highlight bash  >}}

cat <<EOF > kind_config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  kind: ClusterConfiguration
  apiServer:
    extraArgs:
      "service-node-port-range": "3000-32767"
- |
  kind: InitConfiguration
  nodeRegistration:
    kubeletExtraArgs:
      node-labels: "ingress-ready=true"
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80 # ingress controller
    hostPort: 5080
    protocol: TCP
  - containerPort: 30010 # grpc nodeport
    hostPort: 8082
    protocol: TCP
  - containerPort: 30011 # grpc router nodeport
    hostPort: 8083
    protocol: TCP
  - containerPort: 443 # minimalistic UI
    hostPort: 5443
    protocol: TCP
EOF

export BASEDOMAIN="jumpstarter.${IP}.nip.io"
export GRPC_ENDPOINT="grpc.${BASEDOMAIN}:8082"
export GRPC_ROUTER_ENDPOINT="router.${BASEDOMAIN}:8083"

kind create cluster  --config kind_config.yaml

helm upgrade jumpstarter --install oci://quay.io/jumpstarter-dev/helm/jumpstarter \
            --create-namespace --namespace jumpstarter-lab \
            --set global.baseDomain=${BASEDOMAIN} \
            --set jumpstarter-controller.grpc.endpoint=${GRPC_ENDPOINT} \
            --set jumpstarter-controller.grpc.routerEndpoint=${GRPC_ROUTER_ENDPOINT} \
            --set global.metrics.enabled=false \
            --set jumpstarter-controller.grpc.nodeport.enabled=true \
            --set jumpstarter-controller.grpc.mode=nodeport \
            --version=0.1.0
        {{< / highlight >}}
    {{< /tab >}}

    {{< tab header="**ArgoCD in OpenShift**" >}}
        <h3>Create namespace</h3>
        First, we must create a namespace for the Jumpstarter installation. This namespace
        should be labeled with argocd.argoproj.io/managed-by=<your-argo-CD-instance>to allow
        ArgoCD to manage the resources in the namespace.</br>

        In this case, using the default openshift-gitops ArgoCD deployment, the command would be:
        {{< highlight bash  >}}
        kubectl create namespace jumpstarter-lab
        kubectl label namespace jumpstarter-lab argocd.argoproj.io/managed-by=openshift-gitops
        {{< / highlight >}}
        </br>
        <h3>Application</h3>

         <b>jumpstarter-controller.controllerSecret</b> and <b>jumpstarter-controller.routerSecret</b>
         are secrets that are used to secure the authentication between clients and the jumpstarter elements.
         <b>These secrets should be unique and not shared between installations</b>. Helm installation can
         auto-generate values for these, but with ArgoCD such mechanism doesn't work. You need to manually
         create these secrets in the namespace where the jumpstarter is installed.

{{< highlight yaml >}}
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: jumpstarter
  namespace: openshift-gitops
spec:
  destination:
    name: in-cluster
    namespace: jumpstarter-lab
  project: default
  source:
    chart: jumpstarter
    helm:
      parameters:
      - name: global.baseDomain
        value: devel.jumpstarter.dev
      - name: global.metrics.enabled
        value: "true"
      - name: jumpstarter-controller.controllerSecret
        value: "pick-a-secret-DONT-USE-THIS-DEFAULT"
      - name: jumpstarter-controller.routerSecret
        value: "again-pick-a-secret-DONT-USE-THIS-DEFAULT"
      - name: jumpstarter-controller.grpc.mode
        value: "route"
    repoURL: quay.io/jumpstarter-dev/helm
    targetRevision: "0.1.0"
{{< / highlight >}}

<h3>Note: CRDs</h3>
ArgoCD needs to be able to manage the CRDs that Jumpstarter uses. This is done by creating a ClusterRole and ClusterRoleBinding that allows the ArgoCD application controller to manage the CRDs.

An alternative to this is to manually create and update the CRDs
that jumpstarter uses.

{{< highlight yaml  >}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    argocds.argoproj.io/name: openshift-gitops
    argocds.argoproj.io/namespace: openshift-gitops
  name: openshift-gitops-argocd-appcontroller-crd
rules:
- apiGroups:
  - 'apiextensions.k8s.io'
  resources:
  - 'customresourcedefinitions'
  verbs:
  - '*'
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    argocds.argoproj.io/name: openshift-gitops
    argocds.argoproj.io/namespace: openshift-gitops
  name: openshift-gitops-argocd-appcontroller-crd
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: openshift-gitops-argocd-appcontroller-crd
subjects:
- kind: ServiceAccount
  name: openshift-gitops-argocd-application-controller
  namespace: openshift-gitops
{{< / highlight >}}

    {{< /tab >}}
{{< /tabpane >}}

## Where can I find help?

Look at the [Community](/community/) page for more information on how to get in touch with the Jumpstarter developers,
we have a Matrix channel where we'll be happy to help you with your installation.