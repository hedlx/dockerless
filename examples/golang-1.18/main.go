package main

import (
	"context"
	"io"

	doless "github.com/hedlx/doless/runtime/golang-1.18"
)

func main() {
	doless.Lambda(func(ctx context.Context, req *doless.Request[io.Reader]) (int, interface{}) {
		return 200, req.Payload
	})
}
