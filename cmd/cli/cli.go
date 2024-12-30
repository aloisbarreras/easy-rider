package main

import (
	"context"
	"easyrider"
	"easyrider/netlify"
	"easyrider/sheets"
	"easyrider/vercel"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
)

// Default settings
const (
	DefaultSource      = "google-sheets"
	DefaultFormat      = "netlify"
	DefaultDestination = "stdout"
)

// main sets up signal handlers and context and delegates to Run to execute the command.
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

	cfg := &easyrider.CLIConfig{}
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	if err := parseFlags(fs, cfg, args); err != nil {
		return err
	}

	switch cmd {
	case "generate":
		loader, err := getLoader(cfg)
		if err != nil {
			return err
		}

		formatter, err := getFormatter(cfg)
		if err != nil {
			return err
		}

		out, err := getOutWriter(cfg)
		if err != nil {
			return err
		}

		return (&easyrider.GenerateCommand{
			Cfg:       cfg,
			Loader:    loader,
			Formatter: formatter,
			Out:       out,
		}).Run(ctx)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func getLoader(cfg *easyrider.CLIConfig) (easyrider.RedirectLoader, error) {
	var source easyrider.RedirectLoader
	switch easyrider.RedirectSource(cfg.Source) {
	case easyrider.RedirectSourceGoogleSheets:
		s := sheets.NewRedirectSource(cfg.SheetID)
		s.ServiceAccountCredentialsFile = cfg.ServiceAccountCredentialsFile
		s.SheetRange = cfg.SheetRange
		source = s
	default:
		return nil, fmt.Errorf("unknown source %q", cfg.Source)
	}

	return source, nil
}

func getFormatter(cfg *easyrider.CLIConfig) (easyrider.RedirectFormatter, error) {
	var formatter easyrider.RedirectFormatter
	switch easyrider.RedirectFormat(cfg.Format) {
	case easyrider.RedirectFormatNetlify:
		formatter = &netlify.RedirectFormatter{}
	case easyrider.RedirectFormatVercel:
		formatter = &vercel.RedirectFormatter{}
	default:
		return nil, fmt.Errorf("unknown format %q", cfg.Format)
	}

	return formatter, nil
}

func getOutWriter(cfg *easyrider.CLIConfig) (io.Writer, error) {
	var out io.Writer
	switch easyrider.RedirectDestination(cfg.Destination) {
	case easyrider.RedirectDestinationStdout:
		out = os.Stdout
	default:
		return nil, fmt.Errorf("unknown destination %q", cfg.Destination)
	}

	return out, nil
}

// parseFlags parses the command line flags and sets the relevant fields in the CLIConfig.
func parseFlags(fs *flag.FlagSet, cfg *easyrider.CLIConfig, args []string) error {
	src := fs.String("source", DefaultSource, "source of the data")
	format := fs.String("format", DefaultFormat, "format of the data")
	dest := fs.String("destination", DefaultDestination, "destination of the data")
	svcAccount := fs.String("service-account", "", "path to the service account credentials file")
	sheetID := fs.String("google-sheet-id", "", "ID of the Google Sheet")
	sheetRange := fs.String("google-sheet-range", "A2:C", "range of the Google Sheet to read")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg.Source = *src
	cfg.Format = *format
	cfg.Destination = *dest
	cfg.ServiceAccountCredentialsFile = *svcAccount
	cfg.SheetID = *sheetID
	cfg.SheetRange = *sheetRange

	return nil
}
