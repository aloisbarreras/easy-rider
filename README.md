![logo](docs/img/easy-rider-logo.png)
# easy-rider

**easy-rider** is a command-line tool that formats HTTP redirect rules for various platforms.

## Why

Many popular web hosting platforms, such as Netlify and Vercel, manage redirects through configuration files that are checked into your codebase.
It can be challenging for non-technical roles (e.g. marketing/sales teams), to edit a file on GitHub or know how to correctly format a JSON or configuration file,
so you are likely spending time responding to messages or Jira tickets to add/edit a redirect on the website. `easy-rider` allows these non-technical
roles to define HTTP redirects via a familiar interface (e.g. Google Sheets) and generate the necessary configuration files for your existing web infrastructure.

## How it works

1. **Store Redirects**: Store your redirects in one of our supported sources.
2. **Generate Redirects**: Use the `generate` command to pull the redirects and format them for your platform.
3. **Deploy Redirects**: Integrate this tool into your deployment pipeline to apply the generated redirect rules to your platform.

### Example
```bash
easy-rider generate -source google-sheets -sheet-id <SHEET_ID> -format netlify
```

This command pulls redirect rules from the specified Google Sheets file and outputs them in Netlify's `_redirects` format.

```
/from.html      /to.html                200
/from2.html     /to2.html               301
```

Or, the same command with `-format vercel` would output:

```json
{
  "redirects": [
    {
      "source": "/from.html",
      "destination": "/to.html",
      "statusCode": 301
    },
    {
      "source": "/from2.html",
      "destination": "/to2.html",
      "statusCode": 301
    }
  ]
}
```

## Supported Sources

- Google Sheets

## Supported Formats

- `netlify`: Generates `_redirects` files.
- `vercel`: Generates `vercel.json` redirect rules.

## License

MIT License. See `LICENSE` for more details.