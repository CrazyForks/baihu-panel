package constant

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"os"
	"sync"
)

type Sentence struct {
	Name string `json:"name"`
	From string `json:"from"`
}

const sentenceFile = "internal/constant/sentence1-10000.json"

var (
	lineCount     int
	lineCountOnce sync.Once
)

// countLines 统计文件行数（只执行一次）
func countLines() {
	lineCountOnce.Do(func() {
		file, err := os.Open(sentenceFile)
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineCount++
		}
	})
}

// GetRandomSentence 随机获取一条古诗词
func GetRandomSentence() string {
	countLines()
	if lineCount == 0 {
		return "欢迎使用白虎面板"
	}

	targetLine := rand.Intn(lineCount)

	file, err := os.Open(sentenceFile)
	if err != nil {
		return "欢迎使用白虎面板"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	for scanner.Scan() {
		if currentLine == targetLine {
			var s Sentence
			if err := json.Unmarshal(scanner.Bytes(), &s); err == nil && s.Name != "" {
				return s.Name
			}
			break
		}
		currentLine++
	}

	return "欢迎使用白虎面板"
}
