package app_test

import (
	"bufio"
	"bytes"
	"log/slog"
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

	t.Run("parse and transform should return error when missing mapping rule", func(t *testing.T) {
		a := app.New(slog.Default())
		headerLine := "\"Description\",\"Duration\",\"Member\",\"Email\",\"Project\",\"Tags\",\"Start date\",\"Start time\",\"Stop date\",\"Stop time\"\n"
		line1 := "\"some work\",\"0:12:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:00:00\",\"2025-11-05\",\"07:12:00\"\n"
		ibuf := bytes.NewBuffer(make([]byte, 0))
		ibuf.WriteString(headerLine)
		ibuf.WriteString(line1)
		reader := bufio.NewReader(ibuf)
		scanner := bufio.NewScanner(reader)
		entries, err := app.ParseTogglTimeEntries(scanner)

		if !assert.NoError(t, err) {
			return
		}
		config := app.AppConfig{
			Mapping: app.Mapping{
				Projects: map[string]app.MappingRule{},
			},
			InputFile:  "",
			OutputFile: "",
		}

		_, err = a.Transform(&config, entries)

		assert.Error(t, err)
		assert.Equal(t, "failed to find project [test project] from mapping rules", err.Error())

	})

	t.Run("test invalid split length validation", func(t *testing.T) {
		line := "\"some work\",\"0:12:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:00:00\",\"2025-11-05\",07:12:00\"\n"
		err := getErrorFromInvalidLine(line)

		assert.Error(t, err)
		assert.Equal(t, "failed to parse line, expected 10 item split, but found 9 on line 2", err.Error())

	})

	t.Run("test invalid duration validation", func(t *testing.T) {
		line := "\"some work\",\"0:1a:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:00:00\",\"2025-11-05\",\"07:12:00\"\n"
		err := getErrorFromInvalidLine(line)

		assert.Error(t, err)
		assert.Equal(t, "failed to parse duration with value 0:1a:00 on line 2", err.Error())

	})

	t.Run("test start date validation", func(t *testing.T) {
		line := "\"some work\",\"0:12:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-0z\",\"07:00:00\",\"2025-11-05\",\"07:12:00\"\n"
		err := getErrorFromInvalidLine(line)

		assert.Error(t, err)
		assert.Equal(t, "failed to parse start date time with value 2025-11-0z 07:00:00 on line 2", err.Error())

	})

	t.Run("test end date validation", func(t *testing.T) {
		line := "\"some work\",\"0:12:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:00:00\",\"2025-11-05\",\"0z:12:00\"\n"
		err := getErrorFromInvalidLine(line)

		assert.Error(t, err)
		assert.Equal(t, "failed to parse end date time with value 2025-11-05 0z:12:00 on line 2", err.Error())
	})

	t.Run("parse and transform should return expected lines", func(t *testing.T) {
		a := app.New(slog.Default())
		headerLine := "\"Description\",\"Duration\",\"Member\",\"Email\",\"Project\",\"Tags\",\"Start date\",\"Start time\",\"Stop date\",\"Stop time\"\n"
		line1 := "\"some work\",\"0:12:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:00:00\",\"2025-11-05\",\"07:12:00\"\n"
		line2 := "\"some work\",\"0:15:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:12:00\",\"2025-11-05\",\"07:27:00\"\n"
		line3 := "\"some more work\",\"0:44:00\",\"some user\",\"example@email.com\",\"test project\",\"-\",\"2025-11-05\",\"07:27:00\",\"2025-11-05\",\"08:11:00\"\n"
		ibuf := bytes.NewBuffer(make([]byte, 0))
		ibuf.WriteString(headerLine)
		ibuf.WriteString(line1)
		ibuf.WriteString(line2)
		ibuf.WriteString(line3)
		reader := bufio.NewReader(ibuf)
		scanner := bufio.NewScanner(reader)
		entries, err := app.ParseTogglTimeEntries(scanner)

		if !assert.NoError(t, err) {
			return
		}
		config := app.AppConfig{
			Mapping: app.Mapping{
				Projects: map[string]app.MappingRule{},
			},
			InputFile:  "",
			OutputFile: "",
		}

		config.Mapping.Projects["test project"] = app.MappingRule{
			ProjectName: "Test Project",
			ProjectTask: "Some Task",
			Role:        "Senior Dev",
		}
		exportLines, err := a.Transform(&config, entries)

		if !assert.NoError(t, err) {
			return
		}

		assert.Len(t, exportLines, 3)

		assert.Equal(t, app.DynamicsHeaders, exportLines[0])
		assert.Equal(t, "Project Service,Test Project,Some Task,Senior Dev,Work,,,Draft,11/05/2025,30.00,\"some work\",\"some work\"\n", exportLines[1])
		assert.Equal(t, "Project Service,Test Project,Some Task,Senior Dev,Work,,,Draft,11/05/2025,45.00,\"some more work\",\"some more work\"\n", exportLines[2])
	})
}

func getErrorFromInvalidLine(line string) error {
	headerLine := "\"Description\",\"Duration\",\"Member\",\"Email\",\"Project\",\"Tags\",\"Start date\",\"Start time\",\"Stop date\",\"Stop time\"\n"
	line1 := line
	ibuf := bytes.NewBuffer(make([]byte, 0))
	ibuf.WriteString(headerLine)
	ibuf.WriteString(line1)
	reader := bufio.NewReader(ibuf)
	scanner := bufio.NewScanner(reader)
	_, err := app.ParseTogglTimeEntries(scanner)
	return err
}
