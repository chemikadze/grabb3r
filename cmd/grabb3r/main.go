package main

import (
	"github.com/chemikadze/grabb3r"
	"log"
	"os"
)

func main() {
	user, _ := os.LookupEnv("LEETCODE_USER")
	password, _ := os.LookupEnv("LEETCODE_PASSWORD")
	src := grabb3r.NewLeetCodeSource(user, password)
	dst := grabb3r.NewFileDestination("./dest/file")
	if err := src.Login(); err != nil {
		panic(err)
	}
	if err := dst.Initialize(); err != nil {
		panic(err)
	}
	sync := grabb3r.NewSyncronizer(src, dst)
	if err := sync.Synchronize(); err != nil {
		if err, ok := err.(grabb3r.HttpError); ok {
			log.Fatalf("HttpError: %v %v", err, err.Body)
		}
		panic(err)
	}
}
