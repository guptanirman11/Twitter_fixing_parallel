package main

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/server"
	"strconv"
)

func main() {
	var mode string
	var numConsumers int

	if len(os.Args) < 2 {
		mode = "s"
	} else {
		numConsumersStr := os.Args[1]
		parsedNumConsumers, err := strconv.Atoi(numConsumersStr)
		if err != nil {
			fmt.Println("Invalid number of consumers:", err)
			os.Exit(1)
		}

		if parsedNumConsumers > 0 {
			mode = "p"
			numConsumers = parsedNumConsumers
		}

	}

	decoder := json.NewDecoder(os.Stdin)

	config := server.Config{
		Encoder:        json.NewEncoder(os.Stdout),
		Decoder:        decoder,
		Mode:           mode,
		ConsumersCount: numConsumers,
	}

	// Run the server with the specified configuration.
	server.Run(config)
}
