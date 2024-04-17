package main

import (
	"fmt"
	"github.com/amuwall/go-request"
	"net/http"
)

func main() {
	client, err := request.NewClient(
		"127.0.0.1",
		request.WithScheme("http"),
		request.WithPort(8080),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	params := map[string]interface{}{
		"msg": "hello world",
	}

	req, err := request.NewRequest(
		http.MethodPost,
		"/api/test",
		request.WithBodyParams(
			request.NewJsonBodyParams(params),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", resp)
}
