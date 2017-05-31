
package main

import (
	"os"
	"fmt"
	"io/ioutil"
	// "strings"
	"log"
	"regexp"
)

// Color globals
const CURLYBRACECOLOR = "black"
const SQUAREBRACECOLOR = "grey"
const COLONCOLOR = "blue"
const COMMACOLOR = "green"
const BOOLEANCOLOR = "pink"
const STRINGSCOLOR = "brown"
const ESCAPECHARCOLOR = "purple"
const NUMBERSCOLOR = "orange"

// Special symbols
const LT = "&lt;"
const GT = "&gt;"
const AMP = "&amp;"
const QUOT = "&quot;"
const APOS = "&apos;"

// Span template
const SPANTEMPLATE = "<span style=\"color:%s\">%s</span>"
const INDENTSPANTEMPLATE = "<span style=\"font-family:monospace; white-space:pre\">%v</span>"

// Regexp 
// Rexexp switch case: https://groups.google.com/forum/#!topic/golang-nuts/QSlnvdmyCvE
var REGEXCURLYBRACE = regexp.MustCompile(`{|}`)
var REGEXJSONKEY = regexp.MustCompile(`\"`) // assume no error
var REGEXACTUALJSONKEY = regexp.MustCompile(`\".*\"`)
var REGEXCOLON = regexp.MustCompile(`:`)
var REGEXCOMMA = regexp.MustCompile(`,`)
var REGEXVALUE = regexp.MustCompile(`.`)
var REGEXVALUETRUE = regexp.MustCompile(`[t][r][u][e]`)
var REGEXVALUEFALSE = regexp.MustCompile(`[f][a][l][s][e]`)
var REGEXVALUENULL = regexp.MustCompile(`[n][u][l][l]`)
var REGEXVALUEARRAY = regexp.MustCompile(`\[.*\]`)
var REGEXWHITESPACE = regexp.MustCompile(`\s`)

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
	rawFileString, rawFileStringLen := readFileToString(filename)
	var token string
	
	for rawFileStringLen > 0  {
		token, rawFileString, rawFileStringLen = scanJSON(rawFileString)
		colorAndFormat(token)

	}
	// fmt.Println(token, rawFileString, rawFileStringLen)
	return ""
}

// string -> string
// Reads a file into a string variable
func readFileToString(filename string) (string, int) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	result := string(fileContent)
	// fmt.Println("Result: ", resultString)
	return result, len(result)
}

// string -> string, []string, int
// Tokenizes JSON input by splitting them up into 
// a list of strings
// Invariant: incoming string is non zero length
func scanJSON(rawFileString string) (string, string, int) {
	var restTokens string 
	var token string
	// var colorCode string
	tokenBegin := string(rawFileString[0])
	// fmt.Println()
	// fmt.Println(tokenBegin, " -> ", rawFileString[1:len(rawFileString)])
	switch {
	case REGEXCURLYBRACE.MatchString(tokenBegin):
		// fmt.Println(tokenBegin, " is a curly")
		token = tokenBegin
		restTokens = rawFileString[1:len(rawFileString)]
	case REGEXJSONKEY.MatchString(tokenBegin):
		// copy whole key
		token = REGEXACTUALJSONKEY.FindString(rawFileString)
		// fmt.Println("found key ", token)
		restTokens = rawFileString[len(token):len(rawFileString)]
	case REGEXCOLON.MatchString(tokenBegin):
		// fmt.Println(tokenBegin, " is a colon")
		token = tokenBegin
		restTokens = rawFileString[1:len(rawFileString)]
	case REGEXCOMMA.MatchString(tokenBegin):
		// fmt.Println(tokenBegin, " is a comma")
		token = tokenBegin
		restTokens = rawFileString[1:len(rawFileString)]
	case REGEXWHITESPACE.MatchString(tokenBegin):
		restTokens = rawFileString[1:len(rawFileString)]
		// fmt.Println("blank found")
	case REGEXVALUE.MatchString(tokenBegin):
		// case that it is a curly brace...
		if tokenBegin == "t" {
			token = string(REGEXVALUETRUE.FindString(rawFileString[0:6]))
			// fmt.Println("value found: ", token)
			restTokens = rawFileString[len(token):len(rawFileString)]
		} else if tokenBegin == "f" {
			token = REGEXVALUEFALSE.FindString(rawFileString[0:6])
			// fmt.Println("value found: ", token)
			restTokens = rawFileString[len(token):len(rawFileString)]
		} else if tokenBegin == "n" {
			token = REGEXVALUENULL.FindString(rawFileString[0:6])
			// fmt.Println("value found: ", token)
			restTokens = rawFileString[len(token):len(rawFileString)]
		} else if tokenBegin == "[" {
			token = REGEXVALUEARRAY.FindString(rawFileString[0:len(rawFileString)])
			// fmt.Println("value found: ", token)
			restTokens = rawFileString[len(token):len(rawFileString)]
			// restTokens = rawFileString[1:len(rawFileString)]
		} 
		// else {
		// 	restTokens = rawFileString[1:len(rawFileString)]
		// 	fmt.Println(tokenBegin, " is a value")
		// }
	}

	return token, restTokens, len(restTokens)
}

// []string -> ???
// From a list of strings, returns an HTML string 
// that will display the original JSON file contents
// with its tokens colored and properly formatted
func colorAndFormat(token string) string {
	htmlString := ""
	if token != "" {
		spanTest := fmt.Sprintf(SPANTEMPLATE, COLONCOLOR, token)
		fmt.Printf("%v\n", spanTest)
	}
	// fmt.Printf(SPANTEMPLATE, "green", APOS)
	// fmt.Print(spanTest)
	// fmt.Print("\n")
	// fmt.Printf(INDENTSPANTEMPLATE, spanTest)
	


	return htmlString
}