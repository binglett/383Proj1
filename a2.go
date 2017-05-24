package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"log"
)

// Color globals
const CURLYBRACECOLOR = ""
const SQUAREBRACECOLOR = ""
const COLONCOLOR = ""
const COMMACOLOR = ""
const BOOLEANCOLOR = ""
const STRINGSCOLOR = ""
const ESCAPECHARCOLOR = ""
const NUMBERSCOLOR = ""

// Special symbols
const LT = "&lt;"
const GT = "&gt;"
const AMP = "&amp;"
const QUOT = "&quot;"
const APOS = "&apos;"

// Span template
const SPANTEMPLATE = "<span style=\"color:%v\">%v</span>"
const INDENTSPANTEMPLATE = "<span style=\"font-family:monospace; white-space:pre\">%v</span>"

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: go run a2 <filename>")
		os.Exit(1)
	} else {
		filename := os.Args[1]
		formatJSON(filename)
	}

	// TMP testing 
}

// string -> string
// Returns the colored and formatted html string of 
// a valid JSON 
func formatJSON(filename string) string {
	rawFileString := reaadFileToString(filename)
	tokens := scanJSON(rawFileString)
	return colorAndFormat(tokens)
}

// string -> string
// Reads a file into a string variable
func reaadFileToString(filename string) []string {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Fields(string(fileContent))
}

// string -> []string
// Tokenizes JSON input by splitting them up into 
// a list of strings
func scanJSON(rawFileString []string) []string {
	var tokens []string 
	tokens = rawFileString
	// TODO split into tokens

	return tokens
}

// []string -> ???
// From a list of strings, returns an HTML string 
// that will display the original JSON file contents
// with its tokens colored and properly formatted
func colorAndFormat(tokens []string) string {
	htmlString := ""
	spanTest := fmt.Sprintf(SPANTEMPLATE, "green", tokens)
	// fmt.Printf(SPANTEMPLATE, "green", APOS)
	// fmt.Print(spanTest)
	// fmt.Print("\n")
	// fmt.Printf(INDENTSPANTEMPLATE, spanTest)
	
	fmt.Printf("%v\n", spanTest)


	return htmlString
}