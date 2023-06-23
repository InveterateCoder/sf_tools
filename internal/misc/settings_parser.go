package misc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type SettingsItem struct {
	Input  map[string]string
	Output map[string]string
}

type Settings struct {
	SFArn                      string
	NextToken                  *string
	ParallelListExecutions     int
	ParallelExecutionHistories int
	FromUnixTimestamp          int64
	ToUnixTimestamp            int64
	States                     map[string]SettingsItem
}

var _settings *Settings

func ParseSFMap() *Settings {
	if _settings == nil {
		fmt.Println("Parsing sf_map.json...")
		file, err := os.Open("sf_map.json")
		if err != nil {
			log.Fatalf("Failed to open \"sf_map.json\":\n\t%s\n", err)
		}
		var settings Settings
		err = json.NewDecoder(file).Decode(&settings)
		file.Close()
		if err != nil {
			log.Fatalf("Failed to parse \"sf_map.json\":\n\t%s\n", err)
		}
		_settings = &settings
		if _settings.SFArn == "" {
			log.Fatal("Missing required \"SFArn\" in \"sf_map.json\"\n")
		}
		if _settings.ParallelListExecutions == 0 {
			log.Fatal("Missing required \"ParallelListExecutions\" in \"sf_map.json\"\n")
		}
		if _settings.ParallelExecutionHistories == 0 {
			log.Fatal("Missing required \"ParallelExecutionHistories\" in \"sf_map.json\"\n")
		}
		if _settings.States == nil {
			log.Fatal("\"States\" must be properly specified in \"sf_map.json\"\n")
		}
		for key, val := range _settings.States {
			if key == "" || val.Input == nil || val.Output == nil {
				log.Fatal("\"States\" must be properly specified in \"sf_map.json\"\n")
			}
		}
	}
	return _settings
}

func FlushSettings() {
	if _settings == nil {
		log.Fatal("Calling UpdateNextToken before settings initialization")
	}
	file, err := os.OpenFile("sf_map.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(_settings); err != nil {
		panic(err)
	}
}
