package main

import "time"

func main() {
	for {
		select {
		default:
			time.Sleep(1 * time.Hour)
		}
	}
}
