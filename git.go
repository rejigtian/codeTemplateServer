package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitRepo 封装 Git 仓库操作
type GitRepo struct {
	path string
}

// NewGitRepo 创建一个新的 Git 仓库管理器
func NewGitRepo(path string) *GitRepo {
	return &GitRepo{path: path}
}

// Init 初始化或更新 Git 仓库
func (g *GitRepo) Init() error {
	templatesPath := filepath.Join(g.path, "templates")

	// 检查是否已经是 Git 仓库
	if _, err := os.Stat(filepath.Join(templatesPath, ".git")); err == nil {
		// 已经是 Git 仓库，检查是否已设置远程仓库
		if err := g.ensureRemote(); err != nil {
			return err
		}
		// 拉取最新代码
		if err := g.pull(); err != nil {
			return err
		}
		return nil
	}

	// 如果远程仓库已存在，直接克隆
	if config.GitRemote != "" {
		if err := g.clone(); err != nil {
			return err
		}
		return nil
	}

	// 如果没有远程仓库，创建新的本地仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = templatesPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to init git repo: %w", err)
	}

	// 设置远程仓库
	if err := g.ensureRemote(); err != nil {
		return err
	}

	// 创建初始提交
	if err := g.add("."); err != nil {
		return err
	}
	if err := g.commit("Initial commit", false); err != nil {
		return err
	}

	// 设置默认分支为main（如果需要）
	if err := g.setDefaultBranch(); err != nil {
		return err
	}

	return nil
}

// ensureRemote 确保远程仓库已设置
func (g *GitRepo) ensureRemote() error {
	// 检查是否已设置远程仓库
	remoteCmd := exec.Command("git", "remote")
	remoteCmd.Dir = filepath.Join(g.path, "templates")
	output, err := remoteCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check git remote: %w", err)
	}

	hasOrigin := false
	remotes := strings.Split(string(output), "\n")
	for _, remote := range remotes {
		if remote == "origin" {
			hasOrigin = true
			break
		}
	}

	if !hasOrigin {
		// 添加远程仓库
		addCmd := exec.Command("git", "remote", "add", "origin", config.GitRemote)
		addCmd.Dir = filepath.Join(g.path, "templates")
		if err := addCmd.Run(); err != nil {
			return fmt.Errorf("failed to add git remote: %w", err)
		}
	} else {
		// 更新远程仓库地址
		setUrlCmd := exec.Command("git", "remote", "set-url", "origin", config.GitRemote)
		setUrlCmd.Dir = filepath.Join(g.path, "templates")
		if err := setUrlCmd.Run(); err != nil {
			return fmt.Errorf("failed to update git remote: %w", err)
		}
	}

	return nil
}

// setDefaultBranch 设置默认分支为main
func (g *GitRepo) setDefaultBranch() error {
	cmd := exec.Command("git", "branch", "-M", "main")
	cmd.Dir = filepath.Join(g.path, "templates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set default branch: %w", err)
	}
	return nil
}

// add 添加文件到暂存区
func (g *GitRepo) add(path string) error {
	cmd := exec.Command("git", "add", path)
	cmd.Dir = filepath.Join(g.path, "templates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git add: %w", err)
	}
	return nil
}

// clone 克隆远程仓库
func (g *GitRepo) clone() error {
	// 删除已存在的目录
	templatesPath := filepath.Join(g.path, "templates")
	if err := os.RemoveAll(templatesPath); err != nil {
		return fmt.Errorf("failed to remove existing directory: %w", err)
	}

	// 克隆仓库
	cmd := exec.Command("git", "clone", config.GitRemote, "templates")
	cmd.Dir = g.path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// pull 拉取最新代码
func (g *GitRepo) pull() error {
	// 先清理本地更改
	resetCmd := exec.Command("git", "reset", "--hard", "HEAD")
	resetCmd.Dir = filepath.Join(g.path, "templates")
	if err := resetCmd.Run(); err != nil {
		return fmt.Errorf("failed to reset changes: %w", err)
	}

	// 拉取最新代码
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = filepath.Join(g.path, "templates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	return nil
}

// commit 提交更改
func (g *GitRepo) commit(message string, autoPush bool) error {
	// 检查是否有更改需要提交
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = filepath.Join(g.path, "templates")
	output, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}
	if len(output) == 0 {
		return nil // 没有更改，不需要提交
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = filepath.Join(g.path, "templates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git commit: %w", err)
	}

	if autoPush {
		// 先拉取最新代码
		if err := g.pull(); err != nil {
			return fmt.Errorf("failed to pull before push: %w", err)
		}

		// 推送到远程仓库
		pushCmd := exec.Command("git", "push", "origin", "main")
		pushCmd.Dir = filepath.Join(g.path, "templates")
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("failed to git push: %w", err)
		}
	}

	return nil
}

// AddTemplate 添加模板文件并提交
func (g *GitRepo) AddTemplate(templateType, fileName string) error {
	if err := g.add(filepath.Join(templateType, fileName)); err != nil {
		return err
	}
	if err := g.add("metadata.json"); err != nil {
		return err
	}
	message := fmt.Sprintf("Add %s template: %s", templateType, fileName)
	return g.commit(message, true) // 自动推送到远程仓库
}

// DeleteTemplate 删除模板文件并提交
func (g *GitRepo) DeleteTemplate(templateType, fileName string) error {
	if err := g.add(filepath.Join(templateType, fileName)); err != nil {
		return err
	}
	if err := g.add("metadata.json"); err != nil {
		return err
	}
	message := fmt.Sprintf("Delete %s template: %s", templateType, fileName)
	return g.commit(message, true) // 自动推送到远程仓库
}
