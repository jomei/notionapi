# notionapi
An API client for the [Notion API](https://developers.notion.com/) implemented in Golang

# Roadmap
- [x] Databases
- [x] Pages
- [x] Pages
- [x] Blocks
- [x] Users
- [x] Search
- [ ] Tests and Examples

# Installation

```
$ go get github.com/jomei/notionapi
```

# Getting started
Follow Notion’s [getting started guide](https://developers.notion.com/docs/getting-started) to obtain an Integration Token.

## Example

Make a new `Client`

```go
import "github.com/jomei/notionapi"


client := notionapi.NewClient("your-integration-token")
```
Then, use client's methods to retrieve or update your content

```go
page, err := client.PageRetrieve(context.Background(), "your-page-id")
if err != nil {
	// do something
}
```
