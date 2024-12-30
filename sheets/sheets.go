package sheets

import (
	"context"
	"easyrider"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"strconv"
)

type RedirectSource struct {
	SheetID                       string
	SheetRange                    string
	client                        *sheets.Service
	ServiceAccountCredentialsFile string
}

// NewRedirectSource returns a new instance of RedirectSource with the given sheet ID.
func NewRedirectSource(sheetID string) *RedirectSource {
	return &RedirectSource{
		SheetID: sheetID,
	}
}

func (s *RedirectSource) Authenticate(ctx context.Context) error {
	b, err := os.ReadFile(s.ServiceAccountCredentialsFile)
	if err != nil {
		return err
	}
	conf, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	client := conf.Client(ctx)
	s.client, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	return nil
}

func (s *RedirectSource) LoadRedirects(ctx context.Context) ([]easyrider.Redirect, error) {
	err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	
	resp, err := s.client.Spreadsheets.Values.Get(s.SheetID, s.SheetRange).Do()
	if err != nil {
		return nil, err
	}

	var redirects []easyrider.Redirect
	for _, row := range resp.Values {
		// all values come in as string from the sheets API,
		// convert status code to an int
		statusCode, err := strconv.Atoi(row[2].(string))
		if err != nil {
			return nil, err
		}

		redirects = append(redirects, easyrider.Redirect{
			From:       row[0].(string),
			To:         row[1].(string),
			StatusCode: statusCode,
		})
	}

	return redirects, nil
}
