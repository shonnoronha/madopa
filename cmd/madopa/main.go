package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shonnnoronha/madopa/pkg/madopa"
)

func main() {
	inputFile := flag.String("input", "", "Input markdown file")
	outputFile := flag.String("output", "", "Output HTML file")
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

	html, err := madopa.ConvertToHTML(string(content))
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
}

func replaceExt(filename, newExt string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext) + newExt
}
