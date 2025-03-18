package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brettearle/aesil/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, app.NewConfig("127.0.0.1", "8989"), os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
