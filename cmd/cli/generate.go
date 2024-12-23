package main

import (
	"context"
	"easyrider"
	"flag"
	"io"
	"os"
)

type GenerateCommand struct{}

func (g *GenerateCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	src := fs.String("source", DefaultSource, "source of the data")
	format := fs.String("format", DefaultFormat, "format of the data")
	dest := fs.String("destination", DefaultDestination, "destination of the data")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var source easyrider.RedirectSource
	switch *src {
	case "google-sheets":
		source = &easyrider.GoogleSheetsSource{}
	default:
		return flag.ErrHelp
	}

	var formatter easyrider.RedirectFormatter
	switch *format {
	case "netlify":
		formatter = &easyrider.NetlifyFormatter{}
	case "vercel":
		formatter = &easyrider.VercelFormatter{}
	default:
		return flag.ErrHelp
	}

	var out io.Writer
	switch *dest {
	case "stdout":
		out = os.Stdout
	default:
		return flag.ErrHelp
	}

	redirects, err := source.LoadRedirects()
	if err != nil {
		return err
	}

	formatted, err := formatter.FormatRedirects(redirects)
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
}
