---
title: Installation
weight: 1
description: >
  This section contains the installation instructions for the Jumpstarter distributed service, and the python framework.
---

## Installing the distributed service

In the documentation website you can find the installation instructions for the Jumpstarter distributed service
in the various supported platforms.

| Platform | Install method | Instructions |
|----------|----------------|--------------|
| Kubernetes  | Helm |[Kubernetes Installation Guide](https://docs.jumpstarter.dev/installation/service/kubernetes-helm.html) |
| OpenShift | Helm | [OpenShift with Helm Installation Guide](https://docs.jumpstarter.dev/installation/service/openshift-helm.html) |
| OpenShift | Helm + ArgoCD | [OpenShift with ArgoCD Installation Guide](https://docs.jumpstarter.dev/installation/service/openshift-argocd.html) |
| Local cluster with Kind | Helm | [Local cluster with Kind Installation Guide](https://docs.jumpstarter.dev/installation/service/kind-helm.html) |
| Local cluster with Minikube | Helm |[Local cluster with Minikube Installation Guide](https://docs.jumpstarter.dev/installation/service/minikube-helm.html) |

## Installing the CLI tools and Exporters

Jumpstarter provides two cli tools:
  * `jmpctl`: a command line tool to manage clients and exporters in the Jumpstarter distributed service.
  * `jmp`: a command line tool for interacting with Jumpstarter as a client: leasing devices, performing
     operations on them, etc.

In addition there is a exporter tool:
  * `jmp-exporter`: This command line tool is used to run exporters in the Jumpstarter, either distributed
  or local workflow.

Installation details for all those tools can be found here: https://docs.jumpstarter.dev/installation/index.html

## Where can I find help?

Look at the [Community](/community/) page for more information on how to get in touch with the Jumpstarter developers,
we have a Matrix channel where we'll be happy to help you.