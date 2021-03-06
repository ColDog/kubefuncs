# Clients

Client libraries are the core concept behind kubefuncs. This describes the specification that all client libraries should implement. Documentation is generated for each client.

## Releases

Each client is versioned separately.

## Specification

### Environment Variables

The following environment variables should be read by default:

- `TOPIC`: The nsq topic that this client should listen on.
- `CHANNEL`: The nsq channel that this client should listen on.
- `NSQ_LOOKUPD_ADDR`, `NSQ_LOOKUPD_PORT`: Configures the lookupd address.
- `NSQ_NSQD_ADDR`, `NSQ_NSQD_PORT`: Configures the nsqd address.
- `HEALTHZ_ADDR`: Healthz port to listen on.

### Events

Clients listen to an nsq topic, pushing a message into the queue is the only way to call a function.

- The message body is the protocol buffer message type event, protubuf messages are in the [proto](proto) package.
- When instantiating a new event a UUID must be assigned as the event id.
- Events have a payload field which uses the protbuf any type to support messaging types.

The message handling should implement this behaviour:

```python
def handle_event(msg, handler):
  try:
    event = unmarshal(msg.body)
  except MarshalError: # on MarshalError give up
    finish(msg).
    return

  error = None
  try:
    response = handler(event)
  except Exception as e:
    error = e

  if event.Return:
    # If there is a return, let's send it something.
    if error:
      response = proto.Error(error)
    if respons is None:
      response = proto.Empty()
    finish(msg)
    send(event.Return, response)
  else:
    # Requeue the message up to the configured max attempts.
    if error:
      retry(msg)
```
