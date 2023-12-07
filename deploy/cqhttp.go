package deploy

import (
	"bytes"
	"encoding/json"
	"io"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"path/filepath"

	"github.com/andybalholm/brotli"
	"github.com/pterm/pterm"
)

var (
	cqhttpDownloadUrl = "https://github.com/Liliya233/go-cqhttp/releases/latest/download/"
)

func isCQHttpCache() bool {
	// brotli名称
	brotliName := plantform.GetCQHttpName() + ".brotli"
	// 获取文件内容
	jsonData, err := utils.DownloadBytesWithMirror(cqhttpDownloadUrl + "hashes.json")
	if err != nil {
		panic(err)
	}
	// 解析文件内容
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	// 远程brotli文件hash
	remoteHash := hashMap[brotliName]
	if remoteHash == "" {
		panic("未能从远程仓库获取 Hash")
	}
	// 对比本地brotli文件hash, 并返回
	return remoteHash == utils.GetFileHash(filepath.Join(utils.GetCacheDir(), "downloads", brotliName))
}

func CQHttpDeploy() {
	pterm.Warning.Printfln("正在检查 go-cqhttp 更新..")
	// brotli名称
	brotliName := plantform.GetCQHttpName() + ".brotli"
	// brotli路径
	brotliDir := filepath.Join(utils.GetCacheDir(), "downloads", brotliName)
	// 检查缓存文件
	if !isCQHttpCache() {
		pterm.Warning.Println("正在为你下载最新的 go-cqhttp, 请保持耐心..")
		err := utils.DownloadFileWithMirror(cqhttpDownloadUrl+brotliName, brotliDir)
		if err != nil {
			panic(err)
		}
	}
	var brotliBytes, execBytes []byte
	var err error
	// 读取brotli
	if brotliBytes, err = utils.GetFileData(brotliDir); err != nil {
		panic(err)
	}
	if execBytes, err = io.ReadAll(brotli.NewReader(bytes.NewReader(brotliBytes))); err != nil {
		panic(err)
	}
	// 写入exe
	if err = utils.WriteFileData(filepath.Join(utils.GetCurrentDir(), plantform.GetCQHttpName()), execBytes); err != nil {
		panic(err)
	}
}
