package app

import (
	"fmt"
	"math"
)

type MappingRule struct {
	ProjectName string `json:"project"`
	ProjectTask string `json:"projectTask"`
	Role        string `json:"role"`
}

func (mappingRule MappingRule) GetCsvLineString(line TogglTimeEntry, minutesDuration float64) string {
	minutesDurationString := fmt.Sprintf("%.2f", roundToNearest15(minutesDuration))
	return "Project Service," +
		mappingRule.ProjectName + "," +
		mappingRule.ProjectTask + "," +
		mappingRule.Role + "," +
		"Work," +
		"," +
		"," +
		"Draft," +
		line.StartDateTime.Format("1/02/2006") + "," +
		minutesDurationString + "," +
		line.Description + "\"," +
		line.Description + "\"" +
		"\n"
}

func roundToNearest15(minutes float64) float64 {
	interval := 15.0
	val := math.Round(minutes/interval) * interval

	if val <= interval {
		val = interval
	}

	return val
}
