package api

import (
	"encoding/json"
	"fmt"
)

// ---------------------------------------------------------------------------
// 枚举类型
// ---------------------------------------------------------------------------

// SubjectType 条目类型
type SubjectType int

const (
	SubjectBook  SubjectType = 1
	SubjectAnime SubjectType = 2
	SubjectMusic SubjectType = 3
	SubjectGame  SubjectType = 4
	SubjectReal  SubjectType = 6
)

// EpisodeType 章节类型
type EpisodeType int

const (
	EpisodeNormal  EpisodeType = 0 // 本篇
	EpisodeSpecial EpisodeType = 1 // 特别篇
	EpisodeOP      EpisodeType = 2 // OP
	EpisodeED      EpisodeType = 3 // ED
	EpisodePreview EpisodeType = 4 // 预告/宣传/广告
	EpisodeMAD     EpisodeType = 5 // MAD
	EpisodeOther   EpisodeType = 6 // 其他
)

// EpisodeStatus 放送状态
type EpisodeStatus string

const (
	EpisodeAir   EpisodeStatus = "Air"
	EpisodeToday EpisodeStatus = "Today"
	EpisodeNA    EpisodeStatus = "NA"
)

// UserGroup 用户组
type UserGroup int

const (
	UserGroupAdmin         UserGroup = 1  // 管理员
	UserGroupBangumiAdmin  UserGroup = 2  // Bangumi 管理猿
	UserGroupSkylightAdmin UserGroup = 3  // 天窗管理猿
	UserGroupMuted         UserGroup = 4  // 禁言用户
	UserGroupBanned        UserGroup = 5  // 禁止访问用户
	UserGroupCharAdmin     UserGroup = 8  // 人物管理猿
	UserGroupWikiAdmin     UserGroup = 9  // 维基条目管理猿
	UserGroupUser          UserGroup = 10 // 用户
	UserGroupWikiUser      UserGroup = 11 // 维基人
)

// ---------------------------------------------------------------------------
// 通用结构
// ---------------------------------------------------------------------------

// Images 图片尺寸集合
type Images struct {
	Large  string `json:"large"`
	Common string `json:"common"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Grid   string `json:"grid"`
}

// Rating 评分信息
type Rating struct {
	Total int            `json:"total"`
	Count map[string]int `json:"count"`
	Score float64        `json:"score"`
}

// Collection 收藏统计
type Collection struct {
	Wish    int `json:"wish"`
	Collect int `json:"collect"`
	Doing   int `json:"doing"`
	OnHold  int `json:"on_hold"`
	Dropped int `json:"dropped"`
}

// ---------------------------------------------------------------------------
// 条目 (Subject)
// ---------------------------------------------------------------------------

// SubjectSmall 条目基本信息
type SubjectSmall struct {
	ID         int         `json:"id"`
	URL        string      `json:"url"`
	Type       SubjectType `json:"type"`
	Name       string      `json:"name"`
	NameCN     string      `json:"name_cn"`
	Summary    string      `json:"summary"`
	AirDate    string      `json:"air_date"`
	AirWeekday int         `json:"air_weekday"`
	Images     Images      `json:"images"`
	Eps        int         `json:"eps"`
	EpsCount   int         `json:"eps_count"`
	Rating     Rating      `json:"rating"`
	Rank       int         `json:"rank"`
	Collection Collection  `json:"collection"`
}

// SubjectMedium 条目 + 角色 + 制作人员
type SubjectMedium struct {
	SubjectSmall
	Characters []CharacterRole `json:"crt"`
	Staff      []PersonRole    `json:"staff"`
}

// SubjectLarge 条目 + 章节 + 讨论 + 日志
type SubjectLarge struct {
	SubjectMedium
	Episodes []Episode `json:"eps"`
	Topics   []Topic   `json:"topic"`
	Blogs    []Blog    `json:"blog"`
}

// ---------------------------------------------------------------------------
// 章节 (Episode)
// ---------------------------------------------------------------------------

// Episode 章节信息
type Episode struct {
	ID       int           `json:"id"`
	URL      string        `json:"url"`
	Type     EpisodeType   `json:"type"`
	Sort     int           `json:"sort"`
	Name     string        `json:"name"`
	NameCN   string        `json:"name_cn"`
	Duration string        `json:"duration"`
	AirDate  string        `json:"airdate"`
	Comment  int           `json:"comment"`
	Desc     string        `json:"desc"`
	Status   EpisodeStatus `json:"status"`
}

// ---------------------------------------------------------------------------
// 讨论 & 日志
// ---------------------------------------------------------------------------

// Topic 讨论版
type Topic struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	MainID    int    `json:"main_id"`
	Timestamp int64  `json:"timestamp"`
	LastPost  int64  `json:"lastpost"`
	Replies   int    `json:"replies"`
	User      User   `json:"user"`
}

// Blog 日志
type Blog struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Image     string `json:"image"`
	Replies   int    `json:"replies"`
	Timestamp int64  `json:"timestamp"`
	Dateline  string `json:"dateline"`
	User      User   `json:"user"`
}

// ---------------------------------------------------------------------------
// 用户 (User)
// ---------------------------------------------------------------------------

// User 用户信息
type User struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Avatar    Images    `json:"avatar"`
	Sign      string    `json:"sign"`
	UserGroup UserGroup `json:"usergroup"`
}

// ---------------------------------------------------------------------------
// 人物 (Character / Person)
// ---------------------------------------------------------------------------

// MonoBase 人物基础模型
type MonoBase struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Name   string `json:"name"`
	Images Images `json:"images"`
}

// Mono 人物
type Mono struct {
	MonoBase
	NameCN   string `json:"name_cn"`
	Comment  int    `json:"comment"`
	Collects int    `json:"collects"`
}

// Person 现实人物
type Person struct {
	Mono
	Info MonoInfo `json:"info"`
}

// Character 虚拟角色
type Character struct {
	Mono
	Info   MonoInfo   `json:"info"`
	Actors []MonoBase `json:"actors"`
}

// MonoInfo 人物详细信息
type MonoInfo struct {
	Birth  string      `json:"birth"`
	Height string      `json:"height"`
	Gender string      `json:"gender"`
	Alias  MonoAlias   `json:"alias"`
	Source *MonoSource `json:"source"`
	NameCN string      `json:"name_cn"`
	CV     string      `json:"cv"`
}

// MonoAlias 人物别名
type MonoAlias struct {
	JP     string `json:"jp"`
	Kana   string `json:"kana"`
	Nick   string `json:"nick"`
	Romaji string `json:"romaji"`
	ZH     string `json:"zh"`
}

// MonoSource 引用来源，可以是单个字符串或字符串数组
type MonoSource struct {
	Values []string
}

func (m *MonoSource) UnmarshalJSON(data []byte) error {
	// 尝试解析为字符串
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.Values = []string{s}
		return nil
	}

	// 尝试解析为数组
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		m.Values = arr
		return nil
	}

	return fmt.Errorf("mono_source: expected string or []string, got %s", string(data))
}

func (m MonoSource) MarshalJSON() ([]byte, error) {
	if len(m.Values) == 1 {
		return json.Marshal(m.Values[0])
	}
	return json.Marshal(m.Values)
}

// ---------------------------------------------------------------------------
// 角色/人员关联
// ---------------------------------------------------------------------------

// CharacterRole 条目的角色关联
type CharacterRole struct {
	Character
	RoleName string `json:"role_name"`
}

// PersonRole 条目的制作人员关联
type PersonRole struct {
	Person
	RoleName string   `json:"role_name"`
	Jobs     []string `json:"jobs"`
}
