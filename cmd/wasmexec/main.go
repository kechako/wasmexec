package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/kechako/wasmexec/cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := &cli.App{}
	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		code := 1
		if coder, ok := err.(interface{ Code() int }); ok {
			code = coder.Code()
		}
		os.Exit(code)
	}
}
