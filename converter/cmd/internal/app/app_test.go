package app_test

import (
	"bufio"
	"bytes"
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

	t.Run("parse toggl time entries should skip header line and parse entries as expected", func(t *testing.T) {
		headerLine := "\"Description\",\"Duration\",\"Member\",\"Email\",\"Project\",\"Tags\",\"Start date\",\"Start time\",\"Stop date\",\"Stop time\"\n"
		line1 := "\"some work\",\"0:25:52\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:08:36\",\"2025-11-05\",\"07:34:28\"\n"
		line2 := "\"some work2\",\"0:25:52\",\"some user 2\",\"example2@email.com\",\"test project 2\",\"-\",\"2025-11-05\",\"07:08:36\",\"2025-11-05\",\"07:34:28\"\n"
		ibuf := bytes.NewBuffer(make([]byte, 0))
		ibuf.WriteString(headerLine)
		ibuf.WriteString(line1)
		ibuf.WriteString(line2)
		reader := bufio.NewReader(ibuf)
		scanner := bufio.NewScanner(reader)
		entries, err := app.ParseTogglTimeEntries(scanner)

		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(entries))
			assert.Equal(t, "some work", entries[0].Description)
			assert.Equal(t, "example@email.com", entries[0].Email)
			assert.Equal(t, "test project", entries[0].Project)

			assert.Equal(t, "some work2", entries[1].Description)
			assert.Equal(t, "example2@email.com", entries[1].Email)
			assert.Equal(t, "test project 2", entries[1].Project)
		}
	})

	t.Run("FirstN should consider special width character", func(t *testing.T) {
		value := "charactère à autre"
		result := app.FirstN(value, 14)
		assert.Equal(t, "charactère à", result)
	})
}
