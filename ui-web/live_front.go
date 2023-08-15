//go:build dev

package main

import (
	"os"
)

func init() {
	front = os.DirFS("front")
}
