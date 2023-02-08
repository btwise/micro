package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/config"
)

func shouldContinue() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("继续 [Y/n]: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return false
	}

	text = strings.TrimRight(text, "\r\n")

	return len(text) == 0 || strings.ToLower(text)[0] == 'y'
}

// CleanConfig performs cleanup in the user's configuration directory
func CleanConfig() {
	fmt.Println("在以下位置清理您的配置目录", config.ConfigDir)
	fmt.Printf("请先考虑备份%s,然后再继续\n", config.ConfigDir)

	if !shouldContinue() {
		fmt.Println("尽早停止")
		return
	}

	fmt.Println("清洁默认设置")
	config.WriteSettings(filepath.Join(config.ConfigDir, "settings.json"))

	// detect unused options
	var unusedOptions []string
	defaultSettings := config.DefaultAllSettings()
	for k := range config.GlobalSettings {
		if _, ok := defaultSettings[k]; !ok {
			valid := false
			for _, p := range config.Plugins {
				if strings.HasPrefix(k, p.Name+".") || k == p.Name {
					valid = true
				}
			}
			if !valid {
				unusedOptions = append(unusedOptions, k)
			}
		}
	}

	if len(unusedOptions) > 0 {
		fmt.Println("以下选项未使用:")

		sort.Strings(unusedOptions)

		for _, s := range unusedOptions {
			fmt.Printf("%s (值: %v)\n", s, config.GlobalSettings[s])
		}

		fmt.Printf("这些选项将从%s中删除\n", filepath.Join(config.ConfigDir, "settings.json"))

		if shouldContinue() {
			for _, s := range unusedOptions {
				delete(config.GlobalSettings, s)
			}

			err := config.OverwriteSettings(filepath.Join(config.ConfigDir, "settings.json"))
			if err != nil {
				fmt.Println("写入settings.json文件时出错: " + err.Error())
			}

			fmt.Println("删除了未使用的选项")
			fmt.Print("\n\n")
		}
	}

	// detect incorrectly formatted buffer/ files
	files, err := ioutil.ReadDir(filepath.Join(config.ConfigDir, "buffers"))
	if err == nil {
		var badFiles []string
		var buffer buffer.SerializedBuffer
		for _, f := range files {
			fname := filepath.Join(config.ConfigDir, "buffers", f.Name())
			file, e := os.Open(fname)

			if e == nil {
				decoder := gob.NewDecoder(file)
				err = decoder.Decode(&buffer)

				if err != nil && f.Name() != "history" {
					badFiles = append(badFiles, fname)
				}
				file.Close()
			}
		}

		if len(badFiles) > 0 {
			fmt.Printf("在%s中检测到格式无效的%d个文件\n", len(badFiles), filepath.Join(config.ConfigDir, "buffers"))
			fmt.Println("这些文件存储光标和撤消历史记录.")
			fmt.Printf("删除%s中格式错误的文件\n", filepath.Join(config.ConfigDir, "buffers"))

			if shouldContinue() {
				removed := 0
				for _, f := range badFiles {
					err := os.Remove(f)
					if err != nil {
						fmt.Println(err)
						continue
					}
					removed++
				}

				if removed == 0 {
					fmt.Println("无法删除文件")
				} else {
					fmt.Printf("删除了%d个格式错误的文件\n", removed)
				}
				fmt.Print("\n\n")
			}
		}
	}

	// detect plugins/ directory
	plugins := filepath.Join(config.ConfigDir, "plugins")
	if stat, err := os.Stat(plugins); err == nil && stat.IsDir() {
		fmt.Printf("找到目录 %s\n", plugins)
		fmt.Printf("插件现在应该存储在%s中\n", filepath.Join(config.ConfigDir, "plug"))
		fmt.Printf("删除%s\n", plugins)

		if shouldContinue() {
			os.RemoveAll(plugins)
		}

		fmt.Print("\n\n")
	}

	fmt.Println("完成清理")
}
