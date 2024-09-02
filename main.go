package main

import (
	"fmt"
	"os"
)

func main() {
	app, err := NewGridApp(30, 15) // Initial grid size
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing app: %v\n", err)
		os.Exit(1)
	}
	defer app.Cleanup()

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
