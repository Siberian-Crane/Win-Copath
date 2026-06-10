# Win-Copath

[README_EN.md](https://github.com/Siberian-Crane/Win-Copath/blob/main/README_EN.md)

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

### Chocolatey

```powershell
choco install wincopath
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

## 注册表位置

启用右键菜单（`wcpath on`）后，软件会在以下位置创建注册表项：

- 计算机\HKEY_CLASSES_ROOT\Directory\Background\shell\Win Copath
- 计算机\HKEY_CLASSES_ROOT\Directory\ContextMenus\Win-Copath

> 说明：实际还会创建 `HKEY_CLASSES_ROOT\Directory\shell\Win Copath`（右键文件夹）和 `HKEY_CLASSES_ROOT\*\shell\Win Copath`（右键文件）两个父项，三者都通过 `ExtendedSubCommandsKey` 指向同一份 `Directory\ContextMenus\Win-Copath` 子菜单定义，避免重复注册。

如需手动排查或清理，可用 `regedit` 打开上述位置，或直接运行 `wcpath off` 移除。

## 开发协作

本项目由 `mimo-v2.5-pro` 和 `MiniMax-M3` 两款模型共同协作开发。

## 参考

本软件的右键菜单注册表设计参考了 PowerShell 7 的实现方式，使用 `ExtendedSubCommandsKey` 共享子菜单定义，避免重复注册。

## License

MIT
