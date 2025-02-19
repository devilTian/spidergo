/*
 * @Author: tianye@shimiotech.cn
 * @Date: 2025-02-19 11:28:17
 */
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)
type assert_data struct {
	Content string
}

type clip struct {
	AssetInfo assert_data
	InPoint int
	OutPoint int
}

type track_data struct {
	Clips []clip
}

type full_data struct {
	Tracks []track_data
}

const StrFormat = `%d                                                                                                                      
%s --> %s
%s

`
func formatMilliseconds(ms int) string {
	// 将毫秒转换为 time.Duration
	duration := time.Duration(ms) * time.Millisecond

	// 提取小时、分钟、秒和毫秒
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	duration -= minutes * time.Minute
	seconds := duration / time.Second
	milliseconds := duration % time.Second / time.Millisecond

	// 格式化为 00:00:00,000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

func main() {
	// go run shell/genSubtitle.go >> /Users/tianye/Documents/_RunVlog/生成的字幕.srt
	content, err := os.ReadFile("/Users/tianye/Movies/Bcut Drafts/4EF8AB24-4D24-4CB6-87F5-3748F0A82840/11-05-17-162--{8f997659-5973-49b6-a588-9ad9b3057d48}.json")
	if err != nil {
		fmt.Println("文件无法读取")
		return
	}
	data := full_data{}
	json.Unmarshal(content, &data)
	clips := data.Tracks[0].Clips
	for i, v := range clips {
		fmt.Printf(StrFormat, i, formatMilliseconds(v.InPoint), formatMilliseconds(v.OutPoint), v.AssetInfo.Content)
	}
}