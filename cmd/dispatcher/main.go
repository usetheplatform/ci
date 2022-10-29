package main

import (
	"time"

	. "github.com/usetheplatform/ci-system/pkg/common"
	"github.com/usetheplatform/ci-system/pkg/dispatcher"
)

// @TODO: Be able to start connections / close connections with runners
// and update dispatcher on that
func main() {
	options, err := dispatcher.GetOptions()
	CheckIfError(err)

	server := dispatcher.NewServer(options)
	dispatcher := dispatcher.NewDispatcher(&dispatcher.DispatcherOptions{
		CheckRunnersTimeout: time.Millisecond * 5,
		RedistributeTimeout: time.Millisecond * 5,
	})

	// @TODO: At this moment there are no runners, design how to add them

	Info("Starting Dispatcher::CheckRunners")
	dispatcher.CheckRunners()

	Info("Starting Dispatcher::Redistribute")
	dispatcher.Redistribute()

	Info("Starting Server::Serve")
	server.Serve()

	// 1. Open ws server and listen to messages from the observer
	// 2. Start runner checker (goroutine)
	// 3. Start redistributor (goroutine)
	// 4. Listen to messages from runners to get notified, when runner is free (goroutine)
	// 5. Save data to fs as backup

}
