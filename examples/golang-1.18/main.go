package main

import (
	doless "github.com/hedlx/doless/runtime/golang-1.18"
)

func main() {
	doless.Lambda(func(req *doless.Request) (int, string) {
		return 200, string(req.Payload)
	})
}
