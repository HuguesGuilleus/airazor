package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/HuguesGuilleus/airazor"
)

var verbose = flag.Bool("v", false, "Enable verbose mode")

func main() {
	path := flag.String("p", "collection.json", "Path file of the JSON collection")
	flag.Parse()

	collection, err := read(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read %q: %v\r\n", *path, err)
		return
	}

	fmt.Println("fetch ...")
	collection.Fetch(&airazor.Config{
		NewContext:   context.Background,
		RoundTripper: http.DefaultTransport,
		LimitBody:    500_000,
	})

	printall("", collection)
}

func read(path string) (*airazor.Collection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return airazor.ParseCollection(data)
}

func printall(name string, collection *airazor.Collection) {
	name += collection.Name + " > "
	for _, child := range collection.Children {
		printall(name, child)
	}

	for _, r := range collection.Requests {
		fmt.Println(name + r.Name)
		fmt.Printf("%s %s\r\n", r.Method, r.URL)

		if r.Error != "" {
			fmt.Println("error:", r.Error)
		} else {
			resp := r.Response
			fmt.Println(resp.StatusCode)

			if *verbose {
				fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
				resp.Header.Write(os.Stdout)
				fmt.Println()
				fmt.Println(string(resp.Body))
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			}

			if len(resp.TestFails) > 0 {
				for _, t := range resp.TestFails {
					fmt.Println("-\t" + strings.Join(strings.Split(t, "\n"), "\r\n\t"))
				}
			}
		}
	}
}
