package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/schwartzmx/gremtune"
	"github.com/schwartzmx/gremtune/interfaces"
)

func main() {

	host := "localhost"
	port := 8182
	hostURL := fmt.Sprintf("ws://%s:%d/gremlin", host, port)
	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}).With().Timestamp().Logger()

	cosmos, err := gremtune.New(hostURL, gremtune.WithLogger(logger), gremtune.NumMaxActiveConnections(10), gremtune.ConnectionIdleTimeout(time.Second*1))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create the cosmos connector")
	}

	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM)
	exit_chan := make(chan struct{})
	go func() {
		for {

			ticker := time.NewTicker(time.Millisecond * 10)
			defer ticker.Stop()

			select {
			case <-signal_chan:
				close(exit_chan)
				return
			case <-ticker.C:
				queryCosmos(cosmos, logger)
			}
		}
	}()

	<-exit_chan
	if err := cosmos.Stop(); err != nil {
		logger.Error().Err(err).Msg("Failed to stop")
	}
	logger.Info().Msg("Teared down")
}

func queryCosmos(cosmos *gremtune.Cosmos, logger zerolog.Logger) {
	res, err := cosmos.Execute("g.addV('Phil')")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return
	}

	for i, chunk := range res {
		jsonEncodedResponse, err := json.Marshal(chunk.Result.Data)

		if err != nil {
			logger.Error().Err(err).Msg("Failed to encode the raw json into json")
			continue
		}
		logger.Info().Str("reqID", chunk.RequestID).Int("chunk", i).Msgf("Received data: %s", jsonEncodedResponse)
	}
}

func queryCosmosAsync(cosmos *gremtune.Cosmos, logger zerolog.Logger) {
	dataChannel := make(chan interfaces.AsyncResponse)

	go func() {
		for chunk := range dataChannel {
			jsonEncodedResponse, err := json.Marshal(chunk.Response.Result.Data)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to encode the raw json into json")
				continue
			}
			logger.Info().Str("reqID", chunk.Response.RequestID).Msgf("Received data: %s", jsonEncodedResponse)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	err := cosmos.ExecuteAsync("g.addV('Phil')", dataChannel)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute async a gremlin command")
		return
	}
}
