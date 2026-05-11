# 语言依赖

白虎面板深度集成了 **Mise** 运行时管理器，这使得它具备多版本语言环境的高灵活性和隔离性。

## 脚本运行环境

白虎面板原生支持以下脚本的定时执行：
- **Python3**, **Node.js**, **Bash** (标准版镜像内置环境)
- 通过 **Mise** 扩展：支持几乎所有主流编程语言的动态安装与切换。


> [!TIP]
> **Minimal 镜像注意**：如果您使用的是 `minimal` 标签的镜像，系统初始不包含 Python 和 Node.js。您需要进入「编程语言」页面手动点击安装您所需的运行时。

## 依赖管理支持

系统内置了高度集成的跨语言依赖管理器，支持自动化安装和管理以下语言的依赖项，并确保在容器内全局可用：

| 语言 | 包管理器 | 功能说明 |
| :--- | :--- | :--- |
| **Python** | pip | 自动使用内置虚拟环境，支持清华源 |
| **Node.js** | npm | 全局安装模式，自动配置 npmmirror 镜像 |
| **Go** | go install | 通过 `go install` 安装二进制工具 |
| **Rust** | cargo | 通过 `cargo install` 安装 Rust 依赖 |
| **Ruby** | gem | 支持 `gem install` 本地安装 |
| **Bun** | bun | 支持 `bun add -g` 全局模式 |
| **PHP** | composer | 支持 `composer global require` |
| **Deno** | deno | 支持 `deno install -g` |
| **.NET** | dotnet | 支持 `dotnet tool install -g` |
| **Elixir/Erlang** | mix | 支持 `mix archive.install` |
| **Lua** | luarocks | 通过 `luarocks` 管理 Lua 包 |
| **Nim** | nimble | 支持 `nimble install` |
| **Dart/Flutter** | pub | 支持 `pub global activate` |
| **Perl** | cpanm | 简单的 `cpanm` 安装支持 |
| **Crystal** | shards | `shards` 项目级别或工具安装 |

## 使用方法

### 1. 安装环境
进入「编程语言」页面，使用 `mise` 一键安装所需的语言及版本。

### 2. 依赖管理
在已安装列表点击「依赖管理」，输入名称（可选版本）即可自动在对应环境内完成安装。

### 3. 多版本切换
对于复杂的项目，您可以通过面板配置不同的任务版本镜像，系统基于 `mise exec` 实现了完善的环境隔离，不同版本的依赖包互不冲突。

## 常用工具安装

如果您需要在面板环境中使用 Ansible 或其他通过 pipx 管理的工具，可以使用以下命令进行快速安装：

### 安装 Ansible

白虎面板推荐通过 `mise` 结合 `pipx` 安装 Ansible，以保持环境隔离且全局可用：

```bash
# 首先安装 pipx
mise use -g pipx@latest

# 使用 pipx 安装 ansible
mise use -g ansible@latest
```

安装完成后，您可以在「脚本管理」或「定时任务」中直接调用 `ansible` 或 `ansible-playbook` 命令。

## PHP 环境特别说明

在 Docker 容器或精简版 Linux 环境中，通过 `mise` 安装某些 PHP 版本时，如果官方没有提供预编译的二进制包，系统会尝试从源码编译。这需要额外的系统级构建依赖（如 `autoconf`, `bison`, `pkg-config` 等）。

如果您在安装 PHP 时遇到 `autoconf not found` 或 `buildconf failed` 等错误，请参考以下方案：

1. **运行预设脚本**: 我们在项目根目录提供了示例脚本 [install_php_env_deps.sh](file:///example/deps/install_php_env_deps.sh)，它会自动安装所需的系统库。
2. **配置为开机任务**: 由于 Docker 容器层通常是只读的或在重启后重置（非挂载目录），建议在面板中创建一个 **触发类型为 `baihu_startup`** 的定时任务，命令填写 `bash example/deps/install_php_env_deps.sh`。这样可以确保每次容器启动时都会自动准备好 PHP 编译环境。

---

## 隔离机制说明

-   白虎面板通过动态注入 `PATH` 环境和 `mise shims` 将语言环境暴露给系统。
-   每个任务在执行前都会根据任务配置自动加载对应的运行时环境变量。
-   **运行时激活**：自动将 `MISE_DATA_DIR` 等环境变量指向宿主机的持久化挂载目录，确保护持久化可用。
