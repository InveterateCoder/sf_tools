package executors

import (
	"flag"
	"fmt"
	"os"
	"sf_tools/internal/processors"
)

func ExecuteMap(command string) {
	parallelListExecutions := flag.Int("pList", 1, "Parallel listing executions (>= 1)")
	parallelExecutionHistories := flag.Int("pExec", 200, "Parallel executions (>= 1)")
	flag.Parse()

	if len(flag.Args()) != 1 ||
		*parallelListExecutions < 1 ||
		*parallelExecutionHistories < 1 {
		fmt.Fprintf(
			os.Stderr,
			"Usage:\n  %s %s\n\n",
			os.Args[0],
			command,
		)
		flag.PrintDefaults()
		return
	}
	processors.Map(*parallelListExecutions, *parallelExecutionHistories)
}
