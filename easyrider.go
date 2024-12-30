package easyrider

import (
	"context"
)

type RedirectSource string
type RedirectFormat string
type RedirectDestination string

// Default settings
const (
	RedirectSourceGoogleSheets RedirectSource      = "google-sheets"
	RedirectFormatNetlify      RedirectFormat      = "netlify"
	RedirectFormatVercel       RedirectFormat      = "vercel"
	RedirectDestinationStdout  RedirectDestination = "stdout"
)

type Redirect struct {
	From       string
	To         string
	StatusCode int
}

type RedirectLoader interface {
	LoadRedirects(ctx context.Context) ([]Redirect, error)
}

type RedirectFormatter interface {
	FormatRedirects([]Redirect) ([]byte, error)
}
