# JSON-specialized HTTP client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/jsonhttpc)](https://pkg.go.dev/github.com/koron-go/jsonhttpc)
[![Actions/Go](https://github.com/koron-go/jsonhttpc/workflows/Go/badge.svg)](https://github.com/koron-go/jsonhttpc/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/jsonhttpc)](https://goreportcard.com/report/github.com/koron-go/jsonhttpc)

Package jsonhttpc provides a way to easily send and receive HTTP requests with JSON bodies.

*   Request and response is encoded/decoded as JSON automatically.
*   Bit handy for repeating requests.
    *   specify base URL - `WithBaseURL()`
    *   specify HTTP Client - `WithClient()`
    *   specify HTTP header - `WithHeader()`
*   Tips to customize requests - for `Do()`
    *   `ContentType() string` on `body` overrides "Content-Type" header.
        (default is "application/json")
*   Responses are automatically JSON-decoded even in case of errors - `Error`
    *   Most of JSON properties are put into
        `Properties map[string]interface{}`
    *   Error responses that follow [RFC7808 Problem Details for HTTP APIs][rfc7808] are a bit easier to handle.

(Japanese)

*   リクエストとレスポンスは自動的にJSONエンコード/デコードされます
*   リクエストを繰り返し行うのに少し便利です
    *   ベースURLを設定できる `WithBaseURL()`
    *   HTTP Clientを設定できる `WithClient()`
    *   ヘッダーを設定できる `WithHeader()`
*   リクエストをカスタマイズするtipsがあります - for `Do()` 
    *   `body` に `ContentType() string` を実装するとContent-Typeヘッダーを変更
        できます (デフォルトは `application/json`)
*   エラーの際もレスポンスは自動的にJSONデコードされます
    *   `Error.Properties` に `map[string]interface{}` で入ります
    *   [RFC7808 Problem Details for HTTP APIs][rfc7808] は少し楽できます

[rfc7808]:https://tools.ietf.org/html/rfc7807

## Install and update

```console
$ go get github.com/koron-go/jsonhttpc@latest
```

## How to use
