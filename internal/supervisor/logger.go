package supervisor

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LogEntry struct {
	Tag       string
	Message   string
	Level     string
	Timestamp time.Time
}

type LogStore struct {
	mu      sync.Mutex
	entries []LogEntry
	maxSize int
}

func NewLogStore(maxSize int) *LogStore {
	return &LogStore{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (s *LogStore) Add(entry LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, entry)
}

type TaggedLogger struct {
	tag     string
	color   *color.Color
	store   *LogStore
	logFile *os.File
}

func NewTaggedLogger(tag string, c *color.Color, store *LogStore, logFile *os.File) *TaggedLogger {
	return &TaggedLogger{
		tag:     tag,
		color:   c,
		store:   store,
		logFile: logFile,
	}
}

func (l *TaggedLogger) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		level := detectLevel(line)

		entry := LogEntry{
			Tag:       l.tag,
			Message:   line,
			Level:     level,
			Timestamp: time.Now(),
		}

		l.store.Add(entry)

		prefix := l.color.Sprintf("[%-6s]", l.tag)
		fmt.Fprintf(os.Stdout, "%s %s\n", prefix, line)

		if l.logFile != nil {
			fmt.Fprintf(l.logFile, "%s [%s] %s\n",
				entry.Timestamp.Format("2006/01/02 15:04:05"),
				l.tag,
				line,
			)
		}
	}

	return len(p), nil
}

func detectLevel(line string) string {
	lower := strings.ToLower(line)
	if strings.Contains(lower, "error") || strings.Contains(lower, "fatal") || strings.Contains(lower, "failed") {
		return "error"
	}
	if strings.Contains(lower, "warn") || strings.Contains(lower, "warning") {
		return "warn"
	}
	return "info"
}

func ensureLogDir() (*os.File, error) {
	err := os.MkdirAll(".uca/logs", 0755)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(".uca/logs/dev.log"); err == nil {
		os.Rename(".uca/logs/dev.log", ".uca/logs/dev.log.1")
	}

	return os.Create(".uca/logs/dev.log")
}
