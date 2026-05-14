package task

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/engigu/baihu-panel/cmd/clibase"
	"github.com/engigu/baihu-panel/internal/constant"
	"github.com/engigu/baihu-panel/internal/database"
	"github.com/engigu/baihu-panel/internal/models"
	"github.com/engigu/baihu-panel/internal/utils"
)

// 打印主帮助
func printMainHelp() {
	fmt.Fprintf(os.Stderr, "\n白虎面板任务命令行管理工具 (Task CLI)\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  baihu task <子命令> [参数]\n\n")
	fmt.Fprintf(os.Stderr, "可用子命令:\n")
	fmt.Fprintf(os.Stderr, "  list       查询并输出任务列表\n")
	fmt.Fprintf(os.Stderr, "  run        手动立即触发执行指定的任务\n")
	fmt.Fprintf(os.Stderr, "  enable     启用指定的任务（同步加入后台调度队列）\n")
	fmt.Fprintf(os.Stderr, "  disable    禁用指定的任务（同步从后台调度队列摘除）\n")
	fmt.Fprintf(os.Stderr, "  status     查看指定任务最近一次执行的完整输出与状态\n")
	fmt.Fprintf(os.Stderr, "  history    查看指定任务近期的多次执行流水记录\n\n")
	fmt.Fprintf(os.Stderr, "使用 'baihu task <子命令> --help' 查看具体子命令的参数说明和示例。\n\n")
}

// Run 任务命令行入口
func Run(args []string) {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printMainHelp()
		return
	}

	subCommand := args[0]
	subArgs := args[1:]

	switch subCommand {
	case "list":
		runList(subArgs)
	case "run":
		runExecute(subArgs)
	case "enable", "disable":
		runToggle(subCommand, subArgs)
	case "status":
		runStatus(subArgs)
	case "history":
		runHistory(subArgs)
	default:
		fmt.Fprintf(os.Stderr, "未知子命令: %s\n", subCommand)
		printMainHelp()
	}
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	namePtr := fs.String("name", "", "按任务名称或备注进行模糊筛选")
	typePtr := fs.String("type", "", "按任务类型筛选 (例如: task, repo)")
	pagePtr := fs.Int("page", 1, "查询页码")
	sizePtr := fs.Int("size", 20, "每页展示条数")

	fs.Usage = func() {
		clibase.PrintSubCommandUsage("白虎面板任务列表查询工具", "baihu task list [参数]", "  baihu task list\n  baihu task list -page 2 -size 10\n  baihu task list -name \"签到\"", fs)
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	clibase.InitContext(false)

	var total int64
	query := database.DB.Model(&models.Task{})
	if *namePtr != "" {
		query = query.Where("name LIKE ? OR remark LIKE ?", "%"+*namePtr+"%", "%"+*namePtr+"%")
	}
	if *typePtr != "" {
		query = query.Where("type = ?", *typePtr)
	}
	query.Count(&total)

	offset := (*pagePtr - 1) * *sizePtr
	if offset < 0 {
		offset = 0
	}

	var tasks []models.Task
	query.Order("created_at DESC").Limit(*sizePtr).Offset(offset).Find(&tasks)

	fmt.Println("====================================================================================================")
	fmt.Printf("%-12s | %-25s | %-18s | %-8s | %-6s\n", "任务ID", "任务名称", "Cron规则", "类型", "状态")
	fmt.Println("----------------------------------------------------------------------------------------------------")
	for _, t := range tasks {
		name := t.Name
		// 截断过长名称以对齐
		if len(name) > 22 {
			name = string([]rune(name)[:20]) + ".."
		}
		cron := t.Schedule
		if cron == "" {
			cron = "-"
		}
		status := "启用"
		if !utils.DerefBool(t.Enabled, true) {
			status = "禁用"
		}
		fmt.Printf("%-12s | %-25s | %-18s | %-8s | %-6s\n", t.ID, name, cron, t.Type, status)
	}
	fmt.Println("====================================================================================================")
	totalPages := (total + int64(*sizePtr) - 1) / int64(*sizePtr)
	if totalPages == 0 {
		totalPages = 1
	}
	fmt.Printf("共查询到 %d 个任务记录，当前展示第 %d/%d 页 (每页 %d 条)。\n", total, *pagePtr, totalPages, *sizePtr)
	fmt.Printf("提示: 追加参数 (例如 '-page 2 -size 50') 即可灵活查看指定页码或调整展示数量。\n")
}

func runExecute(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	fs.Usage = func() {
		clibase.PrintSubCommandUsage("白虎面板手动任务触发工具", "baihu task run <任务ID>", "  baihu task run a1b2c3d4", nil)
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	parsedArgs := fs.Args()
	if len(parsedArgs) < 1 {
		fmt.Fprintf(os.Stderr, "错误: 缺少目标任务ID。\n")
		fs.Usage()
		return
	}
	taskID := parsedArgs[0]

	clibase.InitContext(false)

	_, err := clibase.CallInternalAPI("POST", "/internal/tasks/execute/"+taskID, map[string]interface{}{})
	if err != nil {
		fmt.Printf(">> 任务触发失败: %v\n", err)
		return
	}

	fmt.Printf(">> 任务 [%s] 触发指令下发成功！已进入后台调度队列排队或执行。\n", taskID)
	fmt.Printf(">> 提示: 可以使用 'baihu task status %s' 查看近期执行输出。\n", taskID)
}

func runToggle(action string, args []string) {
	fs := flag.NewFlagSet(action, flag.ExitOnError)
	actionName := "启用"
	targetEnabled := true
	if action == "disable" {
		actionName = "禁用"
		targetEnabled = false
	}

	fs.Usage = func() {
		clibase.PrintSubCommandUsage(fmt.Sprintf("白虎面板任务%s工具", actionName), fmt.Sprintf("baihu task %s <任务ID>", action), fmt.Sprintf("  baihu task %s a1b2c3d4", action), nil)
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	parsedArgs := fs.Args()
	if len(parsedArgs) < 1 {
		fmt.Fprintf(os.Stderr, "错误: 缺少目标任务ID。\n")
		fs.Usage()
		return
	}
	taskID := parsedArgs[0]

	clibase.InitContext(false)

	_, err := clibase.CallInternalAPI("POST", "/internal/tasks/toggle/"+taskID, map[string]interface{}{
		"enabled": targetEnabled,
	})
	if err != nil {
		fmt.Printf(">> 切换状态操作失败: %v\n", err)
		return
	}

	fmt.Printf(">> 任务 [%s] 已成功%s！\n", taskID, actionName)
}

func runStatus(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.Usage = func() {
		clibase.PrintSubCommandUsage("白虎面板任务执行状态与日志查看工具", "baihu task status <任务ID> [日志ID]", "  baihu task status a1b2c3d4\n  baihu task status a1b2c3d4 log_123456", nil)
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	parsedArgs := fs.Args()
	if len(parsedArgs) < 1 {
		fmt.Fprintf(os.Stderr, "错误: 缺少目标任务ID。\n")
		fs.Usage()
		return
	}
	taskID := parsedArgs[0]
	var specificLogID string
	if len(parsedArgs) > 1 {
		specificLogID = parsedArgs[1]
	}

	clibase.InitContext(false)

	var taskLog models.TaskLog
	query := database.DB.Where("task_id = ?", taskID)
	if specificLogID != "" {
		query = query.Where("id = ?", specificLogID)
	}
	res := query.Order("created_at DESC").Limit(1).Find(&taskLog)
	if res.Error != nil || res.RowsAffected == 0 {
		if specificLogID != "" {
			fmt.Printf("找不到任务 [%s] 指定日志ID [%s] 的记录。\n", taskID, specificLogID)
		} else {
			fmt.Printf("找不到任务 [%s] 的任何执行记录。\n", taskID)
		}
		return
	}

	var task models.Task
	database.DB.Where("id = ?", taskID).Limit(1).Find(&task)
	taskName := taskID
	if task.Name != "" {
		taskName = task.Name
	}

	statusText := "运行中"
	switch taskLog.Status {
	case constant.TaskStatusSuccess:
		statusText = "成功"
	case constant.TaskStatusFailed:
		statusText = "失败"
	case constant.TaskStatusTimeout:
		statusText = "超时"
	case constant.TaskStatusCancelled:
		statusText = "已取消"
	}

	fmt.Println("====================================================================================================")
	fmt.Printf("任务名称: %s (ID: %s)\n", taskName, taskID)
	fmt.Printf("日志记录: %s\n", taskLog.ID)
	fmt.Printf("执行命令: %s\n", string(taskLog.Command))
	fmt.Printf("最终状态: %s (耗时: %d 毫秒, 退出码: %d)\n", statusText, taskLog.Duration, taskLog.ExitCode)
	if taskLog.StartTime != nil {
		fmt.Printf("开始时间: %s\n", taskLog.StartTime.Time().Format("2006-01-02 15:04:05"))
	}
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Println("[日志输出内容]")

	// 解压
	decompressed, err := utils.DecompressFromBase64(string(taskLog.Output))
	if err != nil {
		fmt.Printf("[无法解压日志输出: %v]\n", err)
	} else {
		// 清理多余回车和终端 ANSI 转义字符
		cleanText := strings.ReplaceAll(decompressed, "\r\n", "\n")
		cleanText = clibase.AnsiRegex.ReplaceAllString(cleanText, "")
		fmt.Println(strings.TrimSpace(cleanText))
	}

	if string(taskLog.Error) != "" {
		fmt.Println("\n[系统捕获异常]")
		fmt.Println(string(taskLog.Error))
	}
	fmt.Println("====================================================================================================")
}

func runHistory(args []string) {
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	limitPtr := fs.Int("limit", 10, "展示的最近历史记录条数")

	fs.Usage = func() {
		clibase.PrintSubCommandUsage("白虎面板任务执行历史查看工具", "baihu task history <任务ID> [参数]", "  baihu task history a1b2c3d4\n  baihu task history a1b2c3d4 -limit 20", fs)
	}

	if err := fs.Parse(args); err != nil {
		return
	}

	parsedArgs := fs.Args()
	if len(parsedArgs) < 1 {
		fmt.Fprintf(os.Stderr, "错误: 缺少目标任务ID。\n")
		fs.Usage()
		return
	}
	taskID := parsedArgs[0]

	clibase.InitContext(false)

	var task models.Task
	database.DB.Where("id = ?", taskID).Limit(1).Find(&task)
	taskName := taskID
	if task.Name != "" {
		taskName = task.Name
	}

	var logs []models.TaskLog
	database.DB.Where("task_id = ?", taskID).Order("created_at DESC").Limit(*limitPtr).Find(&logs)

	fmt.Println("====================================================================================================")
	fmt.Printf("任务流水: %s (ID: %s) 的近期执行记录 (最多展示 %d 条)\n", taskName, taskID, *limitPtr)
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Printf("%-20s | %-8s | %-6s | %-12s | %-20s\n", "日志ID", "状态", "退出码", "耗时", "开始时间")
	fmt.Println("----------------------------------------------------------------------------------------------------")

	if len(logs) == 0 {
		fmt.Println("未查询到任何历史执行记录。")
	} else {
		for _, l := range logs {
			statusText := "运行中"
			switch l.Status {
			case constant.TaskStatusSuccess:
				statusText = "成功"
			case constant.TaskStatusFailed:
				statusText = "失败"
			case constant.TaskStatusTimeout:
				statusText = "超时"
			case constant.TaskStatusCancelled:
				statusText = "已取消"
			}

			startStr := "-"
			if l.StartTime != nil {
				startStr = l.StartTime.Time().Format("2006-01-02 15:04:05")
			}
			durationStr := fmt.Sprintf("%d ms", l.Duration)

			fmt.Printf("%-20s | %-8s | %-6d | %-12s | %-20s\n", l.ID, statusText, l.ExitCode, durationStr, startStr)
		}
	}
	fmt.Println("====================================================================================================")
	fmt.Printf("提示: 结合命令 'baihu task status %s <日志ID>' 查看特定历史日志内容。\n", taskID)
}
