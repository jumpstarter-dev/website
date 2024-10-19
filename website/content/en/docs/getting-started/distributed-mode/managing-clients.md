---
title: Managing Clients
weight: 20
description: >
  This section explains how to use the `jmpctl` CLI to create clients and manage them.
---


## Creating a client
{{< highlight bash >}}
jmpctl client create john --namespace jumpstarter-lab > john.yaml

# allow the client to be able to use all drivers
cat >> john.yaml <<EOF
drivers:
  allow: []
  unsafe: True
EOF
{{< / highlight >}}

This `john.yaml` should be shared with the user (John in this case) that will be
interacting with the Jumpstarter service as a consumer (developer, tester, etc.)
leasing exporters and developing/running tests. This file should be installed
to `~/.config/jumpstarter/clients/john.yaml`.

A client config can also be created and configured as a secret for a robot
CI/CD accounts that need to perform automated testing.

For additional details on the client configuration, please refer to the
[Client config](https://docs.jumpstarter.dev/config.html#clients) documentation on the docs
site.

On the kubernetes API, the new  client is represented as a `client.jumpstarter.dev` custom resource:
```bash
kubectl get client john -n jumpstarter-lab -o yaml
```

{{< highlight yaml >}}
apiVersion: jumpstarter.dev/v1alpha1
kind: Client
metadata:
  creationTimestamp: "2024-10-19T12:02:32Z"
  generation: 1
  name: john
  namespace: jumpstarter-lab
  resourceVersion: "20400"
  uid: 4ea788ec-d964-4b1c-9e80-61b15e808426
spec: {}
status:
  credential:
    name: john-client
  endpoint: grpc.jumpstarter.192.168.1.10.nip.io:8082
{{< / highlight >}}

## Listing clients

Listing can be performed through the CLI or the Kubernetes API.
{{< highlight bash >}}
jmpctl client list --namespace jumpstarter-lab
NAME   AGE
john   3m27s
{{< / highlight >}}
or
{{< highlight bash >}}
kubectl get client -n jumpstarter-lab
NAME   AGE
john   3m27s
{{< / highlight >}}

## Deleting a client
When a client is no longer needed, or we want to revoke access to the Jumpstarter service, we can
delete it.

{{< highlight bash >}}
jmpctl client delete john --namespace jumpstarter-lab
{{< / highlight >}}
or
{{< highlight bash >}}
kubectl delete client john -n jumpstarter-lab
{{< / highlight >}}
