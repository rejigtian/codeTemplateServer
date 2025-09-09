package main

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	GitRemote string `json:"gitRemote"` // Git 远程仓库地址
	Port      int    `json:"port"`      // 服务器端口
}

var config Config

func loadConfig() error {
	// 首先尝试从环境变量读取
	if remote := os.Getenv("GIT_REMOTE"); remote != "" {
		config.GitRemote = remote
	}
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
			return nil
		}
	}

	// 如果环境变量不存在，尝试从配置文件读取
	data, err := os.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在时创建默认配置
			config.GitRemote = "git@github.com:rejigtian/codeTemplateFiles.git" // 默认值
			config.Port = 8080                                                  // 默认端口
			return saveConfig()
		}
		return err
	}

	return json.Unmarshal(data, &config)
}

func saveConfig() error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}
