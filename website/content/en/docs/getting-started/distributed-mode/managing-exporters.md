---
title: Managing Exporters
weight: 23
description: >
  This section explains how to use the `jmpctl` CLI to create exporters and manage them.
---

## Creating an exporter

{{< highlight bash >}}
jmpctl exporter create my-exporter --namespace jumpstarter-lab > my-exporter.yaml
{{< / highlight >}}

This `my-exporter.yaml` should be installed in the box