package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/guigateixeira/general-auth/application"
)

func main() {
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)

	if err != nil {
		fmt.Println("Failed to start app:", err)
	}
}
