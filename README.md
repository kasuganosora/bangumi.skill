# Bangumi Skill

AI Agent 友好的 [Bangumi (bgm.tv)](https://bgm.tv) 命令行工具，用于管理番剧追番、条目搜索、角色/声优查询等功能。作为 Ai Agent Skill 使用，也可独立运行。

## 快速开始

### 1. 下载

从 [Releases](https://github.com/kasuganosora/bangumi.skill/releases) 下载对应平台的 zip 包，解压到 `skill/` 目录。

### 2. 设置令牌

```bash
skill/bangumi auth login --token "你的令牌"
```

令牌申请：https://next.bgm.tv/demo/access-token

### 3. 开始使用

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

## 构建

```bash
cd cli
go build -o bin/bangumi .
go test ./...
```

## CI / Release

Push 到 main/master 触发 lint + test。打 `v*` tag 自动构建 6 平台（Linux/Windows/macOS × x86_64/ARM64）并发布 Release。

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
└── .github/workflows/ # CI/CD
```
