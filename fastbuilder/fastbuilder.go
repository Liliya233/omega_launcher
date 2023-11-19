package fastbuilder

import (
	"omega_launcher/launcher"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

// 配置启动参数
func setupCmdArgs(cfg *launcher.Config) []string {
	args := []string{"--no-update-check", "-c", cfg.RentalCode}
	// 是否需要租赁服密码
	if cfg.RentalPasswd != "" {
		args = append(args, "-p")
		args = append(args, cfg.RentalPasswd)
	}
	// 是否启动Omega
	if cfg.StartOmega {
		args = append(args, "-O")
	}
	// 尝试使用 Token
	if IsToken(cfg.FBToken) {
		args = append(args, "--plain-token")
		args = append(args, cfg.FBToken)
	}
	// 验证服务器
	if cfg.AuthServer != "" {
		args = append(args, "-A")
		args = append(args, cfg.AuthServer)
	}
	return args
}

func Run(cfg *launcher.Config) {
	// 启动前保存一次配置
	launcher.SaveConfig(cfg)
	// 读取验证服务器返回的Token并保存
	go func() {
		for {
			if currentToken := loadCurrentFBToken(); IsToken(currentToken) {
				cfg.FBToken = currentToken
				launcher.SaveConfig(cfg)
				break
			}
			time.Sleep(time.Second)
		}
	}()
	// 重启间隔
	restartTime := 0
	for {
		// 记录启动时间
		startTime := time.Now()
		// 启动时提示信息
		pterm.Success.Println("正在启动 Omega/FastBuilder, 请根据其提示进行操作")
		// 给予执行权限
		os.Chmod(getFBExecPath(), 0755)
		// 启动命令
		cmd := exec.Command(getFBExecPath(), setupCmdArgs(cfg)...)
		cmd.Dir = filepath.Join(utils.GetCurrentDataDir())
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		// 启动并持续运行FastBuilder
		if err := cmd.Start(); err != nil {
			pterm.Fatal.WithFatal(false).Println("Omega/FastBuilder 启动时出现错误")
			panic(err)
		}
		// 判断是否正常退出
		if err := cmd.Wait(); err == nil {
			pterm.Success.Println("Omega/FastBuilder 已正常退出, 启动器将结束运行")
			return
		}
		// 为了避免频繁请求, 崩溃后将等待一段时间后重启, 可手动跳过等待
		if restartTime < 3600 {
			restartTime = restartTime + 60
		}
		// 运行时间大于1分钟时视为正常运行, 设置为立即重启
		if time.Since(startTime) > time.Minute {
			restartTime = 0
		}
		// 等待输入或计时结束
		if result, isUser := utils.GetInputYNInTime("是否需要重启 Omega/FastBuilder?", int32(restartTime)); result {
			if isUser {
				restartTime = 0
			}
			continue
		}
		pterm.Success.Println("已选择无需重启 Omega/FastBuilder, 启动器将结束运行")
		break
	}
}
