/*
	CMPT383 Project 1
	JSON Formatter
	Bonnie Ng (301223584)
	20170530
*/

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
const NORMALCOLOR = "green"
const CURLYBRACECOLOR = "black"
const SQUAREBRACECOLOR = "purple"
const COLONCOLOR = "blue"
const COMMACOLOR = "maroon"
const BOOLEANCOLOR = "pink"
const ESCAPECHARCOLOR = "orange"
const NUMBERSCOLOR = "yellow"

// Types
type JSONElem int
const LEFTCURLYBRACE JSONElem = 1
const RIGHTCURLYBRACE JSONElem = 2
const SQUAREBRACE JSONElem = 3
const COMMA JSONElem = 4
const COLON JSONElem = 5
const BOOLNULL JSONElem = 6
const KEY JSONElem = 7
const ARRAY JSONElem = 8

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
var REGEXLEFTCURLYBRACE = regexp.MustCompile(`{`)
var REGEXRIGHTCURLYBRACE = regexp.MustCompile(`}`)
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

}

// string -> string
// Returns the colored and formatted html string of 
// a valid JSON 
func formatJSON(filename string) {
	rawFileString, rawFileStringLen := readFileToString(filename)
	var token string
	var tokenType JSONElem
	scopeNumber := 0

	begin := "<span style=\"font-family:monospace; white-space:pre\">"
	fmt.Printf("%v", begin)

	for rawFileStringLen > 0  {
		token, rawFileString, rawFileStringLen, tokenType = scanJSON(rawFileString)
		colorAndFormat(token, tokenType, &scopeNumber)

	}
	end := "</span>"
	fmt.Printf("%v", end)
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

// string -> string, string, int, string
// Tokenizes JSON input by returning the next token,
// the rest of the unparsed string, the length of the 
// unparsed string, and what JSON element the token is (for further processing)
// Invariant: incoming string is non zero length
func scanJSON(rawFileString string) (string, string, int, JSONElem) {
	var restTokens string 
	var token string
	var tokenType JSONElem
	tokenBegin := string(rawFileString[0])
	switch {
	case REGEXLEFTCURLYBRACE.MatchString(tokenBegin):
		token = tokenBegin
		tokenType = LEFTCURLYBRACE
	case REGEXRIGHTCURLYBRACE.MatchString(tokenBegin):
		token = tokenBegin
		tokenType = RIGHTCURLYBRACE
	case REGEXJSONKEY.MatchString(tokenBegin):
		token = REGEXACTUALJSONKEY.FindString(rawFileString)
		tokenType = KEY
	case REGEXCOLON.MatchString(tokenBegin):
		token = tokenBegin
		tokenType = COLON
	case REGEXCOMMA.MatchString(tokenBegin):
		token = tokenBegin
		tokenType = COMMA 
	case REGEXWHITESPACE.MatchString(tokenBegin):
		// Do nothing
	case REGEXVALUE.MatchString(tokenBegin):
		if tokenBegin == "t" {
			token = string(REGEXVALUETRUE.FindString(rawFileString[0:6]))
			tokenType = BOOLNULL
		} else if tokenBegin == "f" {
			token = REGEXVALUEFALSE.FindString(rawFileString[0:6])
			tokenType = BOOLNULL
		} else if tokenBegin == "n" {
			token = REGEXVALUENULL.FindString(rawFileString[0:6])
			tokenType = BOOLNULL
		} else if tokenBegin == "[" {
			token = REGEXVALUEARRAY.FindString(rawFileString[0:len(rawFileString)])
			tokenType = ARRAY
		} 
	}
	if len(token) > 0 {
		restTokens = rawFileString[len(token):len(rawFileString)]		
	} else {
		restTokens = rawFileString[1:len(rawFileString)]		
	}
	return token, restTokens, len(restTokens), tokenType
}

// string JSONElem int -> string
// From a token, including its type, and the current scope returns an HTML string 
// that will display the original JSON file contents properly colored and indented
// with its tokens colored and properly formatted
func colorAndFormat(token string, tokenType JSONElem, scopeNumber *int) string {
	htmlString := ""
	var spanTest string 
	// whitespace := makeWhitespace(*scopeNumber)
	if token != "" {
		if tokenType == LEFTCURLYBRACE {			
			spanTest = fmt.Sprintf(SPANTEMPLATE, CURLYBRACECOLOR, token)
			spanTest += "</br>"
			*scopeNumber++
		} else if tokenType == RIGHTCURLYBRACE {
			spanTest = fmt.Sprintf(SPANTEMPLATE, CURLYBRACECOLOR, token)
			spanTest = "</br>" + spanTest
			*scopeNumber--
		} else if tokenType == COMMA {
			spanTest = fmt.Sprintf(SPANTEMPLATE, COMMACOLOR, token)	
			spanTest += "</br>"
		} else if tokenType == COLON {
			spanTest = fmt.Sprintf(SPANTEMPLATE, COLONCOLOR, token)	
			spanTest = " " + spanTest + " "
		} else if tokenType == BOOLNULL {
			spanTest = fmt.Sprintf(SPANTEMPLATE, BOOLEANCOLOR, token)	
		} else if tokenType == KEY {
			// tmp := processKey(token)
			spanTest = fmt.Sprintf(SPANTEMPLATE, NORMALCOLOR, token)	
		} else {
			spanTest = fmt.Sprintf(SPANTEMPLATE, NORMALCOLOR, token)
		}
		// fmt.Printf("%s%v", whitespace, spanTest)
		fmt.Printf("%v", spanTest)
	}
	return htmlString
}

// replaces special characters
// colors special characters
func processKey(token string) string {
	
	return ""
}

// !!!
func makeWhitespace(scopeNumber int) string {
	whitespace := ""
	for scopeNumber	> 0 {
		whitespace += "\t"
		scopeNumber--
	}
	return whitespace
}