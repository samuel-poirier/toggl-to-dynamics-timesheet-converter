package app

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const header string = "Time Source,Project,Project Task,Role,Type,Subcontract,Subcontract line,Entry Status,Date,Duration,Description,External Comments"

func (a *App) Export(config *AppConfig, togglExportLines []TogglTimeEntry) error {
	file, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(header + "\n")

	if err != nil {
		return err
	}

	for _, line := range togglExportLines {
		mappingRule, ok := config.Mapping.Projects[line.project]
		if !ok {
			return fmt.Errorf("faild to find project [%s] from mapping rules", line.project)
		}

		minutesDuration := line.stopDateTime.Sub(line.startDateTime).Minutes()
		minutesDurationString := fmt.Sprintf("%.2f", roundToNearest15(minutesDuration))
		_, err := writer.WriteString(
			"Project Service," +
				mappingRule.ProjectName + "," +
				mappingRule.ProjectTask + "," +
				mappingRule.Role + "," +
				"Work," +
				"," +
				"," +
				"Draft," +
				line.startDateTime.Format("1/02/2006") + "," +
				minutesDurationString + "," +
				"," +
				line.description + "\"" +
				"\n")

		if err != nil {
			return err
		}
	}

	// Flush the buffer to ensure all data is written to the file
	err = writer.Flush()
	return nil
}

func roundToNearest15(minutes float64) float64 {
	interval := 15.0
	val := math.Round(minutes/interval) * interval

	if val <= interval {
		val = interval
	}

	return val
}
