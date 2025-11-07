package app_test

import (
	"testing"

	"github.com/samuel-poirier/toggl-to-dynamics-timesheet-converter/cmd/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	t.Run("get line string should return expected format and rounding", func(t *testing.T) {
		line := "\"some work\",\"0:25:52\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:08:36\",\"2025-11-05\",\"07:34:28\""
		entry, err := app.ParseLineToTogglTimeEntry(line, 0)

		if assert.NoError(t, err) {
			assert.Equal(t, "some work", entry.Description)
			assert.Equal(t, "example@email.com", entry.Email)
			assert.Equal(t, "test project", entry.Project)
		}
	})
}
