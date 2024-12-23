package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

// Default settings
const (
	DefaultSource      = "google-sheets"
	DefaultFormat      = "netlify"
	DefaultDestination = "stdout"
)

type Config struct {
	Source      string
	Format      string
	Destination string
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	return Config{
		Source:      DefaultSource,
		Format:      DefaultFormat,
		Destination: DefaultDestination,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	if err := Run(ctx, os.Args[1:]); errors.Is(err, flag.ErrHelp) {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, args []string) error {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	switch cmd {
	case "generate":
		return (&GenerateCommand{}).Run(ctx, args)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}
