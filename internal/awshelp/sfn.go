package awshelp

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type StateDetails struct {
	Name   string
	Input  any
	Output any
}

type SF struct {
	region            string
	stateMachineArn   string
	maxResults        int64
	fromUnixTimestamp int64
	toUnixTimestamp   int64
	sess              *session.Session
	client            *sfn.SFN
}

func NewSF(stateMachineArn string, fromUnixTimestamp, toUnixTimestamp int64) (sf SF) {
	sf = SF{
		region:            "sa-east-1",
		stateMachineArn:   stateMachineArn,
		maxResults:        1000,
		fromUnixTimestamp: fromUnixTimestamp,
		toUnixTimestamp:   toUnixTimestamp,
		sess: session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})),
	}
	sf.client = sfn.New(sf.sess, &aws.Config{Region: &sf.region})
	return
}

func (sf *SF) TrimForJson(events []*sfn.HistoryEvent) (ret []map[string]any) {
	for _, event := range events {
		item := map[string]any{}
		val := reflect.ValueOf(event).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			if field.Kind() == reflect.Pointer && !field.IsNil() {
				item[val.Type().Field(i).Name] = field.Interface()
			}
		}
		ret = append(ret, item)
	}
	return
}

func (sf *SF) ListExecutions(nextToken *string, statusFilter string) ([]*sfn.ExecutionListItem, *string) {
	params := sfn.ListExecutionsInput{
		StateMachineArn: &sf.stateMachineArn,
		MaxResults:      &sf.maxResults,
	}
	if nextToken != nil && *nextToken != "" {
		params.NextToken = nextToken
	}
	if statusFilter != "" {
		params.StatusFilter = &statusFilter
	}
	ret, err := sf.client.ListExecutions(&params)
	if err != nil {
		panic(err)
	}
	var filtered []*sfn.ExecutionListItem
	if sf.fromUnixTimestamp == 0 && sf.toUnixTimestamp == 0 {
		filtered = ret.Executions
	} else {
		for _, exec := range ret.Executions {
			if (sf.fromUnixTimestamp == 0 || exec.StartDate.Unix() >= sf.fromUnixTimestamp) &&
				(sf.toUnixTimestamp == 0 || exec.StartDate.Unix() < sf.toUnixTimestamp) {
				filtered = append(filtered, exec)
			}
		}
	}
	return filtered, ret.NextToken
}

func (sf *SF) GetExecutionHistory(executionArn string, includeExecutionData bool) (events []*sfn.HistoryEvent) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *sfn.ExecutionDoesNotExist:
				fmt.Fprintf(os.Stderr, "\n%s\n", err)
			default:
				panic(err)
			}
		}
	}()
	params := sfn.GetExecutionHistoryInput{
		ExecutionArn:         &executionArn,
		IncludeExecutionData: &includeExecutionData,
		MaxResults:           &sf.maxResults,
	}
	for {
		ret, err := sf.client.GetExecutionHistory(&params)
		if err != nil {
			panic(err)
		}
		events = append(events, ret.Events...)
		if ret.NextToken == nil {
			break
		}
		params.NextToken = ret.NextToken
	}
	return
}

func (sf *SF) DescribeExecution(executionArn string) *sfn.DescribeExecutionOutput {
	params := sfn.DescribeExecutionInput{
		ExecutionArn: &executionArn,
	}
	ret, err := sf.client.DescribeExecution(&params)
	if err != nil {
		panic(err)
	}
	return ret
}

func (sf *SF) StartExecution(jsonInput *string) *sfn.StartExecutionOutput {
	params := sfn.StartExecutionInput{
		StateMachineArn: &sf.stateMachineArn,
		Input:           jsonInput,
	}
	ret, err := sf.client.StartExecution(&params)
	if err != nil {
		panic(err)
	}
	return ret
}

func (sf *SF) RestartExecution(exec *sfn.ExecutionListItem) *sfn.StartExecutionOutput {
	desc := sf.DescribeExecution(*exec.ExecutionArn)
	return sf.StartExecution(desc.Input)
}

func (sf *SF) GetStateDetailsByName(events []*sfn.HistoryEvent, fn func(name string) bool) (ret StateDetails, ok bool) {
	enterSet, exitSet := false, false
	for i := 0; i < len(events); i++ {
		if events[i].StateEnteredEventDetails != nil && fn(*events[i].StateEnteredEventDetails.Name) {
			ret.Name = *events[i].StateEnteredEventDetails.Name
			if err := json.Unmarshal([]byte(*events[i].StateEnteredEventDetails.Input), &ret.Input); err != nil {
				panic(err)
			}
			enterSet = true
		} else if events[i].StateExitedEventDetails != nil && fn(*events[i].StateExitedEventDetails.Name) {
			if ret.Name != *events[i].StateExitedEventDetails.Name {
				panic("Ambiguous name")
			}
			if err := json.Unmarshal([]byte(*events[i].StateExitedEventDetails.Output), &ret.Output); err != nil {
				panic(err)
			}
			exitSet = true
		}
		if enterSet && exitSet {
			ok = true
			break
		}
	}
	return
}
