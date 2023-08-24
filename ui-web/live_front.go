//go:build dev

// to run it: go run -tags=dev .

package main

import (
	"os"
)

func init() {
	front = os.DirFS("front")
}
