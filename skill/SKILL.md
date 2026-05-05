---
name: bangumi
description: Manage Bangumi (番组计划) anime tracking. Use this skill whenever the user mentions anime tracking, Bangumi, 番组计划, bgm.tv, wants to search for anime, mark episodes as watched, manage their watchlist, view daily broadcast schedules, check their collection, or explore anime characters and voice actors (声优). Also trigger for requests involving anime progress tracking, rating anime, or finding new anime to watch.
---

# Bangumi Skill

通过 `skill/bangumi` (Linux/macOS) 或 `skill/bangumi.exe` (Windows) 命令行工具管理 Bangumi (bgm.tv) 账号

## 前置条件

首次使用前必须配置令牌：

```bash
skill/bangumi auth login --token <access_token>
```

令牌申请: https://next.bgm.tv/demo/access-token

> 如果令牌失效，CLI 会自动清除并提示重新设置。

## 可选配置

```bash
skill/bangumi config proxy http://127.0.0.1:1081   # HTTP 代理
skill/bangumi config timeout 60                     # 请求超时（秒），默认60
skill/bangumi config proxy                          # 查看当前配置
```

## 全局标志

所有命令都支持：
- `--proxy <url>` — HTTP 代理（临时覆盖）
- `--token <token>` — 令牌（临时覆盖，不存储）
- `--output json` — JSON 格式输出（默认人类友好文本）
- `-v` / `--debug` — 日志级别

## 核心命令

### 1. 搜索

```bash
# 搜索动画，可选关键词 + 筛选
skill/bangumi search subjects "AIR" --sort rank --limit 5
skill/bangumi search subjects --filter-type 2 --filter-air-date ">=2026-05-01" --sort rank
skill/bangumi search subjects --filter-rating ">=7" --sort score --limit 10

# 搜索角色和人物
skill/bangumi search characters "神尾观铃"
skill/bangumi search persons "神尾观铃"
```

**筛选参数：**
| 参数 | 说明 |
|------|------|
| `--filter-type 1/2/3/4/6` | 书籍/动画/音乐/游戏/三次元 |
| `--filter-air-date ">=2020-01-01"` | 播出日期 |
| `--filter-rating ">=7"` | 评分筛选 |
| `--filter-rank "<=100"` | 排名筛选 |
| `--sort match/heat/rank/score` | 排序方式 |

### 2. 每日放送（无需令牌）

```bash
skill/bangumi calendar
```

### 3. 条目详情

所有命令默认接受作品名称（自动搜索匹配），用 `--id` 精确指定：

```bash
# 名称搜索（推荐）
skill/bangumi subject get "AIR"
skill/bangumi subject characters "AIR"
skill/bangumi subject persons "AIR"
skill/bangumi subject relations "AIR"

# ID 精确查找
skill/bangumi subject get --id 12
```

### 4. 追番管理 ⭐

**完整追番流程：**

```bash
# 标记"在看"
skill/bangumi collection update "AIR" --type 3

# 标记第1话看过
skill/bangumi collection update-episode "AIR" --ep 1 --type 2

# 看完评分
skill/bangumi collection update "AIR" --type 2 --rate 9

# 加标签
skill/bangumi collection update "AIR" --tags "经典,催泪"
```

**收藏类型：** `1=想看 2=看过 3=在看 4=搁置 5=抛弃`

**更多命令：**
```bash
# 查看某人收藏
skill/bangumi collection list <用户名> --subject-type 2

# 查看观看进度
skill/bangumi collection episodes "AIR"

# 查看某条收藏详情
skill/bangumi collection get <用户名> <条目ID>
```

### 5. 章节信息

```bash
# 列出所有章节
skill/bangumi episode list "AIR"
skill/bangumi episode list --id 12 --type 0     # 仅本篇

# 章节详情
skill/bangumi episode get <章节ID>
```

### 6. 角色和人物

```bash
# 角色
skill/bangumi character get "神尾观铃"
skill/bangumi character subjects "神尾观铃"      # 角色出演作品
skill/bangumi character persons "神尾观铃"       # 角色声优
skill/bangumi character collect "神尾观铃"       # 收藏角色
skill/bangumi character uncollect "神尾观铃"     # 取消收藏

# 人物（声优/导演）
skill/bangumi person get "神尾观铃"
skill/bangumi person subjects "神尾观铃"         # 参与作品
skill/bangumi person characters "神尾观铃"       # 配音角色
```

### 7. 用户信息

```bash
skill/bangumi user me                              # 当前用户
skill/bangumi user get <用户名>                    # 指定用户
```

### 8. 目录管理（收藏夹）

```bash
skill/bangumi index create "我的推荐" --desc "个人推荐列表"
skill/bangumi index add-subject <目录ID> <条目ID> --comment "神作"
skill/bangumi index subjects <目录ID>
skill/bangumi index collect <目录ID>
```

## 输出格式

默认输出为人类友好文本（AI 可直接阅读）。加 `--output json` 获取结构化 JSON。

## 注意事项

- 所有通过名称搜索的命令（如 `"AIR"`）会在 Bangumi 中自动搜索匹配，取第一个结果。如需精确指定，用 `--id`。
- 章节标记 `update-episode` 支持 `--ep N` 直接指定第 N 话，CLI 会自动查找对应的章节 ID。
- `calendar` 命令不需要令牌，其他命令需要令牌。
- 评论等长文本可用 `--comment-file <路径>` 从文件读取。
