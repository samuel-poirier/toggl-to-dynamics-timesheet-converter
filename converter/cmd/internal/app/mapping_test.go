package app_test

import (
	"testing"
	"time"

	"github.com/samuel-poirier/toggl-to-dynamics-timesheet-converter/cmd/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestLineStringConversion(t *testing.T) {
	t.Run("get line string should return expected format and rounding", func(t *testing.T) {
		mappingRule := app.MappingRule{
			ProjectName: "Test Project",
			Role:        "Test Role",
			ProjectTask: "Test Task",
		}

		start, _ := time.Parse("01/02/2006", "11/04/2025")
		end := start.Add(time.Minute * 5)

		entry := app.NewTogglTimeEntry(
			"Some work done",
			time.Minute*6,
			"member",
			"email",
			"Test Project",
			"tags",
			start,
			end,
		)

		minutes := 6.0

		line := mappingRule.GetCsvLineString(entry, minutes)
		assert.Equal(t, "Project Service,Test Project,Test Task,Test Role,Work,,,Draft,11/04/2025,15.00,\"Some work done\",\"Some work done\"\n", line)
	})
}
