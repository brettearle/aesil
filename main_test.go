package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brettearle/aesil/internal/app"
	"github.com/brettearle/aesil/internal/wait"
)

func TestEndPoints(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	go app.Run(ctx, app.NewConfig("127.0.0.1", "8989"), os.Stdout, os.Stderr)
	wait.ForServerReady(ctx, time.Second, "http://127.0.0.1:8989/ping")

	t.Run("/ping", func(t *testing.T) {
		t.Parallel()
		want := "Pong"
		client := http.Client{}
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8989/ping", nil)
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
