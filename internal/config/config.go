package config

import (
	"errors"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var ConfigDir string

// InitConfigDir finds the configuration directory for micro according to the XDG spec.
// If no directory is found, it creates one.
func InitConfigDir(flagConfigDir string) error {
	var e error

	microHome := os.Getenv("MICRO_CONFIG_HOME")
	if microHome == "" {
		// The user has not set $MICRO_CONFIG_HOME so we'll try $XDG_CONFIG_HOME
		xdgHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgHome == "" {
			// The user has not set $XDG_CONFIG_HOME so we should act like it was set to ~/.config
			home, err := homedir.Dir()
			if err != nil {
				return errors.New("查找主目录时出错\n无法加载配置文件: " + err.Error())
			}
			xdgHome = filepath.Join(home, ".config")
		}

		microHome = filepath.Join(xdgHome, "micro")
	}
	ConfigDir = microHome

	if len(flagConfigDir) > 0 {
		if _, err := os.Stat(flagConfigDir); os.IsNotExist(err) {
			e = errors.New("错误: " + flagConfigDir + " 不存在。默认为 " + ConfigDir + ".")
		} else {
			ConfigDir = flagConfigDir
			return nil
		}
	}

	// Create micro config home directory if it does not exist
	// This creates parent directories and does nothing if it already exists
	err := os.MkdirAll(ConfigDir, os.ModePerm)
	if err != nil {
		return errors.New("创建配置目录时出错: " + err.Error())
	}

	return e
}
