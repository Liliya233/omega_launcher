package launcher

import (
	"omega_launcher/utils"
	"path/filepath"
)

// 启动器仓库配置文件结构
type ConfigRepo struct {
	Name         string `json:"名称"`
	Url          string `json:"URL"`
	IsPreRelease bool   `json:"是否预构建"`
	UseMirror    bool   `json:"是否使用镜像"`
}

// 启动器配置文件结构
type Config struct {
	Repo              *ConfigRepo `json:"仓库信息"`
	RentalCode        string      `json:"租赁服号"`
	RentalPasswd      string      `json:"租赁服密码"`
	FBToken           string      `json:"FBToken"`
	AuthServer        string      `json:"验证服务器"`
	EnableCQHttp      bool        `json:"是否开启go-cqhttp"`
	StartOmega        bool        `json:"是否启动Omega"`
	UpdateFB          bool        `json:"是否更新FB"`
	BlockCQHttpOutput bool        `json:"是否屏蔽go-cqhttp的内容输出"`
}

// 保存配置文件
func SaveConfig(config *Config) {
	utils.WriteJsonData(filepath.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), config)
}
