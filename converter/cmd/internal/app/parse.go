package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type Mapping struct {
	Projects map[string]MappingRule `json:"projects"`
}

type AppConfig struct {
	Mapping    Mapping
	InputFile  string
	OutputFile string
}

func (a *App) ParseArgs(allArgs []string) (*AppConfig, error) {

	a.logger.Info("running with arguments", slog.Any("allArgs", allArgs))

	mappingIndex := findIndex("-m", allArgs)

	if mappingIndex == -1 || len(allArgs) <= mappingIndex {
		return nil, fmt.Errorf("missing mapping file path argument provided with -m")
	}

	inputFileIndex := findIndex("-i", allArgs)

	if inputFileIndex == -1 || len(allArgs) <= inputFileIndex {
		return nil, fmt.Errorf("missing mapping file path argument provided with -i")
	}

	outputFileIndex := findIndex("-o", allArgs)

	if outputFileIndex == -1 || len(allArgs) <= outputFileIndex {
		return nil, fmt.Errorf("missing output file path argument provided with -o")
	}

	var mapping Mapping
	mappingJson, err := os.ReadFile(allArgs[mappingIndex+1])

	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to read mapping file"))
	}

	err = json.Unmarshal(mappingJson, &mapping)

	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to read mapping file"))
	}

	config := &AppConfig{
		InputFile:  allArgs[inputFileIndex+1],
		OutputFile: allArgs[outputFileIndex+1],
		Mapping:    mapping,
	}

	_, err = os.Stat(config.InputFile)

	if err != nil {
		return nil, fmt.Errorf("toggl csv file not found with path %s", config.InputFile)
	}

	return config, nil
}

func findIndex(lookup string, arr []string) int {

	for i, v := range arr {
		if v == lookup {
			return i
		}
	}

	return -1
}
