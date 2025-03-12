package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/shonnnoronha/madopa/internal/renderer"
	"github.com/shonnnoronha/madopa/pkg/madopa"
)

func main() {
	inputFile := flag.String("input", "", "Input markdown file")
	outputFile := flag.String("output", "", "Output HTML file")
	serverFlag := flag.Bool("serve", false, "Serve the generated HTML file")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: File file is required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFile == "" {
		*outputFile = replaceExt(*inputFile, ".html")
	}

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error while reading File %v\n", err)
		os.Exit(1)
	}

	html, err := madopa.Convert(string(content), renderer.NewHTMLRenderer(nil))
	if err != nil {
		fmt.Printf("Error parsing markdow %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, []byte(html), 0644)
	if err != nil {
		fmt.Printf("Error writing to the file %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s \n", *inputFile, *outputFile)

	if *serverFlag {
		serverHTML(*outputFile)
	}
}

func replaceExt(filename, newExt string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext) + newExt
}

func serverHTML(htmlFile string) {
	fileServer := http.FileServer(http.Dir(filepath.Dir(htmlFile)))
	http.Handle("/", fileServer)

	baseName := filepath.Base(htmlFile)
	fmt.Printf("Serving %s at http://localhost:3000/%s\n", htmlFile, baseName)
	fmt.Println("Press Ctrl+C to stop the server")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
