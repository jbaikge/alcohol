package main

import "os"

func env(key, def string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = def
	}
	return
}
