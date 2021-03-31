
# go-chan-over-chan

Examples to understand the Channels-Over-Channels concept offered by Golang

## summary

There are uses for this channel-over-channel strategy, but the ack one is simple and powerful. 
Further, in many cases when you need to “return” something to another goroutine, sending it a `chan` on which it can return 
a value is often the easiest way to do it. This pattern can even be useful when you want to wait for a goroutine to ack its completion. 
Note, however, that you can also do ack-ing with a `sync.WaitGroup`.

Below are some code examples using the 3 strategies. In each case, We’ll simulate the work using a simple `time.Sleep`.

## run

### Style 1: Using a Channel Inside a Channel

Here's the simplest of the patterns in action. Generally this style will be easiest to read and understand, but it has some limits:

- Each `doStuff` goroutine sleeps for a set amount of time. You can't change the sleep time when you send on `ch`
- Each `doStuff` goroutine can only receive a `chan time.Duration` – no more data than that

```bash
go run ./chan-in-chan/main.go
```

### Style 2: Using a Channel Stored Inside a Struct

This code will look almost identical to the previous snippet, with 2 exceptions:

- The ack channel will be stored inside a `struct`
- The sleep time will be stored inside that same `struct`, so we can pass it over the `channel`
    - This makes the code more flexible, because we can tell `doStuff` how long to sleep when we send to it, rather than when we start it

```bash
go run ./chan-in-struct/main.go
```

### Style 3: Using a Channel Inside a Function Closure

This code will look different from the previous examples, because the `doStuff` function won’t know anything about a return channel. 
That fact is both good and bad. On the up side, you can change your code later to do anything you want inside that function 
(e.g. good for testing!), but on the down side, you can’t pass dynamic `time.Durations` into the `doStuff` goroutines, as you could in the previous example.

```bash
go run ./chan-in-func-closure/main.go
```

## links

- https://www.goin5minutes.com/blog/channel_over_channel/
