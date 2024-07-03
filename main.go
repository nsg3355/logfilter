package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LogEntry struct {
	Log    string `json:"log"`
	Stream string `json:"stream"`
	Time   string `json:"time"`
	File   string `json:"file"` // 파일명을 저장하기 위한 필드
}

func filterLogs(logFile, keyword string, errors *[]LogEntry) error {
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var logEntry LogEntry
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(logEntry.Log), strings.ToLower(keyword)) {
			logEntry.File = logFile // 로그 항목에 파일명을 추가
			*errors = append(*errors, logEntry)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	logDir := "./logs"                    // 로그 파일들이 있는 디렉토리 경로
	outputFilePath := "./error_logs.json" // 에러 로그를 저장할 파일 경로
	keyword := "panic"
	var errors []LogEntry

	// 디렉토리 내 모든 로그 파일을 읽음
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "-json.log") {
			if err := filterLogs(path, keyword, &errors); err != nil {
				fmt.Println("Error:", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through log directory:", err)
		return
	}

	// 에러 로그를 JSON 파일로 저장
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(errors); err != nil {
		fmt.Println("Error encoding JSON:", err)
	} else {
		fmt.Println("Error logs saved to", outputFilePath)
	}
}
