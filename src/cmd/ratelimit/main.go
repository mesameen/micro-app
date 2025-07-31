package main

import (
	"fmt"

	"golang.org/x/time/rate"
)

func main() {
	limit := 3
	burst := 3
	limiter := rate.NewLimiter(rate.Limit(limit), burst)
	for i := 0; i < 100; i++ {
		if limiter.Allow() {
			fmt.Println("allowed")
		} else {
			fmt.Println("not allowed")
		}
	}
}
