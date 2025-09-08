package main

import (
	"encoding/json"
	"os"
	"sync"
	_ "time"
)

type TemplateInfo struct {
	DisplayName string `json:"displayName"`
	FileName    string `json:"fileName"`
	Type        string `json:"type"`
	CreateTime  int64  `json:"createTime"` // 添加创建时间字段
}

type TemplateMetadata struct {
	sync.RWMutex
	Templates map[string]TemplateInfo `json:"templates"` // 使用 map 存储
	FilePath  string                  `json:"-"`
}

func NewTemplateMetadata(filePath string) *TemplateMetadata {
	return &TemplateMetadata{
		Templates: make(map[string]TemplateInfo),
		FilePath:  filePath,
	}
}

func (m *TemplateMetadata) Load() error {
	data, err := os.ReadFile(m.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, m)
}

func (m *TemplateMetadata) Save() error {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.FilePath, data, 0644)
}
