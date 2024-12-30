package vercel

import (
	"bytes"
	"easyrider"
	"encoding/json"
)

// Vercel redirects format is a JSON object with "redirects" property that is an array of objects
// with "source", "destination" and "statusCode" properties
type vercelRedirect struct {
	Redirects []struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		StatusCode  int    `json:"statusCode"`
	} `json:"redirects"`
}

type RedirectFormatter struct{}

func (v *RedirectFormatter) FormatRedirects(redirects []easyrider.Redirect) ([]byte, error) {
	var formattedRedirects vercelRedirect
	for _, r := range redirects {
		formattedRedirects.Redirects = append(formattedRedirects.Redirects, struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
			StatusCode  int    `json:"statusCode"`
		}{
			Source:      r.From,
			Destination: r.To,
			StatusCode:  r.StatusCode,
		})
	}

	b := &bytes.Buffer{}
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(formattedRedirects)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
