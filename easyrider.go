package easyrider

import (
	"encoding/json"
	"strings"
)

type Redirect struct {
	From       string
	To         string
	StatusCode string
}

type RedirectSource interface {
	LoadRedirects() ([]Redirect, error)
}

type GoogleSheetsSource struct {
	SheetID string
	APIKey  string
}

func (s *GoogleSheetsSource) LoadRedirects() ([]Redirect, error) {
	// This is a stub implementation that returns some redirects
	return []Redirect{
		{
			From:       "/from.html",
			To:         "/to.html",
			StatusCode: "200",
		},
		{
			From:       "/from2.html",
			To:         "/to2.html",
			StatusCode: "301",
		},
	}, nil
}

type RedirectFormatter interface {
	FormatRedirects([]Redirect) ([]byte, error)
}

type NetlifyFormatter struct{}

func (n *NetlifyFormatter) FormatRedirects(redirects []Redirect) ([]byte, error) {
	//Netlify _redirects format looks like:
	// /from.html	/to.html		200

	var formattedRedirects []string
	for _, r := range redirects {
		formattedRedirects = append(formattedRedirects, r.From+"	"+r.To+"		"+r.StatusCode)
	}

	return []byte(strings.Join(formattedRedirects, "\n") + "\n"), nil
}

type VercelRedirects struct {
	Redirects []struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		StatusCode  int    `json:"statusCode"`
	} `json:"redirects"`
}

type VercelFormatter struct{}

func (v *VercelFormatter) FormatRedirects(redirects []Redirect) ([]byte, error) {
	// Vercel format is a JSON object with "redirects" property that is an array of objects
	//	with "source", "destination" and "statusCode" properties
	var formattedRedirects VercelRedirects
	for _, r := range redirects {
		formattedRedirects.Redirects = append(formattedRedirects.Redirects, struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
			StatusCode  int    `json:"statusCode"`
		}{
			Source:      r.From,
			Destination: r.To,
			StatusCode:  301,
		})
	}

	// pretty print json
	b, err := json.MarshalIndent(formattedRedirects, "", "  ")
	if err != nil {
		return nil, err
	}

	return b, nil
}
