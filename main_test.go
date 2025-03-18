package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brettearle/aesil/internal/app"
)

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func TestEndPoints(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	go app.Run(ctx, app.NewConfig("127.0.0.1", "6969"), os.Stdout, os.Stderr)
	waitForReady(ctx, time.Second, "http://127.0.0.1:6969/ping")

	t.Run("/ping", func(t *testing.T) {
		want := "Pong"
		client := http.Client{}
		req, _ := http.NewRequest("GET", "http://127.0.0.1:6969/ping", nil)
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("got error %v", err)
		}
		defer res.Body.Close()
		gotBytes, err := io.ReadAll(res.Body)
		got := string(gotBytes)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
