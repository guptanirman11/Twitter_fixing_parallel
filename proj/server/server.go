package server

import (
	"encoding/json"
	"fmt"
	"proj2/feed"
	"proj2/queue"
	"sync"
	"time"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

type Response struct {
	Success bool
	ID      int
	Feed    []feed.FeedItem `json:"feed"`
}

// Run starts up the twitter server based on the configuration
// information provided and only returns when the server is fully
// shutdown.

func Run(config Config) {
	if config.Mode == "p" {
		runParallel(config)
	} else if config.Mode == "s" {
		runSequential(config)
	}
}

func runParallel(config Config) {
	server := feed.NewFeed()
	taskQueue := queue.NewLockFreeQueue()
	var wg sync.WaitGroup
	done := false

	// Spawn consumer goroutines
	for i := 0; i < config.ConsumersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer(taskQueue, server, config.Encoder, &done)
		}()
	}

	go producer(config.Decoder, taskQueue)

	for !done {
		time.Sleep(time.Millisecond)
	}

	wg.Wait() // Wait for all consumers to finish
}

func producer(decoder *json.Decoder, taskQueue *queue.LockFreeQueue) {
	for {
		var request queue.Request
		if err := decoder.Decode(&request); err != nil {

			continue
		}

		taskQueue.Enqueue(&request)
	}
}

func consumer(taskQueue *queue.LockFreeQueue, server feed.Feed, encoder *json.Encoder, done *bool) {
	for {
		request := (taskQueue.Dequeue())

		if request == nil {
			if *done {
				return // Exit when done flag is set
			}

			time.Sleep(time.Millisecond)
			continue
		}

		switch request.Command {
		case "ADD":
			AddPost(*request, server, encoder)
		case "REMOVE":
			RemovePost(*request, server, encoder)
		case "CONTAINS":
			ContainsPost(*request, server, encoder)
		case "FEED":
			GetFeed(*request, server, encoder)
		case "DONE":
			*done = true
			return
		}
	}
}

func runSequential(config Config) {
	server := feed.NewFeed()

	for {
		var request queue.Request
		if err := config.Decoder.Decode(&request); err != nil {
			fmt.Println(err)
			continue
		}

		switch request.Command {
		case "ADD":
			AddPost(request, server, config.Encoder)
		case "REMOVE":
			RemovePost(request, server, config.Encoder)
		case "CONTAINS":
			ContainsPost(request, server, config.Encoder)
		case "FEED":
			GetFeed(request, server, config.Encoder)
		case "DONE":
			return
		}
	}
}

func AddPost(request queue.Request, server feed.Feed, encoder *json.Encoder) {

	server.Add(request.Body, request.Timestamp)

	response := Response{
		Success: true,
		ID:      request.ID,
	}

	// Encoding and sending the response.
	EncodeAndSendResponse(response, encoder)
}

func RemovePost(request queue.Request, server feed.Feed, encoder *json.Encoder) {

	success := server.Remove(request.Timestamp)

	response := Response{
		Success: success,
		ID:      request.ID,
	}
	EncodeAndSendResponse(response, encoder)
}

func ContainsPost(request queue.Request, server feed.Feed, encoder *json.Encoder) {

	success := server.Contains(request.Timestamp)

	response := Response{
		Success: success,
		ID:      request.ID,
	}
	EncodeAndSendResponse(response, encoder)
}

func GetFeed(request queue.Request, server feed.Feed, encoder *json.Encoder) {

	feedData := server.GetFeedData()

	response := Response{
		Success: true,
		ID:      request.ID,
		Feed:    feedData,
	}

	// Encoding and sending the response.
	EncodeAndSendResponse(response, encoder)
}

func EncodeAndSendResponse(response Response, encoder *json.Encoder) {
	if err := encoder.Encode(response); err != nil {
		fmt.Println("error in encoding and sending response")
		fmt.Println(err)
	}
}
