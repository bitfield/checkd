This example shows how to use the `checkd` library to build a very simple metrics-gathering daemon which exposes two metrics:

* Current Unix time
* Number of times the check has been called

While you can build a `checkd` daemon with just one `main.go` file, this example has one Go file per check, plus a `main.go` to start all the checks running. It's probably easier to organise your code this way if you want to build one daemon that checks multiple things in parallel.

First, let's look at the `counter.go` file. This defines a `checkCounter` function:

```go
func checkCounter() {
	checkd.Counter("check_calls_total", "Number of times the checker has been called").Inc()
}
```

How does it work? The `checkd.Counter` function creates a Prometheus counter metric with the given name and help text, and returns the metric for you to use. But if the metric has already been registered, it simply returns the existing metric instead. This means you don't have to keep track of the metric yourself between calls to the check function; `checkd` takes care of that for you.

Since this check function just records the number of times it has been called, it calls `Inc()` on the counter metric.

We also need some way to tell `checkd` to call this check function, and how often, so we do that in the `init()` function:

```go
func init() {
	checkd.Every(time.Minute, checkCounter)
}
```

You can read this out loud as "Every minute, checkCounter". `checkd.Every` registers the specified check function (in this case, `checkCounter`), and schedules it to be called at the specified interval (in this case, every minute).

For a slightly different example, let's look at the time checker in `time.go`:

```go
func init() {
	checkd.Every(10*time.Second, func() {
		checkd.Gauge("unix_time_seconds", "Current Unix time").SetToCurrentTime()
	})
}
```

Instead of defining a named check function and passing it to `Every`, this just passes it an anonymous function literal. It's up to you which style you prefer.

`checkd.Gauge` behaves just like `checkd.Counter`, registering, caching, and returning the specified metric.

With all our checks registered and scheduled, how do we actually start them running? Here's `main.go`:

```go
func main() {
	checkd.Start()
}
```

Very simple! `checkd.Start` starts each registered check running in a separate goroutine, sleeping for the specified time in between each call to the check function. It also starts a Prometheus metrics listener running on the default port (8666). `Start` runs forever, or until it encounters a fatal error, which will be logged before exit.
