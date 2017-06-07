package main

import (
	"time"

	_ "github.com/dup2X/dxtrace"
)

func main() {
	tk := time.NewTicker(time.Millisecond * 100)
	for range tk.C {
	}
}
