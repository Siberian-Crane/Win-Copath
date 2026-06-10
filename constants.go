package main

const (
	appName       = "Win Copath"
	menuKeyName   = "Win Copath"
	clipboardFile = "wcpath-clip.ps1"
)

// 路径格式定义
type pathFormat struct {
	flag  string // 命令行参数标识
	label string // 右键菜单显示文本
}

var pathFormats = []pathFormat{
	{flag: "fwd", label: "/ (正斜杠)"},
	{flag: "bck", label: "\\ (反斜杠)"},
	{flag: "qfwd", label: `"/" (引号正斜杠)`},
	{flag: "qbck", label: `"\\" (引号反斜杠)`},
}
