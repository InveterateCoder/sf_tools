package misc

import (
	"encoding/json"
	"errors"
	"os"
	"sf_tools/internal/awshelp"
	"sync"
	"time"
)

type State struct {
	ExecutionArn string               `json:"execution_arn"`
	Timestamp    time.Time            `json:"timestamp"`
	Details      awshelp.StateDetails `json:"details"`
}

const _file_name = "states.json"

var appendMutex sync.Mutex
var _states *[]State

func GetStates() *[]State {
	if _states == nil {
		var (
			file *os.File
			err  error
		)
		file, err = os.Open(_file_name)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				file, err = os.OpenFile(_file_name, os.O_CREATE|os.O_WRONLY, 0755)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				_states = &[]State{}
				err = json.NewEncoder(file).Encode(_states)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			defer file.Close()
			var states []State
			err = json.NewDecoder(file).Decode(&states)
			if err != nil {
				panic(err)
			}
			_states = &states
		}
	}
	return _states
}

func FlushStates() {
	if _states == nil {
		return
	}
	file, err := os.OpenFile(_file_name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(_states); err != nil {
		panic(err)
	}
}

func AppendToStates(state State) {
	appendMutex.Lock()
	defer appendMutex.Unlock()
	*_states = append(*_states, state)
}
