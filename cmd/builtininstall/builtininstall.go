package builtininstall

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/engigu/baihu-panel/internal/logger"
	"github.com/engigu/baihu-panel/internal/utils"
)

// Run 执行内建包安装逻辑
func Run(args []string) {
	fs := flag.NewFlagSet("builtininstall", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n白虎面板内建依赖安装工具\n\n")
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  baihu builtininstall\n\n")
		fmt.Fprintf(os.Stderr, "说明:\n")
		fmt.Fprintf(os.Stderr, "  此命令将为系统中已安装的 Node.js 和 Python 环境安装白虎面板所需的内建依赖包。\n\n")
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	logger.Infof("[Builtin] 开始为 mise 环境安装内建包...")

	// 1. 确定内建包路径
	// 优先使用 /www/builtin (Docker 环境)，否则尝试相对于二进制文件的当前目录
	builtinPath := "/www/builtin"
	if _, err := os.Stat(builtinPath); os.IsNotExist(err) {
		// 回退到当前目录下的 builtin
		pwd, _ := os.Getwd()
		builtinPath = filepath.Join(pwd, "builtin")
	}

	if _, err := os.Stat(builtinPath); os.IsNotExist(err) {
		logger.Errorf("[Builtin] 找不到内建包目录: %s", builtinPath)
		return
	}

	// 2. 安装 Node.js 包
	installForLanguage("node", filepath.Join(builtinPath, "nodejs"))

	// 3. 安装 Python 包
	installForLanguage("python", filepath.Join(builtinPath, "python"))

	logger.Infof("[Builtin] 内建包安装流程完成")
}

func installForLanguage(lang, pkgPath string) {
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		logger.Warnf("[Builtin] %s 的内建包目录不存在: %s", lang, pkgPath)
		return
	}

	versions, err := utils.ListMiseInstalledVersions(lang)
	if err != nil {
		logger.Errorf("[Builtin] 获取 %s 的 mise 版本列表失败: %v", lang, err)
		return
	}

	if len(versions) == 0 {
		logger.Infof("[Builtin] 未发现已安装的 %s 版本，跳过", lang)
		return
	}

	for _, v := range versions {
		logger.Infof("[Builtin] 正在为 %s@%s 安装内建包...", lang, v)
		
		var subCmdArgs []string
		if lang == "node" {
			// 使用 npm i -g 进行全局安装
			subCmdArgs = []string{"npm", "i", "-g", pkgPath}
		} else {
			// python 改为标准安装 (非 -e)，避免 Docker 内软链接可能导致的路径丢失问题
			subCmdArgs = []string{"pip", "install", "--force-reinstall", pkgPath}
		}

		// 构建参数列表: [mise, exec, lang@v, --, cmd...]
		fullArgs := utils.BuildMiseCommandArgsSimple(subCmdArgs, lang, v)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", append([]string{"/c"}, fullArgs...)...)
		} else {
			cmd = exec.Command(fullArgs[0], fullArgs[1:]...)
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Errorf("[Builtin] 为 %s@%s 安装失败: %v\n输出: %s", lang, v, err, string(out))
		} else {
			logger.Infof("[Builtin] 为 %s@%s 安装成功", lang, v)
		}
	}
}
