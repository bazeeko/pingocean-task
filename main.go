package main

import (
	"context"
	"flag"
	"fmt"
	"time"
)

const MAX_GOROUTINES = 10

type urlSlice []string

func (s *urlSlice) String() string {
	return ""
}

func (s *urlSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	searchString := flag.String("str", "pingocean", "a string to search")

	var urls urlSlice
	flag.Var(&urls, "url", "URL(s) in which the search will be executed")

	flag.Parse()

	if len(urls) == 0 {
		urls = append(urls, "https://pingocean.com/")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchResult := Search(*searchString, urls, ctx)

	for k, v := range searchResult {
		fmt.Printf("URL: %s, Occurences: %d\n", k, v)
	}

	// fmt.Println(SearchStringInURL(url, str))
}
