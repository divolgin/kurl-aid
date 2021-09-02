package log

import "fmt"

var (
	currentIndent = ""
)

func IndentMore() {
	currentIndent = currentIndent + "    "
}

func IndentLess() {
	if len(currentIndent) < 4 {
		currentIndent = ""
	} else {
		currentIndent = currentIndent[:len(currentIndent)-4]
	}
}

func Printf(format string, args ...interface{}) {
	if currentIndent == "" {
		fmt.Printf("• ")
	} else {
		fmt.Printf("%s", currentIndent)
	}
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}
