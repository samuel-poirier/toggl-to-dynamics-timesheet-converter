package app

import (
	"bufio"
	"fmt"
	"log/slog"
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

	linesGroupedMap := make(map[string]TogglTimeEntry)
	minutesPerGroupMap := make(map[string]float64)

	for _, line := range togglExportLines {
		key := line.GetGroupingKey()

		linesGroupedMap[key] = line
		minutesDuration := line.stopDateTime.Sub(line.startDateTime).Minutes()

		if value, ok := minutesPerGroupMap[key]; ok {
			minutesPerGroupMap[key] = value + minutesDuration
		} else {
			minutesPerGroupMap[key] = minutesDuration
		}

	}

	totalWeekMinutes := 0
	for key, line := range linesGroupedMap {
		mappingRule, ok := config.Mapping.Projects[line.project]

		if !ok {
			return fmt.Errorf("failed to find project [%s] from mapping rules", line.project)
		}

		minutesDuration, ok := minutesPerGroupMap[key]

		if !ok {
			return fmt.Errorf("failed to find calculated minutes from grouped lines for project [%s]", line.project)
		}

		totalWeekMinutes += int(minutesDuration)

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

	a.logger.Info("weekly total hours", slog.Any("hours", totalWeekMinutes/60))

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
