package app

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type App struct {
	logger *slog.Logger
}

type TogglTimeEntry struct {
	description   string
	duration      time.Duration
	member        string
	email         string
	project       string
	tags          string
	startDateTime time.Time
	stopDateTime  time.Time
}

func NewTogglTimeEntry(
	description string,
	duration time.Duration,
	member,
	email,
	project,
	tags string,
	startDateTime time.Time,
	stopDateTime time.Time,
) TogglTimeEntry {
	return TogglTimeEntry{
		description,
		duration,
		member,
		email,
		project,
		tags,
		startDateTime,
		stopDateTime,
	}
}

func (e TogglTimeEntry) GetGroupingKey() string {
	return e.startDateTime.Format("2006-01-02") + e.project + e.description
}

func (a *App) LoadTogglCsvExportLines(filePath string) ([]TogglTimeEntry, error) {
	file, err := os.Open(filePath) // Replace "example.txt" with your file path

	if err != nil {
		return nil, err
	}

	defer file.Close() // Ensure the file is closed when the function exits

	// Create a new scanner to read lines
	scanner := bufio.NewScanner(file)

	entries := make([]TogglTimeEntry, 0)
	// Iterate over each line
	lineCount := 0

	for scanner.Scan() {

		lineCount++

		if lineCount == 1 {
			continue
		}

		line := scanner.Text() // Get the current line as a string
		line = firstN(line, len(line)-1)
		lineSplit := strings.Split(line, "\",\"")

		duration, err := parseDuration(lineSplit[1])

		if err != nil {
			return nil, fmt.Errorf("failed to parse duration with value %s on line %v", lineSplit[1], lineCount)
		}

		startDateTimeString := fmt.Sprintf("%s %s", lineSplit[6], lineSplit[7])
		startDateTime, err := time.Parse("2006-01-02 15:04:05", startDateTimeString)

		if err != nil {
			return nil, fmt.Errorf("failed to parse start date time with value %s on line %v", startDateTimeString, lineCount)
		}

		stopDateTimeString := fmt.Sprintf("%s %s", lineSplit[8], lineSplit[9])
		stopDateTime, err := time.Parse("2006-01-02 15:04:05", stopDateTimeString)

		if err != nil {
			return nil, fmt.Errorf("failed to parse end date time with value %s on line %v", stopDateTimeString, lineCount)
		}

		entry := NewTogglTimeEntry(
			firstN(lineSplit[0], 100),
			duration,
			lineSplit[2],
			lineSplit[3],
			lineSplit[4],
			lineSplit[5],
			startDateTime,
			stopDateTime,
		)

		entries = append(entries, entry)
	}

	return entries, nil
}

func firstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}

func parseDuration(d string) (time.Duration, error) {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(d, "%d:%d:%d", &hours, &minutes, &seconds)
	if err != nil {
		return time.Duration(0), err
	}

	// Convert the parsed components into a time.Duration
	duration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second

	return duration, nil
}

func New(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}
