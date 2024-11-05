---
title: Managing Exporters
weight: 20
description: >
  This section explains how to use the `jmpctl` CLI to create exporters and manage them.
---

## Creating an exporter

Creating an exporter via `jmpctl` means registering a new exporter with the Jumpstarter service, and getting a base configuration file that can be installed in the exporter
host machine.

{{< highlight bash >}}
jmpctl exporter create my-exporter --namespace jumpstarter-lab > my-exporter.yaml
{{< / highlight >}}

The `my-exporter.yaml` should be configured with the desired exported drivers filling up the
export section, see https://docs.jumpstarter.dev/config.html#exporter-config for more details.

Details about existing drivers can be found in:
https://docs.jumpstarter.dev/reference/contrib/index.html

If you don't have the hardware ready yet but you want to try things out you
can setup the exporter with something like the following example which
will provide a few mock interfaces to play with:

{{< highlight yaml >}}
apiVersion: jumpstarter.dev/v1alpha1
kind: ExporterConfig
endpoint: grpc.jumpstarter.192.168.1.10.nip.io:8082
token: << token data >>
export:
    storage:
        type: jumpstarter.drivers.storage.driver.MockStorageMux
    power:
        type: jumpstarter.drivers.power.driver.MockPower
    echonet:
        type: jumpstarter.drivers.network.driver.EchoNetwork
    can:
        type: jumpstarter_driver_can.driver.Can
        config:
            channel: 1
            interface: "virtual"
{{< / highlight >}}

Once the exporter configuration is ready it should be installed in the
exporter host machine at
`/etc/jumpstarter/exporters/my-exporter.yaml`.

{{% alert title="Note" color="info" %}}
Remember, the exporter is Linux service that exports the interfaces to the target DUT(s)
(serial ports, video interfaces, bluetooth, anything that Jumpstarter has a driver for,
and the exporter service can reach via linux device or network). In this case the exporter
service calls back to the Jumpstarter service to report the available interfaces and
wait for commands.
{{% /alert %}}

### Running the exporter

A simple way to run the exporter and observe behavior before installing it more permanently is to run:

{{< highlight bash >}}
sudo podman run --rm -ti --name my-exporter --net=host  --privileged \
                -e JUMPSTARTER_GRPC_INSECURE=1 \
                -v /run/udev:/run/udev -v /dev:/dev -v /etc/jumpstarter:/etc/jumpstarter \
                quay.io/jumpstarter-dev/jumpstarter:main \
                jmp exporter run my-exporter

INFO:jumpstarter.exporter.exporter:Registering exporter with controller
INFO:jumpstarter.exporter.exporter:Currently not leased

<CTRL+C to stop>
{{< / highlight >}}

This will run the exporter in the foreground and allow you to observe the behavior.
Please note that we are running the exporter as root, with access to all devices and network
interfaces, this is most likely necessary for the exporter to access the interface
devices that will allow it to control the target DUT(s).

Once we are satisfied with the behavior we can install the exporter as a systemd service:

Create a systemd service file at `/etc/containers/systemd/my-exporter.container` with
the following content:

{{< highlight ini >}}
[Unit]
Description=My exporter

[Container]
ContainerName=my-exporter
Environment=JUMPSTARTER_GRPC_INSECURE=1
Exec=jmp exporter run my-exporter
Image=quay.io/jumpstarter-dev/jumpstarter:main
Network=host
PodmanArgs=--privileged
Volume=/run/udev:/run/udev
Volume=/dev:/dev
Volume=/etc/jumpstarter:/etc/jumpstarter

[Install]
WantedBy=multi-user.target default.target
{{< / highlight >}}

Then enable and start the service:

{{< highlight bash >}}
sudo systemctl enable --now my-exporter
{{</ highlight >}}

### Observing the exporter

Once the exporter is running, it will be available for clients to lease and use, also
can be observed via kubectl:

{{< highlight bash >}}
kubectl get exporter my-exporter -n jumpstarter-lab -o yaml
{{< / highlight >}}

{{< highlight yaml >}}
apiVersion: jumpstarter.dev/v1alpha1
kind: Exporter
metadata:
  creationTimestamp: "2024-10-19T18:26:35Z"
  generation: 1
  name: my-exporter
  namespace: jumpstarter-lab
  resourceVersion: "63609"
  uid: 9bfa4f35-7f56-4615-9dd7-fdf3fbd0650e
spec: {}
status:
  conditions:
  - lastTransitionTime: "2024-10-19T18:36:36Z"
    message: ""
    observedGeneration: 1
    reason: Register
    status: "True"
    type: Registered
  - lastTransitionTime: "2024-10-19T18:31:52Z"
    message: ""
    observedGeneration: 1
    reason: Connect
    status: "True"
    type: Online
  credential:
    name: my-exporter-exporter
  devices:
  - labels:
      jumpstarter.dev/client: jumpstarter.drivers.composite.client.CompositeClient
    uuid: e7d828ac-929a-471b-a2d4-5fe44b83e498
  - labels:
      jumpstarter.dev/client: jumpstarter.drivers.storage.client.StorageMuxClient
      jumpstarter.dev/name: storage
    parent_uuid: e7d828ac-929a-471b-a2d4-5fe44b83e498
    uuid: cc8905dc-fd89-4cf8-a2e2-f4978073cd56
  - labels:
      jumpstarter.dev/client: jumpstarter.drivers.power.client.PowerClient
      jumpstarter.dev/name: power
    parent_uuid: e7d828ac-929a-471b-a2d4-5fe44b83e498
    uuid: 05a5b94b-be77-4d61-aace-622ca68b413c
  - labels:
      jumpstarter.dev/client: jumpstarter.drivers.network.client.NetworkClient
      jumpstarter.dev/name: echonet
    parent_uuid: e7d828ac-929a-471b-a2d4-5fe44b83e498
    uuid: 811dada0-e276-4652-89fa-554e59112f7a
  - labels:
      jumpstarter.dev/client: jumpstarter_driver_can.client.CanClient
      jumpstarter.dev/name: can
    parent_uuid: e7d828ac-929a-471b-a2d4-5fe44b83e498
    uuid: 336144b3-cf57-47ae-a630-0fe84f3773b6
  endpoint: grpc.jumpstarter.192.168.1.10.nip.io:8082
  {{< / highlight >}}

## Deleting an exporter
Once a exporter is decomisioned it can be deleted with:

{{< highlight bash >}}
jmpctl exporter delete my-exporter --namespace jumpstarter-lab
{{< / highlight >}}
or via kubectl:

{{< highlight bash >}}
kubectl delete exporter my-exporter -n jumpstarter-lab
{{< / highlight >}}
