# Monitoring

Currently basic prometheus metrics are exported by the functions. More work is desired here to integrate tracing as well as a robust set of prometheus metrics.

However, if you have custom instrumentation, you can wrap the client library with a new client that provides custom instrumentation very easily.
