# JSON specified HTTP client

[![GoDoc](https://godoc.org/github.com/koron-go/jsonhttpc?status.svg)](https://godoc.org/github.com/koron-go/jsonhttpc)
[![Actions/Go](https://github.com/koron-go/jsonhttpc/workflows/Go/badge.svg)](https://github.com/koron-go/jsonhttpc/actions?query=workflow%3AGo)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/jsonhttpc/master.svg)](https://circleci.com/gh/koron-go/jsonhttpc/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/jsonhttpc)](https://goreportcard.com/report/github.com/koron-go/jsonhttpc)

Package jsonhttpc provides JSON specialized HTTP Client.

*   request and response is encoded/decoded as JSON automatically.
*   bit convenient for repeating requests.
    *   specify base URL - `WithBaseURL()`
    *   specify HTTP Client - `WithClient()`
    *   specify HTTP header - `WithHeader()`
*   tips to customize requests - `Do()`
    *   `ContentType() string` on `body` overrides "Content-Type" header.
        (default is "application/json")
*   error response is decoded as JSON - `Error`
    *   Most of JSON properties are put into
        `Properties map[string]interface{}`
    *   There are some fields for
        [RFC7808 Problem Details for HTTP APIs][rfc7808]

(Japanese)

*   リクエストとレスポンスは自動的にJSONエンコード/デコードされます
*   リクエストを繰り返し行うのに少し便利です
    *   ベースURLを設定できる `WithBaseURL()`
    *   HTTP Clientを設定できる `WithClient()`
    *   ヘッダーを設定できる `WithHeader()`
*   リクエストをカスタマイズするtipsがあります `Do()` 
    *   `body` に `ContentType() string` を実装するとContent-Typeヘッダーを変更
        できます (デフォルトは `application/json`)
*   エラーの際もレスポンスは自動的にJSONデコードされます
    *   `Error.Properties` に `map[string]interface{}` で入ります
    *   [RFC7808 Problem Details for HTTP APIs][rfc7808] は少し楽できます

## How to use

[rfc7808]:https://tools.ietf.org/html/rfc7807
