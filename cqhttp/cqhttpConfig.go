package cqhttp

import (
	"fmt"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/pterm/pterm"
	v2 "gopkg.in/yaml.v2"
)

//go:embed assets/config.yml
var defaultConfigBytes []byte

// Copy from go-cqhttp
// Account 账号配置
type Account struct {
	Uin                    int64  `yaml:"uin"`
	Password               string `yaml:"password"`
	SignServer             string `yaml:"sign-server"`
	IsLowVersionSignServer bool   `yaml:"is-below-110"`
	SignServerKey          string `yaml:"key"`
}

// Config 总配置文件
type CQHttpConfig struct {
	Account *Account `yaml:"account"`
}

// 从cqhttp配置里读取QQ账密信息
func getCQConfig() *CQHttpConfig {
	cfg := &CQHttpConfig{}
	data, err := os.ReadFile(filepath.Join(GetCQHttpDir(), "config.yml"))
	if err != nil {
		return nil
	}
	if err := v2.Unmarshal(data, cfg); err != nil {
		return nil
	}
	return cfg
}

// 写入cqhttp配置
func writeCQConfig(cfgStr string) {
	err := utils.WriteFileData(filepath.Join(GetCQHttpDir(), "config.yml"), []byte(cfgStr))
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("更新 go-cqhttp 配置文件时遇到问题")
		panic(err)
	}
}

// 更新cqhttp配置文件的地址
func updateCQConfigAddress(wsAddress string) {
	cfgStr := strings.Replace(string(defaultConfigBytes), "[WS地址]", wsAddress, 1)
	if cqCfg := getCQConfig(); cqCfg != nil {
		// 保留账密信息
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", fmt.Sprint(cqCfg.Account.Uin), 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", cqCfg.Account.Password, 1)
		cfgStr = strings.Replace(cfgStr, "[SignServer地址]", cqCfg.Account.SignServer, 1)
		cfgStr = strings.Replace(cfgStr, "[IsLowVersionSignServer]", fmt.Sprint(cqCfg.Account.IsLowVersionSignServer), 1)
		cfgStr = strings.Replace(cfgStr, "[SignServerKey]", cqCfg.Account.SignServerKey, 1)
	} else {
		// 默认配置
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", "1233456", 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", "", 1)
		cfgStr = strings.Replace(cfgStr, "[SignServer地址]", "-", 1)
		cfgStr = strings.Replace(cfgStr, "[IsLowVersionSignServer]", "false", 1)
		cfgStr = strings.Replace(cfgStr, "[SignServerKey]", "114514", 1)
	}
	// 写入新配置
	writeCQConfig(cfgStr)
}

// 初始化cqhttp
func initCQConfig() {
	// 移除token等文件
	utils.RemoveFile(filepath.Join(GetCQHttpDir(), "device.json"))
	utils.RemoveFile(filepath.Join(GetCQHttpDir(), "session.token"))
	// 要求输入cqhttp配置信息
	cfgStr := strings.Replace(string(defaultConfigBytes), "[QQ账号]", fmt.Sprint(utils.GetInt64Input("请输入QQ账号")), 1)
	cfgStr = strings.Replace(cfgStr, "[QQ密码]", utils.GetPswInput("请输入QQ密码"), 1)
	cfgStr = strings.Replace(cfgStr, "[WS地址]", "null", 1)
	SignServer := utils.GetInput("请输入 Sign Server 地址 (没有或使用启动器配置请留空)")
	if SignServer == "" {
		SignServer = "-"
	}
	cfgStr = strings.Replace(cfgStr, "[SignServer地址]", SignServer, 1)
	cfgStr = strings.Replace(cfgStr, "[IsLowVersionSignServer]", "false", 1)
	cfgStr = strings.Replace(cfgStr, "[SignServerKey]", "114514", 1)
	// 写入新配置
	writeCQConfig(cfgStr)
}
