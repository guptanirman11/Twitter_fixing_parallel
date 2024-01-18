
# Project 2: A Simple Twitter Client/Server System

This project is intended to the use and implementation
of parallel data structures using low-level primitives.

## Single-User Twitter Feed

In this Parallel project I used the following Go
concurrent constructs:

  - `go` statement
  - `sync.Mutex` and its associated methods.
  - `sync/atomic` package. You may use any of the atomic operations.
  - `sync.WaitGroup` and its associated methods.
  - `sync.Cond` and its associated methods.


## Part 1: Twitter Feed

Here in this project I tried re-developing the data structure that represents a user's
feed. My implementation redefines it as a singly linked list.

I implemented few methods of a `feed`
(i.e. the `Add`, `Remove`, and `Contains` methods). I used the
internal representations for `type feed struct` and `type post struct`
in your implementation. *

Tested the implementation of feed by using the test file called
`feed_test.go`.

  - `TestSimpleSeq`
  - `TestAdd`
  - `TestContains`
  - `TestRemove`



Sample run of the `SimpleSeq` and `TestAdd` tests:

    //Run top-level tests matching "SimpleSeq" such as "TestSimpleSeq" and "TestAdd".
    $ go test -v -run "SimpleSeq|TestAdd"
    === RUN   TestSimpleSeq
    --- PASS: TestSimpleSeq (0.00s)
    === RUN   TestAdd
    --- PASS: TestAdd (0.00s)
    PASS
    ok    hw4/feed    0.078s


## Part 2: Thread Safety using a Read-Write Lock

A read/write lock mechanism allows multiple readers to access a data
structure concurrently, but only a single writer is allowed to access
the data structures at a time. Implemented a read/write lock library that
**only** uses a **single condition variable** and **mutex** for its
synchronization mechanisms. Go provides a Read/Write lock that is
implemented using atomics:

  - [R/W Mutex in Go](https://golang.org/pkg/sync/#RWMutex)

As with the Go implementation, I implemeted the following four methods
associated with my lock:

  - `Lock()`
  - `Unlock()`
  - `RLock()`
  - `RUnlock()`

I limited the max number of readers to `32`.

### Coarse Grain Feed

Now, I made the feed library inside of `feed.go`
thread-safe by using my implementation of a read/write lock. I needed
to think about the appropriate places to call the various read/write
locking and unlocking methods in the `feed` methods. 


## Part 3: A Twitter Server

Using the `feed` library from Part 1, I implemented a server that
processes **requests**, which perform actions (e.g., add a post, remove
a post, etc.) on a single feed. These requests come from a client
program (e.g., the twitter mobile app) where a user can request to
modify their feed. The server sends **responses** back to the client
with a result (e.g., notification that a post was added and/or removed
successfully, etc.).

For the client and server I agreed on the JSON format to send the requests and
responses. 


The client and server use decoders and encoders to parse
the JSON string into a type easily usable in Go. 
`json.Decoder`, which acts as a streaming buffer of requests, and
`json.Encoder`, which acts as a streaming buffer of responses. This
model is a simplified version of a real-life client-server model used
heavily in many domains such as web development. At a high-level, you
can think of your program as a simple “server” in the client-server
model illustrated below:

Requests (i.e., tasks in our program) are sent from a “client” (e.g., a
redirected file on the command line, a task generator program piped into
 program, etc.) via os.Stdin. The “server” (i.e., program) will
process these requests and send their results back to the client via
os.Stdout. This model is a mimicking a real-life client-server model
used heavily in many domains such as web development; however, we are
not actually implementing web client-server system in this assignment.

### Requests and Responses Format

The basic format for the requests coming in from `json.Decoder` will be
of the following format:

``` json
{ 
"command": string, 
"id": integer, 
... data key-value pairings ... 
}
```

A request will always have a `"command"` and `"id"` key. The `"command"`
key holds a string value that represents the type of feed task. The
`"id"` represents an unique identification number for this request.
Requests are processed asynchronously by the server so requests can be
processed out of order from how they are received from `json.Decoder`;
therefore, the `"id"` acts as a way to tell the client that result
coming back from the server is a response to an original request with
this specific `"id"` value. Thus, **it is not your responsibility to
maintain this order and you must not do anything to maintain it in your
program**.

The remaining key-value pairings represent the data for a specific
request. The following subsections will go over the various types of
requests.

### Add Request

An add request adds a new post to the feed data structure. The
`"command"` value will always be the string `"ADD"`. The data fields
include a key-value pairing for the message body (`"body": string`) and
timestamp (`"timestamp": number`). For example,

``` json
{ 
  "command": "ADD", 
  "id": 342, 
  "body": "just setting up my twttr", 
  "timestamp": 43242423
}
```

After completing a `"ADD"` request, the goroutine assigned the request
will send a response back to the client via `json.Encoder` acknowledging
the add was successful. The response is a JSON object that includes a
success key-value pair (`"success": boolean`). For an add request, the
value is always true since you can add an infinite number of posts with the
original identification number.

``` json
{ 
  "success": true, 
  "id": 342
}
```

### Remove Request

A remove request removes a post from the feed data structure. The
`"command"` value will always be the string `"REMOVE"`. The data fields
include a key-value pairing for the timestamp (`"timestamp": number`)
that represents the post that should be removed. For example,

``` json
{ 
  "command": "REMOVE", 
  "id": 2361, 
  "timestamp": 43242423
}
```

After completing a `"REMOVE"` task, the goroutine assigned the task will
send a response back to the client via `json.Encoder` acknowledging the
remove was successful or unsuccessful. The response is a JSON object
that includes a success key-value pair (`"success": boolean`). For a
remove request, the value is `true` if the post with the requested
timestamp was removed, otherwise assign the key to `false`with the original
identification number.

``` json
{ 
  "success": true, 
  "id": 2361
}
```

### Contains Request

A contains request checks to see if a feed post is inside the feed data
structure. The `"command"` value will always be the string `"CONTAINS"`.
The data fields include a key-value pairing for the timestamp
(`"timestamp": number`) that represents the post to check. For example,

``` json
{ 
  "command": "CONTAINS", 
  "id": 2362, 
  "timestamp": 43242423
}
```

After completing a `"CONTAINS"` task, the goroutine assigned the task
will send a response back to the client via `json.Encoder` acknowledging
whether the feed contains that post. The response is a JSON object that
includes a success key-value pair (`"success": boolean`). For a contains
request, the value is `true` if the post with the requested timestamp is
inside the feed, otherwise assign the key to `false`with the original
identification number 

``` json
{ 
  "success": false, 
  "id": 2362
}
```

**Note**: Assuming we removed the post previously.

### Feed Request

A feed request returns all the posts within the feed. The `"command"`
value will always be the string `"FEED"`. Their are no data fields for
this request. For example,

``` json
{ 
  "command": "FEED", 
  "id": 2, 
}
```

After completing a `"FEED"` task, the goroutine assigned the task will
send a response back to the client via `json.Encoder` with all the posts
currently in the feed. The response is a JSON object that includes a
success key-value pair (`"feed": [objects]`). For a feed request, the
value is a JSON array that includes a JSON object for each feed post.
Each JSON object will include a `"body"` key (`"body": string`) that
represents a post’s body and a `"timestamp"` key (`"timestamp": number`)
that represents the timestamp for the post with the original identification
number.

``` json
{ 
  "id": 2
  "feed":[
        { 
          "body": "This is my second twitter post", 
          "timestamp": 43242423
        },
        {
          "body": "This is my first twitter post", 
          "timestamp": 43242420
        }
        ]
}
```

### Done Request

If client will no longer send requests then it sends a done request. The
`"command"` value will always be the string `"DONE"`. Their are no data
fields for this request. For example,

``` json
{ 
  "command": "DONE" 
}
```

This notifies server it needs to *shutdown* (i.e., close down the
program). A done request signals to the main goroutine that no further
processing is necessary after this request is received. No response is
sent back to the client.

### Implementing the Server


When a goroutine calls the `Run` function, it will start the server
based on the configuration passed to the function. This function does not return (i.e., it is a blocking function)
until the server is shutdown (i.e., it receives the `"DONE"` request).

### Parallel Version: Tasks Queue

Inside the `Run` function, if `config.Mode == "p"` then the server will
run the parallel version. This version is implemented using a *task
queue*. This task queue is another work distribution technique and my
first exposure to the producer-consumer model. In this model, the
producer will be the main goroutine and its job is to collect a series
of tasks and place them in a queue structure to be executed by consumers
(also known as workers). The consumers will be spawned goroutines. I implemented the parallelization as follows:

1.  The main goroutine begins by spawning a specified
    `config.ConsumersCount` goroutines, where each will begin executing
    a function called `func consumer(...)`. Each goroutine will either
    begin doing work or go to sleep in a conditional wait if there is no
    work to begin processing yet. 
2.  After spawning the consumer goroutines, the main goroutine will call
    a function `func producer(...)`. Inside the producer function, the
    main goroutine reads in from `json.Decoder` a series of tasks (i.e.,
    requests). For the sake of explicitness, the tasks will be feed
    operations for a single user-feed that the program manages. The main
    goroutine will place the tasks inside of a queue data structure and
    do the following:
      - If there is a consumer goroutine waiting for work then place a
        task inside the queue and wake one consumer up.
      - Otherwise, the main gorountine continues to place tasks into the
        queue. Eventually, the consumers will grab the tasks from the
        queue at later point in time.
3.  Inside the `func consumer(...)` function each consumer goroutine
    will try to grab one task from the queue. The consumer will then
    process the request and send the appropriate response back. When a
    consumer finishes executing its task, it checks the queue to grab
    another task. If there are no tasks in the queue then it will need
    to wait for more tasks to process or exit its function if there are
    no remaining tasks to complete.

### Additional Queue Requirements

I implemented this queue data structure so that both the main and
worker goroutines have access to retrieve and modify it. All work is
placed in this queue so workers can grab tasks when necessary. The actual enqueuing and
dequeuing of items is done in a unbounded `lock-free manner`
(i.e., non-blocking).


### Sequential Version

I wroter a sequential version of this program where the
main goroutine processes and executes all the requested tasks without
spawning any gorountines.

## Part 4: The Twitter Client

Inside the `twitter/twitter.go`, I defined a simple Twitter client
that has the following usage and required command-line argument:

``` 
Usage: twitter <number of consumers>
    <number of consumers> = the number of goroutines (i.e., consumers) to be part of the parallel version.  
```

The program needs to create a `server.Config` value based off the above
arguments. If `<number of consumers>` is not entered then this means you
need to run your sequential version of the server; otherwise, run the
parallel version. The `json.Decoder` and `json.Encoder` should be
created by using `os.Stdin` and `os.Stdout`. Once the `Run` function returns than the
program exits.

Assumptions: No error checking is needed. All tasks read-in will be in
the correct format with all its specified data. All command line
arguments will be given as specified.

## Part 5: Benchmarking Performance

In Part 5, I tested the execution time of parallel
implementation by averaging the elapsed time of the twitter tests. I ran the timings on a CS cluster.

## Part 6: Performance Measurement

Inside the `proj2/benchmark` directory, you will see the a file called
`benchmark.go`. This program copies over the all requests test cases you
saw from `twitter/twitter_test.go`(i.e., extra-small, small, medium,
large, and extra-large). The benchmark program allows  to execute one
of these test cases using sequential or parallel versions and
outputs the elapsed time for executing that test. 

### Generation of speedup graphs

I used the `benchmark.go` program to produce a speedup graph for
the different test-cases by varying the number of threads. The set of
threads will be `{2,4,6,8,12}`. 

    
    > \[Speedup = \frac{\text{wall-clock time of serial execution}}{\text{wall-clock time of parallel execution}}\]
    


## Part 7: Performance Analysis

Made a report (pdf document, text file, etc.) summarizing 
results from the experiments and the conclusions you draw from them.

  - A brief description of the project (i.e., an explanation what is
    implemented in feed.go, server.go, twitter.go).

  - Instructions on how to run your testing script. We should be able to
    just run your script. However, if we need to do another step then
    please let us know in the report.

  -   - explaination of the results of 
        graph.
    
        - Evaluate the impact of the linked-list implementation on performance.
        - Explore potential enhancements in synchronization techniques, queuing, and producer/consumer components.
        - Assess hardware effects on benchmark performance.

