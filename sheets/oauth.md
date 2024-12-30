### Example OAuth Flow code
This code will read a `credentials.json` file from disk and prompt the user with a URL to grant OAuth access.
The functions are in `oauth.go`.

```go
package main 

func authenticate() {
    b, err := os.ReadFile("/Users/aloisbarreras/Downloads/credentials.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }
    
    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client := getClient(config)
    
    srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }
    
    s.client = srv
}
```
