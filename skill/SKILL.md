---
name: bangumi
description: Manage Bangumi (番组计划) anime tracking. Use this skill whenever the user mentions anime tracking, Bangumi, 番组计划, bgm.tv, wants to search for anime, mark episodes as watched, manage their watchlist, view daily broadcast schedules, check their collection, or explore anime characters and voice actors (声优). Also trigger for requests involving anime progress tracking, rating anime, or finding new anime to watch.
---

# Bangumi Skill

通过 `skill/bangumi` (Linux/macOS) 或 `skill/bangumi.exe` (Windows) 命令行工具管理 Bangumi (bgm.tv) 账号

## ⚡ 核心原则

**1. 直接传名字，不要先搜索再标记！**

```bash
# ✅ 正确 — 一条命令搞定
skill/bangumi collection update "史上最强的大魔王转生为村民A" --type 3
skill/bangumi collection update-episode "史上最强的大魔王转生为村民A" --ep 1 --type 2

# ❌ 错误 — 不要这样！
skill/bangumi search subjects "xxx"    # ← 不需要！
skill/bangumi collection update 12345 --type 3  # ← 多此一举！
```

**2. "标记在看第N话" = 两条命令：**

```bash
skill/bangumi collection update "作品名" --type 3           # 标记在看
skill/bangumi collection update-episode "作品名" --ep N --type 2  # 标记第N话看过
```

**3. 名称找不到时先搜索确认：**

```bash
skill/bangumi search subjects "关键词" --limit 3   # 确认准确名称
```

---

## 前置条件

```bash
skill/bangumi auth login --token <access_token>
```

令牌申请: https://next.bgm.tv/demo/access-token

> 令牌失效时 CLI 自动清除并提示重新设置。

## 可选配置

```bash
skill/bangumi config proxy http://127.0.0.1:1081
skill/bangumi config timeout 60
```

## 全局标志

- `--proxy <url>` — HTTP 代理
- `--token <token>` — 令牌（覆盖 token.json）
- `--output json` — JSON 输出（默认 AI 友好文本）
- `-v` / `--debug` — 日志

---

## 命令参考

### 搜索

```bash
skill/bangumi search subjects "AIR" --sort rank --limit 5
skill/bangumi search subjects --filter-type 2 --filter-air-date ">=2026-05-01" --sort rank
skill/bangumi search subjects --filter-rating ">=7" --sort score --limit 10
skill/bangumi search characters "神尾观铃"
skill/bangumi search persons "神尾观铃"
```

| 参数 | 说明 |
|------|------|
| `--filter-type 1/2/3/4/6` | 书籍/动画/音乐/游戏/三次元 |
| `--filter-air-date ">=2026-05-01"` | 播出日期 |
| `--filter-rating ">=7"` | 评分筛选 |
| `--filter-rank "<=100"` | 排名筛选 |
| `--sort match/heat/rank/score` | 排序 |

### 每日放送（无需令牌）

```bash
skill/bangumi calendar
```

### 条目详情

输入作品名自动搜索，用 `--id` 精确指定：

```bash
skill/bangumi subject get "AIR"
skill/bangumi subject characters "AIR"    # 角色列表
skill/bangumi subject persons "AIR"       # 制作人员
skill/bangumi subject relations "AIR"     # 关联作品
skill/bangumi subject get --id 12         # ID精确查找
```

### 追番管理 ⭐

```bash
skill/bangumi collection update "AIR" --type 3              # 在看
skill/bangumi collection update-episode "AIR" --ep 1 --type 2  # 第1话看过
skill/bangumi collection update "AIR" --type 2 --rate 9     # 看过+评分
skill/bangumi collection update "AIR" --tags "经典,催泪"     # 加标签
```

收藏类型: `1=想看 2=看过 3=在看 4=搁置 5=抛弃`

**查看收藏（不传参数默认自己）：**

```bash
skill/bangumi collection list                           # 自己的收藏
skill/bangumi collection list sai --subject-type 2 --collection-type 3  # sai在看的动画
skill/bangumi collection episodes "AIR"                 # 章节观看进度
skill/bangumi collection get sai 12                     # sai对条目12的收藏详情
```

**查看收藏的角色/人物：**

```bash
skill/bangumi collection characters              # 自己收藏的角色
skill/bangumi collection characters sai          # sai收藏的角色
skill/bangumi collection persons                 # 自己收藏的人物（声优/导演）
skill/bangumi collection persons sai             # sai收藏的人物
```

> ⚠️ collection 子命令的 `<用户名>` 参数接受的是**用户名**（如 `sai`），不是用户ID数字。

### 章节

```bash
skill/bangumi episode list "AIR"                # 全部章节
skill/bangumi episode list --id 12 --type 0     # 仅本篇
skill/bangumi episode get <章节ID>
```

### 角色

```bash
skill/bangumi character get "神尾观铃"          # 详情 + 出演作品 + 声优（一次返回聚合数据）
skill/bangumi character subjects "神尾观铃"     # 仅出演作品
skill/bangumi character persons "神尾观铃"      # 仅声优
skill/bangumi character collect "神尾观铃"      # 收藏（需令牌）
skill/bangumi character uncollect "神尾观铃"    # 取消收藏
```

> **注意：** `character get` 会自动附带该角色的出演作品列表和声优列表，无需再单独调用 subjects/persons。

### 人物（声优/导演）

```bash
skill/bangumi person get "神尾观铃"             # 详情 + 参与作品 + 配音角色（一次返回聚合数据）
skill/bangumi person subjects "神尾观铃"        # 仅参与作品
skill/bangumi person characters "神尾观铃"      # 仅配音角色
skill/bangumi person collect "神尾观铃"         # 收藏
```

> **注意：** `person get` 会自动附带该人物的参与作品列表和配音角色列表，无需再单独调用 subjects/characters。

### 用户

```bash
skill/bangumi user me                           # 当前用户
skill/bangumi user get sai                      # 指定用户
```

### 目录

```bash
skill/bangumi index create "推荐" --desc "列表"
skill/bangumi index add-subject <目录ID> <条目ID> --comment "备注"
skill/bangumi index subjects <目录ID>
skill/bangumi index collect <目录ID>
```

---

## 注意事项

- 所有名称搜索取第一个匹配结果，用 `--id` 精确指定
- `update-episode --ep N` 自动查找章节 ID，无需手动查
- `calendar` 无需令牌，其他命令需令牌
- 长文本用 `--comment-file <路径>` 从文件读取
