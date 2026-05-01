package logs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const logPath = ".uca/logs/dev.log"
const prevLogPath = ".uca/logs/dev.log.1"

func View(source string, level string, lines int) error {
	entries, err := readLogFile(logPath)
	if err != nil {
		return fmt.Errorf("no logs found — run 'uca dev' first")
	}

	filtered := filterEntries(entries, source, level)

	start := 0
	if len(filtered) > lines {
		start = len(filtered) - lines
	}

	for _, entry := range filtered[start:] {
		printEntry(entry)
	}

	return nil
}

func Tail(source string, level string) error {
	f, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("no logs found — run 'uca dev' first")
	}
	defer f.Close()

	f.Seek(0, io.SeekEnd)

	fmt.Println("Tailing logs... (Ctrl+C to stop)")

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		entry := parseLine(line)
		if matchesFilter(entry, source, level) {
			printEntry(entry)
		}
	}
}

func Clear() error {
	os.Remove(logPath)
	os.Remove(prevLogPath)
	fmt.Println("Logs cleared")
	return nil
}

type logEntry struct {
	timestamp string
	tag       string
	message   string
}

func readLogFile(path string) ([]logEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []logEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := parseLine(scanner.Text())
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

func parseLine(line string) logEntry {
	line = strings.TrimSpace(line)

	entry := logEntry{message: line}

	if len(line) > 22 && line[19] == ' ' && line[20] == '[' {
		closeBracket := strings.Index(line[20:], "]")
		if closeBracket > 0 {
			entry.timestamp = line[:19]
			entry.tag = line[21 : 20+closeBracket]
			entry.message = strings.TrimSpace(line[20+closeBracket+1:])
		}
	}

	return entry
}

func filterEntries(entries []logEntry, source string, level string) []logEntry {
	if source == "" && level == "" {
		return entries
	}

	var filtered []logEntry
	for _, e := range entries {
		if matchesFilter(e, source, level) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func matchesFilter(entry logEntry, source string, level string) bool {
	if source != "" && entry.tag != source {
		return false
	}
	if level != "" {
		entryLevel := detectLogLevel(entry.message)
		if entryLevel != level {
			return false
		}
	}
	return true
}

func detectLogLevel(msg string) string {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "error") || strings.Contains(lower, "fatal") || strings.Contains(lower, "failed") {
		return "error"
	}
	if strings.Contains(lower, "warn") || strings.Contains(lower, "warning") {
		return "warn"
	}
	return "info"
}

func printEntry(entry logEntry) {
	var tagColor *color.Color
	switch entry.tag {
	case "server":
		tagColor = color.New(color.FgCyan)
	case "vite":
		tagColor = color.New(color.FgYellow)
	case "agent":
		tagColor = color.New(color.FgGreen)
	default:
		tagColor = color.New(color.FgWhite)
	}

	if entry.tag != "" {
		prefix := tagColor.Sprintf("[%-6s]", entry.tag)
		fmt.Printf("%s %s\n", prefix, entry.message)
	} else {
		fmt.Println(entry.message)
	}
}
