package constant

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"math/rand"
	"sync"
)

//go:embed sentence1-10000.json
var sentenceData []byte

type Sentence struct {
	Name string `json:"name"`
	From string `json:"from"`
}

var (
	lineCount    int
	lineOffsets  []int64
	sentenceOnce sync.Once
)

// initSentences 初始化：统计行数和记录每行偏移
func initSentences() {
	sentenceOnce.Do(func() {
		scanner := bufio.NewScanner(bytes.NewReader(sentenceData))
		var offset int64 = 0
		for scanner.Scan() {
			lineOffsets = append(lineOffsets, offset)
			// scanner.Bytes() returns the line without newline
			offset += int64(len(scanner.Bytes())) + 1 // +1 for newline
			lineCount++
		}
	})
}

// GetRandomSentence 随机获取一条古诗词
func GetRandomSentence() string {
	initSentences()
	if lineCount == 0 {
		return "欢迎使用白虎面板"
	}

	targetIndex := rand.Intn(lineCount)
	start := lineOffsets[targetIndex]

	// 找到行结束位置（通过下一个偏移量或文件末尾）
	var end int64
	if targetIndex < lineCount-1 {
		end = lineOffsets[targetIndex+1] - 1 // -1 移除换行符
	} else {
		end = int64(len(sentenceData))
	}

	// 确保不越界并移除末尾可能的回车符 (Windows \r\n)
	line := sentenceData[start:end]
	line = bytes.TrimRight(line, "\r\n")

	var sData []string
	if err := json.Unmarshal(line, &sData); err == nil && len(sData) >= 1 {
		name := sData[0]
		from := ""
		if len(sData) >= 2 {
			from = sData[1]
		}

		if from != "" {
			return "\"" + name + "\"—— " + from
		}
		return name
	}

	return "欢迎使用白虎面板"
}
