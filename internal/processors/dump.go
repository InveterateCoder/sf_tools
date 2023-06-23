package processors

import (
	"fmt"
	"sf_tools/internal/awshelp"
	"sf_tools/internal/misc"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/service/sfn"
)

func Dump() {
	settings := misc.ParseSFMap()
	states := misc.GetStates()
	sf := awshelp.NewSF(settings.SFArn, settings.FromUnixTimestamp, settings.ToUnixTimestamp)
	totalExecs := 0
	proccessedExecs := 0
	var logMutex sync.Mutex
	log := func() {
		logMutex.Lock()
		defer logMutex.Unlock()
		proccessedExecs++
		str := fmt.Sprintf("\rTotal appended: %d\tGot: %d\tProcessed: %d", len(*states), totalExecs, proccessedExecs)
		fmt.Printf("\r%s", strings.Repeat(" ", len(str)+20))
		fmt.Print(str)
	}
	var listWaitGroup sync.WaitGroup
	execsReturned := false
	for i := 0; i == 0 || settings.NextToken != nil; i++ {
		execs, nextToken := sf.ListExecutions(settings.NextToken, "SUCCEEDED")
		totalExecs += len(execs)
		settings.NextToken = nextToken
		if len(execs) > 0 {
			if !execsReturned {
				execsReturned = true
			}
			listWaitGroup.Add(1)
			go func(execs []*sfn.ExecutionListItem) {
				defer listWaitGroup.Done()
				var historyWaitGroup sync.WaitGroup
				for j := 0; j < len(execs); j++ {
					historyWaitGroup.Add(1)
					go func(exec *sfn.ExecutionListItem) {
						defer historyWaitGroup.Done()
						defer log()
						events := sf.GetExecutionHistory(*exec.ExecutionArn, true)
						if stateDetails, ok := sf.GetStateDetailsByName(events, misc.Selector); ok {
							misc.MutateState(&stateDetails)
							state := misc.State{
								ExecutionArn: *exec.ExecutionArn,
								Timestamp:    *exec.StartDate,
								Details:      stateDetails,
							}
							misc.AppendToStates(state)
						}
					}(execs[j])
					if (j+1)%settings.ParallelExecutionHistories == 0 {
						historyWaitGroup.Wait()
					}
				}
				historyWaitGroup.Wait()
			}(execs)
		} else {
			if execsReturned {
				settings.NextToken = nil
				break
			}
			fmt.Printf("\rSkipping execs: %d", (i+1)*1000)
		}
		if (i+1)%settings.ParallelListExecutions == 0 {
			listWaitGroup.Wait()
			misc.FlushSettings()
			misc.FlushStates()
			totalExecs = 0
			proccessedExecs = 0
		}
	}
	listWaitGroup.Wait()
	misc.FlushSettings()
	misc.FlushStates()
	totalExecs = 0
	fmt.Println()
}
