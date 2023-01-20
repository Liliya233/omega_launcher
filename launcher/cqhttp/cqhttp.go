package cqhttp

import (
	"bufio"
	"fmt"
	"io"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path"

	"github.com/pterm/pterm"
)

func CQHttpEnablerHelper() {
	// 无法创建目录则退出
	if !utils.MkDir(path.Join(utils.GetCurrentDir(), "cqhttp_storage")) {
		panic("无法创建 cqhttp_storage 目录")
	}
	// 提示与确认
	pterm.Info.Print("现在你可以进行文件上传的操作了, 输入 y 继续配置 go-cqhttp: ")
	utils.GetInputYN()
	// 导入配置
	UnPackCQHttpRunAuth()
	// 配置文件路径
	cqCfgFp := path.Join(utils.GetCurrentDir(), "cqhttp_storage", "config.yml")
	// 如果go-cqhttp配置文件不存在, 则执行初始化操作
	if utils.IsFile(cqCfgFp) {
		pterm.Info.Print("已读取到 go-cqhttp 配置文件, 要使用吗? 使用请输入 y, 需要重新设置请输入 n: ")
		if utils.GetInputYN() {
			// 在使用上次的配置时，将读取cqhttp配置文件的账密，然后对cqhttp配置文件进行更新
			if re := getCQConfig(cqCfgFp); re != nil {
				updateCQConfig(re, "null")
			}
		} else {
			cqhttpInit()
		}
	} else {
		cqhttpInit()
	}
}

func Run(botCfg *defines.LauncherConfig) {
	// 不存在cqhttp目录则退出
	if !utils.IsDir(GetCQHttpDir()) {
		panic("cqhttp_storage 目录不存在, 请使用启动器配置一次 go-cqhttp")
	}
	// 如果不存在cqhttp程序则解压
	if !utils.IsFile(path.Join(GetCqHttpExec())) {
		if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
			pterm.Error.Println("解压 go-cqhttp 时遇到问题")
			panic(err)
		}
	}
	// 读取Omega配置
	utils.MkDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置"))
	qGroupLinkFp, qGroupLinkCfg := getOmegaQGroupLinkConfig()
	qGuildLinkFp, qGuildLinkCfg := getOmegaQGuildLinkConfig()
	// 检查群服互通配置文件的地址是否可用
	if qGroupLinkCfg.Configs.Address == "" || !utils.IsAddressAvailable(qGroupLinkCfg.Configs.Address) {
		port, err := utils.GetAvailablePort()
		if err != nil {
			pterm.Error.Println("无法为 go-cqhttp 获取可用端口")
			panic(err)
		}
		qGroupLinkCfg.Configs.Address = fmt.Sprintf("127.0.0.1:%d", port)
	}
	utils.MkDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置", "群服互通"))
	updateOmegaConfig(qGroupLinkFp, qGroupLinkCfg)
	// 频服互通直接使用群服互通的地址
	qGuildLinkCfg.Configs.Address = qGroupLinkCfg.Configs.Address
	utils.MkDir(path.Join(utils.GetCurrentDataDir(), "omega_storage", "配置", "第三方", "Liliya233", "频服互通"))
	updateOmegaConfig(qGuildLinkFp, qGuildLinkCfg)
	// 启动前, 将Omega配置内的IP地址同步到go-cqhttp配置文件
	updateCQConfig(getCQConfig(path.Join(GetCQHttpDir(), "config.yml")), qGroupLinkCfg.Configs.Address)
	// 配置启动参数
	args := []string{"-faststart"}
	// 配置执行目录
	cmd := exec.Command(GetCqHttpExec(), args...)
	cmd.Dir = path.Join(GetCQHttpDir())
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cqhttp_out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 从管道中获取并打印cqhttp输出内容
	stopOutput := false
	go func() {
		reader := bufio.NewReader(cqhttp_out)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			if stopOutput {
				return
			}
			fmt.Print(readString + "\033[0m")
		}
	}()
	// 启动并持续运行cqhttp
	go func() {
		pterm.Success.Println("正在启动 go-cqhttp")
		err := cmd.Start()
		if err != nil {
			pterm.Error.Println("go-cqhttp 启动时出现错误")
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			pterm.Error.Println("go-cqhttp 运行时出现错误")
			panic(err)
		}
	}()
	// 等待cqhttp启动完成
	WaitConnect(qGroupLinkCfg.Configs.Address)
	// 配置完成后, 根据设置决定是否关闭go-cqhttp输出
	if botCfg.BlockCQHttpOutput {
		pterm.Warning.Println("将屏蔽 go-cqhttp 的输出内容")
		stopOutput = true
	}
	// 打包data文件
	PackCQHttpRunAuth(qGroupLinkFp, qGuildLinkFp)
	// 打印提示消息
	pterm.Success.Println("go-cqhttp 已经启动成功了, 可前往 omega_storage 文件夹对相关组件进行进一步配置")
	pterm.Info.Println("若要为服务器配置 go-cqhttp, 请执行以下的操作：")
	pterm.Info.Println("1. 在服务器使用启动器配置 go-cqhttp, 直至启动器要求进行文件上传操作")
	pterm.Info.Println("2. 将 上传这个文件到云服务器以使用云服务器的群服互通.data 上传至服务器 omega_storage 目录下")
	pterm.Info.Println("3. 在服务器上进行确认, 此时应该配置成功了")
	pterm.Info.Println("如果遇到 go-cqhttp 相关的问题, 可前往 https://docs.go-cqhttp.org/ 寻找可用信息")
}
