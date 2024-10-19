---
title: Distributed Mode
description: Getting started with the Jumpstarter distributed mode
#categories: [Examples, Placeholders]
tags: [test, docs]
weight: 3
---

Jumpstarter's distributed mode is useful for labs, where teams need to collaborate
and share devices, for development and CI/CD testing.

The distributed service is provided from Kubernetes and provides:

* A registry of:
  * Exporters
  * Clients: users, and robot-users.

* The service to let clients Lease exporters.

* The routing between clients and exporters.

