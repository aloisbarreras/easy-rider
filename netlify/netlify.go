package netlify

import (
	"easyrider"
	"fmt"
	"strings"
)

type RedirectFormatter struct{}

func (n *RedirectFormatter) FormatRedirects(redirects []easyrider.Redirect) ([]byte, error) {
	//Netlify _redirects format looks like:
	// /from.html	/to.html		200

	var formattedRedirects []string
	for _, r := range redirects {
		formattedRedirects = append(formattedRedirects, fmt.Sprintf("%s %s %d", r.From, r.To, r.StatusCode))
	}

	return []byte(strings.Join(formattedRedirects, "\n") + "\n"), nil
}
