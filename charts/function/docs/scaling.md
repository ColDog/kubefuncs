# Scaling

Currently a horizontal pod autoscaler is deployed for every function to scale it up. The metric used is CPU utilization.

Future plans are to integrate custom metrics with the NSQ server to scale based on NSQ queue size, then functions should ideally spin up and down quickly depending on the size of the queue they are processing.
