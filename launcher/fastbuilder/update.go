package fastbuilder

import (
	"bytes"
	"io"
	"omega_launcher/defines"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

// Fastbuilder远程仓库地址
var STORAGE_REPO = ""

// 仓库选择
func selectRepo(cfg *defines.LauncherConfig, reselect bool) {
	if reselect || cfg.Repo < 1 || cfg.Repo > 4 {
		// 不再于列表提示自用仓库
		pterm.Info.Printf("当前可选择的仓库有：\n1. Github 仓库\n2. Github 镜像仓库\n")
		pterm.Info.Printf("请输入序号来选择一个仓库: ")
		cfg.Repo = utils.GetIntInputInScope(1, 4)
	}
	if cfg.Repo == 1 {
		pterm.Info.Println("将使用 Github 仓库进行更新")
		STORAGE_REPO = defines.REMOTE_REPO
	} else if cfg.Repo == 2 {
		pterm.Info.Println("将使用 Github 镜像仓库进行更新")
		STORAGE_REPO = defines.MIRROR_REPO
	} else if cfg.Repo == 3 {
		pterm.Info.Println("将使用本地仓库进行更新 (自用)")
		STORAGE_REPO = defines.LOCAL_REPO
	} else if cfg.Repo == 4 {
		pterm.Info.Println("将使用云裳仓库进行更新 (不一定最新)")
		STORAGE_REPO = defines.YSCLOUD_REPO
	} else {
		panic("无效的仓库, 请重新进行选择")
	}
}

// 下载FB
func download() {
	var execBytes []byte
	var err error
	// 获取写入路径与远程仓库url
	path := getFBExecPath()
	url := STORAGE_REPO + GetFBExecName()
	// 下载
	compressedData := utils.DownloadSmallContent(url)
	// 官网并没有提供brotli, 所以对读取操作进行修改
	if execBytes, err = io.ReadAll(bytes.NewReader(compressedData)); err != nil {
		panic(err)
	}
	// 写入文件
	if err := utils.WriteFileData(path, execBytes); err != nil {
		panic(err)
	}
}

// 升级FB
func Update(cfg *defines.LauncherConfig, reselect bool) {
	selectRepo(cfg, reselect)
	pterm.Warning.Println("正在从指定仓库获取更新信息..")
	targetHash := getRemoteFBHash(STORAGE_REPO)
	currentHash := getCurrentFBHash()
	//fmt.Println(targetHash)
	//fmt.Println(currentHash)
	if targetHash == currentHash {
		pterm.Success.Println("太好了, 你的 Fastbuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 Fastbuilder, 请保持耐心..")
		download()
	}
}
