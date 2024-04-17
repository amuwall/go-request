# go-request

`go-request` is a wrapper for `net/http` request to make it easier to send HTTP request.

## Install

```shell
go get -u github.com/amuwall/go-request
```

## Quick Start

```go
// Create client, default scheme is https and port is 443
client, err := request.NewClient("127.0.0.1")

// Create request
req, err := request.NewRequest(http.MethodPost, "/api/test")

// Send request
resp, err := client.do(req)
```

## Advanced

### Client Options

Here are some client options. When call `NewClient`, you can use these to set client.

* WithScheme
* WithPort
* WithTransport
* WithTimeout
* WithClientCertificateBlock
* WithClientCertificateFile
* WithTLSServerName
* WithSkipVerifyCertificates

Example:

```go
client, err := request.NewClient(
    "127.0.0.1",
    request.WithScheme("http"),
    request.WithScheme("80"),
)
```

### Request Options

Here are some request options. When call `NewRequest`, you can use these to set request.

* WithHost
* WithHeaders
* WithQueryParams
* WithBodyParams

Example:

```go
req, err := request.NewRequest(
    http.MethodPost,
    "/api/test",
    request.WithBodyParams(
        request.NewJsonBodyParams(map[string]string{"msg": "hello"}),
    ),
)
```

### Request Body Params

Here are two defined params, `JsonBodyParams` and `FormBodyParams`.

#### JsonBodyParams

```go
params := request.NewJsonBodyParams(map[string]string{"msg": "hello"})
```

#### FormBodyParams

```go
params := request.NewFormBodyParams(map[string]string{"msg": "hello"})
params.AddFile("file", "test.txt", f) // Send file
```

### Response

You can use `.StatusCode` to get response status code, use `.Header` to get response header, and use `.RawBody` to get response body.  
If your response body is JSON, you can call `.UnmarshalJSONBody` method to unmarshal response body

```go
err := resp.UnmarshalJSONBody(val)
```
