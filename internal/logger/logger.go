package logger

import (
	"encoding/json"
	"opus-api/internal/types"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// CleanupOldLogs deletes all existing logs
func CleanupOldLogs() {
	os.RemoveAll(types.LogDir)
	os.MkdirAll(types.LogDir, 0755)
}

// RotateLogs keeps only 5 newest log folders
func RotateLogs() {
	if !types.DebugMode {
		return
	}
	entries, err := os.ReadDir(types.LogDir)
	if err != nil || len(entries) < 5 {
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	for i := 0; i < len(entries)-4; i++ {
		os.RemoveAll(filepath.Join(types.LogDir, entries[i].Name()))
	}
}

// CreateLogFolder creates a log folder and returns its path
func CreateLogFolder(requestID string) (string, error) {
	if !types.DebugMode {
		return "", nil
	}
	folderName := time.Now().Format("2006-01-02T15-04-05") + "_" + requestID
	logFolder := filepath.Join(types.LogDir, folderName)
	if err := os.MkdirAll(logFolder, 0755); err != nil {
		return "", err
	}
	return logFolder, nil
}

// WriteJSONLog writes JSON log file
func WriteJSONLog(logFolder, fileName string, data interface{}) {
	if !types.DebugMode || logFolder == "" {
		return
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filepath.Join(logFolder, fileName), jsonBytes, 0644)
}

// WriteTextLog writes text log file
func WriteTextLog(logFolder, fileName, content string) {
	if !types.DebugMode || logFolder == "" {
		return
	}
	os.WriteFile(filepath.Join(logFolder, fileName), []byte(content), 0644)
}

// AppendLog appends content to a log file
func AppendLog(logFolder string, fileName string, content string) error {
	if !types.DebugMode || logFolder == "" {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(logFolder, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}