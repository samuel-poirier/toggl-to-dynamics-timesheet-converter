package app

import (
	"fmt"
	"log/slog"
)

const DynamicsHeaders string = "Time Source,Project,Project Task,Role,Type,Subcontract,Subcontract line,Entry Status,Date,Duration,Description,External Comments\n"

func (a *App) Transform(config *AppConfig, togglExportLines []TogglTimeEntry) ([]string, error) {

	linesGroupedMap := make(map[string]TogglTimeEntry)
	minutesPerGroupMap := make(map[string]float64)

	for _, line := range togglExportLines {
		key := line.GetGroupingKey()

		linesGroupedMap[key] = line
		minutesDuration := line.StopDateTime.Sub(line.StartDateTime).Minutes()

		if value, ok := minutesPerGroupMap[key]; ok {
			minutesPerGroupMap[key] = value + minutesDuration
		} else {
			minutesPerGroupMap[key] = minutesDuration
		}

	}

	totalWeekMinutes := 0
	transformedLines := []string{
		DynamicsHeaders,
	}

	for key, line := range linesGroupedMap {
		mappingRule, ok := config.Mapping.Projects[line.Project]

		if !ok {
			return nil, fmt.Errorf("failed to find project [%s] from mapping rules", line.Project)
		}

		// the key exists for sure here
		minutesDuration := minutesPerGroupMap[key]

		totalWeekMinutes += int(minutesDuration)

		lineString := mappingRule.GetCsvLineString(line, minutesDuration)
		transformedLines = append(transformedLines, lineString)
	}

	a.logger.Info("weekly total hours", slog.Any("hours", totalWeekMinutes/60))

	return transformedLines, nil
}
