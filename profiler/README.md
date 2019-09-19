# Profiler samples

This directory contains [Profiler](https://cloud.google.com/profiler/docs) samples.
These samples are configured to collect profiling data and transmit that data to your
Google Cloud Platform project. You can use the
[Profiler UI](https://cloud.google.com/profiler/docs/using-profiler)
to analyze the collected profiles. For the Go language, your analysis options include
CPU time, heap, allocated heap, contention, and thread analysis. You can even
compare sets of profiles.

For information on configuring your applications to collect and transmit profiling data to
your Google Cloud Platform project, the the following guides:

+ [Go](https://cloud.google.com/profiler/docs/profiling-go)
+ [Java](https://cloud.google.com/profiler/docs/profiling-java)
+ [Node.js](https://cloud.google.com/profiler/docs/profiling-nodejs)
+ [Python](https://cloud.google.com/profiler/docs/profiling-python)

## Using these samples

You can execute these samples, without change, on Google Cloud Platform.
To execute these samples outside of Google Cloud Platform, you must perform additional setup.
For information, see
[Profiling applications runing outside of Google Cloud Platform](https://cloud.google.com/profiler/docs/profiling-external).

For detailed information on using these samples, see
[Profiling samples](https://cloud.google.com/profiler/docs/samples).

## Samples

### profiler_quickstart

The sample `profiler_quickstart` is configured to run the `hello-profiler` service.
This is a very simple service that is used by the
[Quickstart](https://cloud.google.com/profiler/docs/quickstart) guide.

[Go Code](/profiler/profiler_quickstart)

### hotapp

The sample `hotapp` is uses an infinite loop with two call stacks.
The default configuration of this sample is used as a benchmarking tool,
to verify that as dependent services change the overall performance is static.
The [Profiler documentation](https://cloud.google.com/profiler/docs)
includes images generated from this sample. The
`docdemo-service` has a specific configuration that adds work.

[Go Code](/profiler/hotapp)

### hotmid

Sample `hotmid` is an application that simulates multiple calls to a library
function made via different call paths. Each of these calls is not
particularly expensive (and so does not stand out on the flame graph). But
in the aggregate these calls add up to a significant time which can be
identified via looking at the flat list of functions' self and total time.

[Go Code](/profiler/hotmid)

