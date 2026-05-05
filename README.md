# Bangumi Skill

[Bangumi (bgm.tv)](https://bgm.tv) 的 AI Agent Skill — 让 Agent 直接操作番剧追番、条目搜索、角色/声优查询。

## 🤖 AI 安装说明（给 Agent 看的）

**一句话：** 下载 Release zip → 解压到项目 `skill/` 目录 → 配置令牌 → Agent 自动启用。

**目录结构：**
```
你的项目/
└── skill/
    ├── SKILL.md          ← Agent 读到这个自动加载技能
    └── bangumi(.exe)     ← SKILL.md 里命令调用的工具
```

**完整命令：**

```bash
# 1. 从 Release 下载对应平台 zip 并解压到 skill/
# https://github.com/kasuganosora/bangumi.skill/releases

# 2. 配置令牌（只需一次）
./skill/bangumi auth login --token "你的令牌"
# 令牌申请: https://next.bgm.tv/demo/access-token

# 3. 完成 — Agent 会自动根据用户意图调用 skill/bangumi
```

> **原理：** Agent 扫描工作目录下的 `skill/SKILL.md`，读取技能定义后，根据用户意图自动选择命令执行。用户无需手动教 Agent 怎么用。

---

## 快速开始

```bash
# 查看今日放送
skill/bangumi calendar

# 搜索动画
skill/bangumi search subjects "AIR" --sort rank

# 查看帮助
skill/bangumi --help
```

## 功能概览

| 功能 | 命令 | 说明 |
|------|------|------|
| 每日放送 | `calendar` | 本周每天播出的动画 |
| 搜索 | `search subjects/characters/persons` | 搜索条目、角色、人物 |
| 条目详情 | `subject get/characters/persons/relations` | 获取作品信息、角色、制作人员 |
| 章节 | `episode list/get` | 查看剧集列表和详情 |
| 追番管理 | `collection update/update-episode/episodes` | 标记在看/看过、评分、标记单集 |
| 角色 | `character get/subjects/persons/collect` | 角色详情、作品、声优 |
| 人物 | `person get/subjects/characters/collect` | 声优/导演详情 |
| 目录 | `index create/edit/add-subject/...` | 管理自定义条目目录 |
| 认证 | `auth login/status/logout` | 令牌管理 |
| 配置 | `config proxy/timeout` | 代理、超时设置 |

## AI 友好设计

所有命令默认通过名称搜索，无需手动查找 ID：

```bash
# AI 直接说作品名即可
skill/bangumi subject get "AIR"
skill/bangumi character get "神尾观铃"

# 追番只需两行
skill/bangumi collection update "AIR" --type 3          # 标记在看
skill/bangumi collection update-episode "AIR" --ep 1 --type 2  # 第1话看过
```

输出默认为人类可读文本，可用 `--output json` 切换为 JSON。

## 开发

```bash
cd cli
go build -o bin/bangumi .
go test ./...
```

Git pre-commit hook 自动执行 `gofmt` + `golangci-lint`。

## CI / Release

Push → lint + test。打 `v*` tag → 构建 6 平台 Release（Linux/Windows/macOS × x86_64/ARM64）。

## 项目结构

```
bangumi.skill/
├── skill/
│   ├── SKILL.md       # Agent Skill 定义
│   └── bangumi(.exe)  # 二进制（构建产物）
├── cli/               # Go 源码
│   ├── api/           # Bangumi API 客户端（50+ 端点）
│   ├── cmd/           # Cobra CLI 命令
│   ├── internal/      # 配置/输出模块
│   └── log/           # slog 日志封装
├── scripts/           # pre-commit hook 等
└── .github/workflows/ # CI/CD
```
