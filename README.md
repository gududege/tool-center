# CLI Tool Center

管理并运行 CLI 工具的桌面应用，基于 Wails v2 (Go + React) 构建。

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.7+-3178C6?logo=typescript)
![Wails](https://img.shields.io/badge/Wails-2.10+-DF0000?logo=wails)

---

## 功能

- **插件化架构**：在 `plugins/` 目录下放置插件即可自动加载，无需重新编译主程序。
- **JSON Schema 驱动表单**：通过 `schema.json` + `uischema.json` 自动生成参数配置界面。
- **CLI 参数映射**：支持 argument、option、switch、option-array 等多种参数类型。
- **任务管理**：支持多任务并发执行，任务与 Tab 解耦。
- **历史参数**：自动保存最近 5 次表单参数，方便复用。
- **命令预览**：实时查看由表单数据生成的命令行。
- **输出面板**：统一查看 stdout / stderr / system 日志。

---

## 架构

```
frontend/         ← React 19 + TypeScript + Vite + TailwindCSS + shadcn/ui
  src/features/
    workspace/      ← DynamicForm（RJSF + Shadow DOM）
    sidebar/        ← 插件导航树
    taskPanel/      ← 任务列表
    outputPanel/    ← 日志输出
internal/         ← Go 后端
  adapter/wails/    ← Wails API 层
  application/      ← 应用服务层
  domain/           ← 领域层（Plugin, Task, Output...）
  infrastructure/   ← 基础设施（文件系统、进程管理）
plugins/          ← 插件目录
main.go           ← Wails 入口
```

---

## 开发

### 前置依赖

- [Go 1.25+](https://go.dev/dl/)
- [Node.js 22+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### 安装依赖

```bash
cd frontend && npm install
```

### 开发模式

```bash
wails dev
```

### 构建

```bash
# 当前平台
wails build

# 指定平台
wails build -platform windows/amd64
wails build -platform linux/amd64
wails build -platform darwin/amd64
wails build -platform darwin/arm64
```

### 测试

```bash
# Go 后端
go test ./...

# 前端
cd frontend
npm run test
npm run build
```

---

## 插件开发

每个插件是一个独立目录，包含 `plugin.json`、`schema.json`、`uischema.json` 和可执行文件。

```text
plugins/my-tool/
├── plugin.json
├── schema.json
├── uischema.json
└── my-tool.exe
```

参考文档：

- [Plugin Authoring Guide](./docs/usage.md)
- [Plugin Authoring Skill](./docs/skill.md)
- [Plugin Specification](./docs/Plugin%20Specification.md)

示例插件：

- [`plugins/cli-print-demo`](./plugins/cli-print-demo) — 展示所有支持的 RJSF widget。
- [`plugins/tia-export`](./plugins/tia-export) — 工业自动化导出工具示例。
- [`plugins/tcping`](./plugins/tcping) — 网络工具示例。

---

## 发布

### 自动发布

推送符合 `v*` 格式的 tag 时，GitHub Actions 会自动交叉编译并发布到 GitHub Releases：

```bash
git tag v1.0.0
git push origin v1.0.0
```

支持的构建目标：

| 平台 | 架构 | 产物 |
| --- | --- | --- |
| Windows | amd64 | `cli-tool-center-windows-amd64.zip` |
| Linux | amd64 | `cli-tool-center-linux-amd64.tar.gz` |
| macOS | amd64 | `cli-tool-center-darwin-amd64.zip` |
| macOS | arm64 | `cli-tool-center-darwin-arm64.zip` |

### 手动快照构建

```bash
gh workflow run release.yml --field snapshot=true
```

---

## 许可证

MIT
