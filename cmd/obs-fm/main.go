package main

import (
	"errors"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("could not execute utility: %s", err)
	}
}

func run() error {
	return errors.New("not implemented")
}
