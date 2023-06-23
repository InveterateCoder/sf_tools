package main

import (
	"fmt"
	"os"
	"sf_tools/internal/executors"
	"strings"
)

const (
	mapCommand     = "map"
	restartCommand = "restart"
)

func getListOfCommands() []string {
	return []string{mapCommand, restartCommand}
}

func isCommandAllowed(command string) bool {
	for _, allowed := range getListOfCommands() {
		if command == allowed {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 || !isCommandAllowed(os.Args[len(os.Args)-1]) {
		fmt.Fprintf(
			os.Stderr,
			"Allowed commands:\n\t[%s]\nExample:\n\t%s %s\n",
			strings.Join(getListOfCommands(), " | "),
			os.Args[0],
			getListOfCommands()[0],
		)
		return
	}
	switch os.Args[len(os.Args)-1] {
	case mapCommand:
		executors.ExecuteMap(mapCommand)
	case restartCommand:
		executors.ExecuteRestart(restartCommand)
	default:
		panic("Didn't catch the executor command")
	}
}
