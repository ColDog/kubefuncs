# Developing Functions

Developing functions is exactly the same as developing with a docker based application. For local development, the following docker-compose file can bring up the dependencies necessary to run a local function:

```yaml
# docker-compose.yaml
version: '2'

# Brings up the following services on the host network with ports:
# nsq: 4150
# lookupd: 4161
# gateway: 8080
services:
  nsqlookup:
    image: nsqio/nsq
    network_mode: host
    command: /nsqlookupd

  nsq:
    image: nsqio/nsq
    network_mode: host
    command: /nsqd --broadcast-address 127.0.0.1 --lookupd-tcp-address=127.0.0.1:4160

  gateway:
    image: coldog/kubefuncs-gateway:latest
    network_mode: host
```
