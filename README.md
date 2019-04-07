# checkd

`checkd` is a Go library for writing programs which collect metrics. For example, suppose you want to monitor the number of processes running on each server in your infrastructure, or the number of pods running in your Kubernetes cluster. You can write a very small and simple program using `checkd`, compile it to a single binary, and run it. The program will run an OpenMetrics listener which exposes your custom metric, so that you can scrape it with Prometheus, or any other OpenMetrics client, such as a Datadog agent.

## Why?

Why is the `checkd` library useful? After all, you can already use tools like `node_exporter` to expose OpenMetrics data, so you can collect, graph, and alert on practically any piece of information about your infrastructure. Kubernetes also exposes lots of useful metrics. And in your own software, you can use the Prometheus client library to expose your key metrics such as request latency and internal statistics. So where does `checkd` come in?

The answer is that `checkd` is designed to make it easy to build programs which exist only to gather and expose metrics. Monitoring tools, in other words.

With traditional monitoring systems such as Nagios and Icinga, you can write your own custom plugins to monitor whatever you want, so long as they conform to a standard output format that the monitoring system can recognise and parse.

In a metrics-based world, `checkd` allows you to do the same kind of thing, creating custom monitoring checks that can do anything you want, evaluate the results the way you want, and expose exactly the metric that you want to monitor.

## Example

Let's look at a simple example:

```go
package main

import "github.com/bitfield/checkd"

func init() {
	checkd.Every(10*time.Second, func() {
		checkd.Gauge("unix_time_seconds", "Current Unix time").SetToCurrentTime()
	})
}

func main() {
    checkd.Start()
}
```

This is literally all the code you need to write to create a customised monitoring daemon. This compiles to a single binary which you can run on your servers, or in your Kubernetes cluster. The daemon runs forever, updating its metric every 10 seconds, and listening on port 8666 for metric scrapes.

It's up to you whether you want to write a single program that does all the custom checks you want for your infrastructure, or many programs, each of which runs in different places checking different things. `checkd` makes it easy either way.

## Check functions

The core of a `checkd` program is the _check function_. This is simply a Go function which takes no parameters and returns no results. It can do absolutely anything at all. To be useful as a monitoring tool, a check function will usually create and update a Prometheus metric, but that's not compulsory.

## Metrics

If you've used the Prometheus client library to instrument an application, you know that there's a certain amount of boilerplate required: you need to set the metric options, register the metric, save the metric somewhere, and update it when necessary.

`checkd` eliminates most of that boilerplate, by giving you a simple function to call to create a metric of the type you want. For example, a counter:

```go
checkd.Counter("check_calls_total", "Number of times the checker has been called").Inc()
```

`Counter` creates and registers a Prometheus counter metric with the specified name and help text. It returns the metric, so that you can update it (for example, by calling `Inc()` on it). But it also stores the metric for you, so that if you call `Counter` again with the same name, you just get the existing metric again, rather than creating a new one.

This condenses a three-stage process (register the metric, save the metric, update the metric) into just one stage. Thus, you can easily write one-line check functions like that in the example.

If you need a gauge metric, use `Gauge` instead of `Counter`.

## Registering check functions

Having created a check function which monitors the metric you're interested in, how do you tell `checkd` about it? You can use the `Every` function to schedule your check function to be called at a specified time interval:

```go
checkd.Every(time.Minute, checkStuff)
```

This is usually done in the `init` function, so that check registration happens at startup, but it's not essential. However, all your checks must be registered before you tell `checkd` to actually start running the checks.

## Starting the checks

To start the checks, call `checkd.Start()`. This is usually done in the `main` function:

```go
func main() {
    checkd.Start()
}
```

`Start` creates a concurrent goroutine for each registered check. Inside the check goroutine, first of all, the check function is called (so all checks run at least once at startup). The goroutine then sleeps for the configured interval before calling the check function again.

`Start` also starts a Prometheus metrics listener exposing any metrics registered by your check functions. The default listening port is 8666, but you can change this by setting the `checkd.Port` value:

```go
checkd.Port = 9999
checkd.Start()
```

## Scraping your metrics

To collect the data exposed by your custom checks, add the listening port as a scrape target in Prometheus or whatever metrics client you're using:

```yaml
scrape_configs:
  - job_name: "my_custom_checks"
    scrape_interval: "10s"
    static_configs:
      - targets:
        - localhost:8666
```

## A more detailed example

Look in the [example](./example) directory to see a complete `checkd` program, which runs two independent checks. You can use this as a starting point for your own `checkd`-based monitoring tools.