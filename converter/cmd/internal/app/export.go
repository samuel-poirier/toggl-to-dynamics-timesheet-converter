package app

import (
	"bufio"
	"os"
)

func (a *App) Export(config *AppConfig, dynamicsExportLines []string) error {
	file, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range dynamicsExportLines {
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	// Flush the buffer to ensure all data is written to the file
	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}
