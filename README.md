# notionapi
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/jomei/notionapi?label=go%20module)](https://github.com/conduitio-labs/notionapi/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/conduitio-labs/notionapi.svg)](https://pkg.go.dev/github.com/conduitio-labs/notionapi)
[![Test](https://github.com/conduitio-labs/notionapi/actions/workflows/test.yml/badge.svg)](https://github.com/conduitio-labs/notionapi/actions/workflows/test.yml)

An API client for the [Notion API](https://developers.notion.com/) implemented in Golang

# Supported APIs
It supports all APIs for Notion API version `2022-06-28`

# Installation

```
$ go get github.com/conduitio-labs/notionapi
```

# Getting started
Follow Notion’s [getting started guide](https://developers.notion.com/docs/getting-started) to obtain an Integration Token.

## Example

Make a new `Client`

```go
import "github.com/conduitio-labs/notionapi"


client := notionapi.NewClient("your-integration-token")
```
Then, use client's methods to retrieve or update your content

```go
page, err := client.Page.Get(context.Background(), "your-page-id")
if err != nil {
	// do something
}
```
