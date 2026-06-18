package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DailyRotateWriter struct {
	mu       sync.Mutex
	basePath string
	file     *os.File
	curDate  string
}

func NewDailyRotateWriter(basePath string) (*DailyRotateWriter, error) {
	w := &DailyRotateWriter{basePath: basePath}
	return w, w.rotateIfNeeded()
}

func (w *DailyRotateWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *DailyRotateWriter) rotateIfNeeded() error {
	date := time.Now().Format("2006-01-02")
	if date == w.curDate && w.file != nil {
		return nil
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	filename := fmt.Sprintf("%s-%s.log", w.basePath, date)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.curDate = date
	return nil
}

func CleanupOldLogs(basePath string, retentionDays int) {
	if retentionDays <= 0 {
		return
	}
	matches, err := filepath.Glob(basePath + "-*.log")
	if err != nil {
		return
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	prefix := basePath + "-"
	for _, filePath := range matches {
		name := strings.TrimSuffix(strings.TrimPrefix(filePath, prefix), ".log")
		day, err := time.Parse("2006-01-02", name)
		if err != nil {
			continue
		}
		if day.Before(cutoff) {
			_ = os.Remove(filePath)
		}
	}
}
