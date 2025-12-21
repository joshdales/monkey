package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var traceEnabled bool

func checkTraceEnabled() {
	trace, err := strconv.ParseBool(os.Getenv("PARSER_TRACE"))
	if err != nil {
		traceEnabled = false
	}

	traceEnabled = trace
}

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	if !traceEnabled {
		return
	}

	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	if !traceEnabled {
		return ""
	}

	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	if !traceEnabled {
		return
	}

	tracePrint("END " + msg)
	decIdent()
}
