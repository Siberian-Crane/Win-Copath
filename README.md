# Win-Copath

Windows 右键菜单路径复制工具。在文件和文件夹的右键菜单中添加 "Win Copath" 选项，支持四种路径格式复制。

## 功能

右键点击文件或文件夹时，显示 "Win Copath" 子菜单，包含四个选项：

| 选项             | 格式      | 示例                            |
|----------------|---------|-------------------------------|
| `/ (正斜杠)`      | 正斜杠     | `C:/Users/test/file.txt`      |
| `\ (反斜杠)`      | 反斜杠     | `C:\Users\test\file.txt`      |
| `"/" (引号正斜杠)`  | 引号包裹正斜杠 | `"C:/Users/test/file.txt"`    |
| `"\\" (引号反斜杠)` | 引号包裹反斜杠 | `"C:\\Users\\test\\file.txt"` |

## 安装

### WinGet (推荐)

```powershell
winget install Siberian-Crane.WinCopath
```

### 手动安装

1. 从 [Releases](https://github.com/Siberian-Crane/Win-Copath/releases) 下载 `wcpath_windows_amd64.exe`
2. 将 `wcpath.exe` 放入一个已添加到 PATH 的目录（如 `C:\Windows\System32`）

## 使用

安装后，**默认不会修改右键菜单**。需要手动启用：

```powershell
# 以管理员身份运行 PowerShell
wcpath on     # 添加 Win Copath 到右键菜单
wcpath off    # 从右键菜单移除
wcpath help   # 查看帮助
```

启用后，右键点击任意文件或文件夹即可看到 "Win Copath" 子菜单。

## 要求

- Windows 10 及以上
- 管理员权限（用于修改注册表）

## 开发

```bash
# 编译
go build -o wcpath.exe .

# 测试复制命令
.\wcpath.exe copy fwd "C:\test\file.txt"
```

## License

MIT
