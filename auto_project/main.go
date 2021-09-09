package main

import "github.com/tal-tech/go-zero/core/breaker"

func main() {
	br := breaker.NewBreaker()
	err := br.Do(func() error {
	})
}
