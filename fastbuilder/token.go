package fastbuilder

import (
	"omega_launcher/launcher"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"
)

// 判断是否为token
func IsToken(token string) bool {
	return strings.HasPrefix(token, "w9/") || strings.HasPrefix(token, "5YQ")
}

// 加载现有的token
func loadCurrentFBToken() string {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	tokenPath := filepath.Join(homedir, ".config", "fastbuilder", "fbtoken")
	// 尝试读取token文件
	if utils.IsFile(tokenPath) {
		if data, err := utils.GetFileData(tokenPath); err == nil {
			return string(data)
		}
	}
	return ""
}

// 删除现有的token
func deleteCurrentFBToken() bool {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	tokenPath := filepath.Join(homedir, ".config", "fastbuilder", "fbtoken")
	// 尝试删除token文件
	return utils.RemoveFile(tokenPath)
}

// 配置Token
func FBTokenSetup(cfg *launcher.Config) {
	// 配置文件
	if IsToken(cfg.FBToken) && utils.GetInputYN("要使用配置文件的 FBToken 吗?") {
		return
	}
	// 设备文件
	if currentFbToken := loadCurrentFBToken(); IsToken(currentFbToken) {
		if utils.GetInputYN("要使用当前设备储存的 FBToken 吗?") {
			cfg.FBToken = currentFbToken
			return
		} else {
			deleteCurrentFBToken()
		}
	}
	// 用户输入
	cfg.FBToken = utils.GetInput("请输入 FBToken (没有或账密登录请留空)")
}
