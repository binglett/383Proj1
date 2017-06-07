/*
	CMPT383 Project 1
	JSON Formatter
	Bonnie Ng (301223584)
	20170606
*/

package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

// Color globals
const NORMALCOLOR = "maroon"
const CURLYBRACECOLOR = "black"
const SQUAREBRACECOLOR = "indianred"
const COLONCOLOR = "grey"
const COMMACOLOR = "maroon"
const BOOLEANCOLOR = "blue"
const ESCAPECHARCOLOR = "fuchsia"
const NUMBERSCOLOR = "green"

// JSON Types
type JSONElem int
const LEFTCURLYBRACE JSONElem = 1
const RIGHTCURLYBRACE JSONElem = 2
const SQUAREBRACE JSONElem = 3
const COMMA JSONElem = 4
const COLON JSONElem = 5
const BOOLNULL JSONElem = 6
const KEY JSONElem = 7
const ARRAY JSONElem = 8
const NUMBER JSONElem = 9

// Special symbols
const LT string = "&lt;"
const GT string = "&gt;"
const AMP string = "&amp;"
const QUOT string = "&quot;"
const APOS string = "&apos;"

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
var REGEXVALUEARRAY = regexp.MustCompile(`\[|\]`)
var REGEXWHITESPACE = regexp.MustCompile(`\s`)
var REGEXNUMBER = regexp.MustCompile(`-|\d`)
var REGEXNUMBERS = regexp.MustCompile(`\-?\d+\.?\d+((e\+)|(e\-)|(E\-)|(E\+)|(e)|(E))?\d+`)
// List of escape characters: http://www.json.org/
var REGEXESCAPECHARS = regexp.MustCompile(`(\\[bfnrt\\\/"])|(\\u(\d\d\d\d))`)


///////////////////////
// Main Function
///////////////////////
func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: go run a2 <filename>")
		os.Exit(1)
	} else {
		filename := os.Args[1]
		formatJSON(filename)
	}

}

//////////////////////////////////////////
// Main Tokenizer and Format Functions
//////////////////////////////////////////

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

// string -> string, string, int, JSONElem
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
		} else if tokenBegin == "[" || tokenBegin == "]" {
			token = REGEXVALUEARRAY.FindString(rawFileString[0:1])
			tokenType = ARRAY
		} else if REGEXNUMBER.MatchString(tokenBegin) {
			token = REGEXNUMBERS.FindString(rawFileString)
			tokenType = NUMBER
		}
	}
	if len(token) > 0 {
		restTokens = rawFileString[len(token):len(rawFileString)]		
	} else {
		restTokens = rawFileString[1:len(rawFileString)]		
	}
	return token, restTokens, len(restTokens), tokenType
}

// string JSONElem *int -> string
// From a token, including its type, and the current scope returns an HTML string 
// that will display the original JSON file contents properly colored and indented
// with its tokens colored and properly formatted
func colorAndFormat(token string, tokenType JSONElem, scopeNumber *int) string {
	htmlString := ""
	var spanTest string 
	if token != "" {
		if tokenType == LEFTCURLYBRACE {			
			spanTest = fmt.Sprintf(SPANTEMPLATE, CURLYBRACECOLOR, token)
			*scopeNumber++
		} else if tokenType == RIGHTCURLYBRACE {
			spanTest = fmt.Sprintf(SPANTEMPLATE, CURLYBRACECOLOR, token)
			*scopeNumber--
			spanTest = "</br>" + makeWhitespace(scopeNumber) + spanTest
		} else if tokenType == COMMA {
			spanTest = fmt.Sprintf(SPANTEMPLATE, COMMACOLOR, token)	
			spanTest += " "
		} else if tokenType == COLON {
			spanTest = fmt.Sprintf(SPANTEMPLATE, COLONCOLOR, token)	
			spanTest = " " + spanTest + " "
		} else if tokenType == BOOLNULL {
			spanTest = fmt.Sprintf(SPANTEMPLATE, BOOLEANCOLOR, token)	
		} else if tokenType == KEY {
			tmp := replaceSpecialCharsInKey(token)
			spanTest = colorKey(tmp, scopeNumber)
		} else if tokenType == ARRAY {
			spanTest = fmt.Sprintf(SPANTEMPLATE, SQUAREBRACECOLOR, token)
		} else if tokenType == NUMBER {
			spanTest = fmt.Sprintf(SPANTEMPLATE, NUMBERSCOLOR, token)
		} else {
			spanTest = fmt.Sprintf(SPANTEMPLATE, NORMALCOLOR, token)
		}
		fmt.Printf("%v", spanTest)
	}
	return htmlString
}

///////////////////////////////////
// Helper Functions
///////////////////////////////////

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

// string -> string 
// Replaces special characters within a token 
func replaceSpecialCharsInKey(token string) string {
	for i := 0; i < len(token); i++ {
		if string(token[i]) == "&" {
			token = string(token[0:i]) + AMP + string(token[i+1:len(token)])
			i += len(AMP) - 1
		} else if string(token[i]) == "<" {
			token = string(token[0:i]) + LT + string(token[i+1:len(token)])
			i += len(LT) - 1
		} else if string(token[i]) == ">" {
			token = string(token[0:i]) + GT + string(token[i+1:len(token)])
			i += len(GT) - 1
		} else if string(token[i]) == "\"" {
			token = string(token[0:i]) + QUOT + string(token[i+1:len(token)])
			i += len(QUOT) - 1
		} else if string(token[i]) == "'" {
			token = string(token[0:i]) + APOS + string(token[i+1:len(token)])
			i += len(APOS) - 1
		} 
	}
	return token
}

// string, *int -> string
// Given a token string, returns html fragment(s) where 
// escape sequences are colored
// Unicode escape info: https://stackoverflow.com/questions/3900919/write-a-program-to-check-if-a-character-forms-an-escape-character-in-c
func colorKey(token string, scopeNumber *int) string {
	var fullHTMLSpan string
	fullHTMLSpan += "\n"
	fullHTMLSpan += makeWhitespace(scopeNumber)
	prevTokenBegin := 0
	for i := 0; i < len(token); i++ {
		if string(token[i]) == "\\" {
			escapeSeq := REGEXESCAPECHARS.FindString(token[i:len(token)])			
			if escapeSeq != "" {
				if prevTokenBegin < (prevTokenBegin + len(escapeSeq)) {
					fullHTMLSpan += fmt.Sprintf(SPANTEMPLATE, NORMALCOLOR, token[prevTokenBegin:i])
				}
				fullHTMLSpan += fmt.Sprintf(SPANTEMPLATE, ESCAPECHARCOLOR, escapeSeq)
				prevTokenBegin = (i + len(escapeSeq))
			}
		}
	}
	if prevTokenBegin <= len(token)+ 1 {
		fullHTMLSpan += fmt.Sprintf(SPANTEMPLATE, NORMALCOLOR, token[prevTokenBegin:len(token)])
	}	
	return fullHTMLSpan
}

// *int -> string
// Given the scope number, creates the numbers of spaces 
// int * 4
func makeWhitespace(scopeNumber *int) string {
	whitespace := ""
	if *scopeNumber > 0 {
		var spaceNumber int = *scopeNumber * 4
		for spaceNumber > 0 {
			whitespace += " "
			spaceNumber--			
		}
	}
	return whitespace
}