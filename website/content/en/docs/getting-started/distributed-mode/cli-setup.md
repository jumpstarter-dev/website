---
title: Service CLI Setup
weight: 10
description: >
  This section explains how to setup the Jumpstarter CLI to work with the Jumpstarter distributed service.
---

The distributed service CLI (jmpctl) is a command line interface that enables easy administration
of your Exporters and Clients.

Please refer to https://github.com/jumpstarter-dev/jumpstarter-controller/releases/latest for
the latest release.

## Installing the CLI
{{< highlight bash  >}}
  export VERSION=0.1.0
  export ARCH=amd64

  curl -L https://github.com/jumpstarter-dev/jumpstarter-controller/releases/download/${VERSION}/jmpctl_${VERSION}_linux_${ARCH} -o /usr/local/bin/jmpctl
  chmod a+x /usr/local/bin/jmpctl
{{< / highlight >}}

## Configuration
The `jmpctl` CLI requires a kubeconfig file with permissions to access the jumpstarter installation
namespace `jumpstarter-lab` is the default if you followed the previous install section.


## Usage
{{< highlight bash >}}
$ jmpctl
Admin CLI for managing jumpstarter

Usage:
  jmpctl [command]

Available Commands:
  client      Manage clients
  completion  Generate the autocompletion script for the specified shell
  exporter    Manage exporters
  help        Help about any command

Flags:
  -h, --help                help for jmpctl
      --kubeconfig string   Path to the kubeconfig file to use
      --namespace string    Kubernetes namespace to operate on (default "default")
      --timeout string      command timeout (default "10s")

Use "jmpctl [command] --help" for more information about a command.

{{< / highlight >}}

