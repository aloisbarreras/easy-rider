package easyrider

import (
	"context"
	"io"
)

// CLIConfig holds the configuration options for the CLI.
type CLIConfig struct {
	Source      string
	Format      string
	Destination string

	// Google Sheets settings
	ServiceAccountCredentialsFile string
	SheetID                       string
	SheetRange                    string
}

type GenerateCommand struct {
	Cfg       *CLIConfig
	Loader    RedirectLoader
	Formatter RedirectFormatter
	Out       io.Writer
}

func (g *GenerateCommand) Run(ctx context.Context) error {
	redirects, err := g.Loader.LoadRedirects(ctx)
	if err != nil {
		return err
	}

	formatted, err := g.Formatter.FormatRedirects(redirects)
	if err != nil {
		return err
	}

	_, err = g.Out.Write(formatted)
	return err
}
