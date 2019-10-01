# Profiler samples

This directory contains [Profiler](https://cloud.google.com/profiler/docs) samples.
These samples are configured to collect profiling data and transmit that data to your
Google Cloud Platform project. You can use the
[Profiler UI](https://cloud.google.com/profiler/docs/using-profiler)
to analyze the collected profiles. For the Go language, your analysis options include
CPU time, heap, allocated heap, contention, and thread analysis. You can even
compare sets of profiles.

For profiling concepts, the capabilities of the Profiler UI, and how to configure your application
to collect and transmit profiling data, see the
[Profiler documentation](https://cloud.google.com/profiler/docs).

See the following configuration guides for details on modifying your application to collect
and transmit profiling data:

+ [Profiling Go applications](https://cloud.google.com/profiler/docs/profiling-go)
+ [Profiling Java applications](https://cloud.google.com/profiler/docs/profiling-java)
+ [Profiling Node.js applications](https://cloud.google.com/profiler/docs/profiling-nodejs)
+ [Profiling Python applications](https://cloud.google.com/profiler/docs/profiling-python)

See
[Profiling applications runing outside of Google Cloud Platform](https://cloud.google.com/profiler/docs/profiling-external) 
for additional configuration steps that are required when you are running your service outside of
Google Cloud Platform.

## Samples

Detailed instructions on executing these samples from Cloud Shell is included in
[Profiling samples](https://cloud.google.com/profiler/docs/samples).

### profiler_quickstart

The sample `profiler_quickstart` is configured to run the `hello-profiler` service.
This is a very simple service that is used by the
[Quickstart](https://cloud.google.com/profiler/docs/quickstart) guide.

[Go Code](/profiler/profiler_quickstart)

### hotapp

The sample `hotapp` is uses an infinite loop with two call stacks.

The [Profiler documentation](https://cloud.google.com/profiler/docs)
includes images generated from this sample. 
If you wish to generate profile data consistent with that included in the Profiler documentation,
run the `hotapp` service with the following command line options:
```
go run main.go -service=docdemo-service -local_work -skew=75 -version=1.75.0
```

[Go Code](/profiler/hotapp)

### hotmid

Sample `hotmid` is an application that simulates multiple calls to a library
function made via different call paths. Each of these calls is not
particularly expensive (and so does not stand out on the flame graph). But
in the aggregate these calls add up to a significant time which can be
identified via looking at the flat list of functions' self and total time.

[Go Code](/profiler/hotmid)

## Executing a sample

To execute a sample and collect profiling data in your GCP project, do the following:

1.  If you have a new GCP, you need to enable the Profiler API for your project. Choose one of the following methods.

    From the Cloud console, go to **APIs & Services** and then click **Enable APIS and Services**.
    Search for **Profiler**.  If the API isn't enabled, click **Enable**.

    From Cloud Shell, run the following command:
 
    ```
    gcloud services enable cloudprofiler.googleapis.com
    ``` 

2.  If you aren't running on GCP, then you need to create a service account. For details on these steps, see 
    [Profiling applications runing outside of Google Cloud Platform](https://cloud.google.com/profiler/docs/profiling-external).
    
3.  From your clone of the GitHub repository, change to the source directory of the program you want to execute.
    For example, the following command changes the working directory to that for the sample `hotapp`:
 
    ```
    cd ~/gopath/src/github.com/GoogleCloudPlatform/golang-samples/profiler/hotapp
    ```

4. If you aren't running on GCP, edit `main.go` and specify your GCP project ID.

5. Start the program:

   ```
   go run main.go
   ```

A few seconds after you start the program, the message `profiler has started` is displayed.
New messages are displayed each time a profile is uploaded to your GCP project.
To stop the program, enter `Ctrl-C`.

### Viewing your data

To view your profile data, do the following:

1. Go to the [Cloud Console](https://console.cloud.google.com).
1. From the Navigation menu, scroll to the **Stackdriver** section and then select **Profiler**. 

Each time you click **Now**, the Profiler UI is refreshed and includes profiles up to the current point in time.

