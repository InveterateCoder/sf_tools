package executors

import (
	"flag"
	"fmt"
	"os"
	"sf_tools/internal/processors"
)

func ExecuteRestart(command string) {
	sfArn := flag.String("sfArn", "", "Step function arn (required)")
	unixTimestampFrom := flag.Int64("from", 0, "Unix timestamp from which the failed executions to be restarted")
	unixTimestampTo := flag.Int64("to", 0, "Unix timestamp until which the failed executions to be restarted")
	flag.Parse()

	if len(flag.Args()) != 1 || *sfArn == "" {
		fmt.Fprintf(
			os.Stderr,
			"Usage:\n  %s %s %s\n\n",
			os.Args[0],
			"-sfArn arn:aws:states:sa-east-1:787732066160:stateMachine:...",
			command,
		)
		flag.PrintDefaults()
		return
	}
	processors.Restart(*sfArn, *unixTimestampFrom, *unixTimestampTo)
}
