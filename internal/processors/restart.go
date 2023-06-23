package processors

import (
	"fmt"
	"sf_tools/internal/awshelp"
	"time"

	"github.com/aws/aws-sdk-go/service/sfn"
)

func Restart(sfArn string, from, to int64) {
	sf := awshelp.NewSF(sfArn, from, to)
	var (
		execs     []*sfn.ExecutionListItem
		nextToken *string
	)
	for i := 0; i == 0 || nextToken != nil; i++ {
		execs, nextToken = sf.ListExecutions(nextToken, "FAILED")
		for j, exec := range execs {
			ret := sf.RestartExecution(exec)
			fmt.Printf("%d | %d\t%s\n", j+1, len(execs), *ret.ExecutionArn)
			time.Sleep(time.Second)
		}
	}
}
