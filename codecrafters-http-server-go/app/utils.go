package main

import "time"

func now() string {
	return time.Now().Format(time.RFC3339)
}
