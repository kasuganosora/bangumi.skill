package api

import "time"

// ===========================================================================
// 分页
// ===========================================================================

// Paged 分页结果
type Paged[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ===========================================================================
// Subject 条目 (v0)
// ===========================================================================

// Subject v0 条目详情
type Subject struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	NameCN        string            `json:"name_cn"`
	Summary       string            `json:"summary"`
	Type          SubjectType       `json:"type"`
	Date          string            `json:"date"`
	Platform      string            `json:"platform"`
	Images        Images            `json:"images"`
	Infobox       interface{}       `json:"infobox"`
	Volumes       int               `json:"volumes"`
	Eps           int               `json:"eps"`
	TotalEpisodes int               `json:"total_episodes"`
	Rating        Rating            `json:"rating"`
	Collection    SubjectCollection `json:"collection"`
	Tags          []SubjectTag      `json:"tags"`
	NSFW          bool              `json:"nsfw"`
	Locked        bool              `json:"locked"`
	Series        bool              `json:"series"`
	SeriesEntry   int               `json:"series_entry"`
}

// SubjectCollection 条目收藏统计
type SubjectCollection struct {
	Wish    int `json:"wish"`
	Collect int `json:"collect"`
	Doing   int `json:"doing"`
	OnHold  int `json:"on_hold"`
	Dropped int `json:"dropped"`
}

// SubjectTag 条目标签
type SubjectTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// V0SubjectRelation 条目关联
type V0SubjectRelation struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	NameCN   string      `json:"name_cn"`
	Relation string      `json:"relation"`
	Type     SubjectType `json:"type"`
}

// V0RelatedSubject 关联条目
type V0RelatedSubject struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	NameCN string      `json:"name_cn"`
	Type   SubjectType `json:"type"`
	Images Images      `json:"images"`
	Staff  string      `json:"staff"`
}

// SubjectCategory 条目分类
type SubjectCategory int

const (
	SubjectCatAnimeTV          SubjectCategory = 2
	SubjectCatAnimeMovie       SubjectCategory = 3
	SubjectCatAnimeOVA         SubjectCategory = 4
	SubjectCatBookComic        SubjectCategory = 1001
	SubjectCatBookNovel        SubjectCategory = 1002
	SubjectCatBookIllustration SubjectCategory = 1003
	SubjectCatMusic            SubjectCategory = 3001
	SubjectCatGamePC           SubjectCategory = 4001
	SubjectCatGameConsole      SubjectCategory = 4002
	SubjectCatGameMobile       SubjectCategory = 4003
	SubjectCatReal             SubjectCategory = 6001
)

// ===========================================================================
// Search
// ===========================================================================

// SearchSubjectRequest 条目搜索请求
type SearchSubjectRequest struct {
	Keyword string              `json:"keyword"`
	Sort    string              `json:"sort,omitempty"` // match, heat, rank, score
	Filter  SearchSubjectFilter `json:"filter,omitempty"`
}

// SearchSubjectFilter 条目搜索过滤
type SearchSubjectFilter struct {
	Type        []SubjectType `json:"type,omitempty"`
	MetaTags    []string      `json:"meta_tags,omitempty"`
	Tag         []string      `json:"tag,omitempty"`
	AirDate     []string      `json:"air_date,omitempty"`
	Rating      []string      `json:"rating,omitempty"`
	RatingCount []string      `json:"rating_count,omitempty"`
	Rank        []string      `json:"rank,omitempty"`
	NSFW        *bool         `json:"nsfw,omitempty"`
}

// SearchCharacterRequest 角色搜索
type SearchCharacterRequest struct {
	Keyword string                `json:"keyword"`
	Filter  SearchCharacterFilter `json:"filter,omitempty"`
}

// SearchCharacterFilter 角色搜索过滤
type SearchCharacterFilter struct {
	NSFW *bool `json:"nsfw,omitempty"`
}

// SearchPersonRequest 人物搜索
type SearchPersonRequest struct {
	Keyword string             `json:"keyword"`
	Filter  SearchPersonFilter `json:"filter,omitempty"`
}

// SearchPersonFilter 人物搜索过滤
type SearchPersonFilter struct {
	Career []string `json:"career,omitempty"`
}

// ===========================================================================
// Character 角色 (v0)
// ===========================================================================

// CharacterType 角色类型
type CharacterType int

const (
	CharTypeCharacter    CharacterType = 1
	CharTypeMechanic     CharacterType = 2
	CharTypeShip         CharacterType = 3
	CharTypeOrganization CharacterType = 4
)

// BloodType 血型
type BloodType int

const (
	BloodA  BloodType = 1
	BloodB  BloodType = 2
	BloodAB BloodType = 3
	BloodO  BloodType = 4
)

// CharacterDetail v0 角色详情
type CharacterDetail struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Type      CharacterType `json:"type"`
	Images    PersonImages  `json:"images"`
	Summary   string        `json:"summary"`
	Locked    bool          `json:"locked"`
	Infobox   []interface{} `json:"infobox"`
	Gender    string        `json:"gender"`
	BloodType *BloodType    `json:"blood_type"`
	BirthYear *int          `json:"birth_year"`
	BirthMon  *int          `json:"birth_mon"`
	BirthDay  *int          `json:"birth_day"`
	Stat      Stat          `json:"stat"`
}

// CharacterFull 角色完整信息（详情 + 出演条目 + 声优）
type CharacterFull struct {
	Detail   CharacterDetail    `json:"detail"`
	Subjects []V0RelatedSubject `json:"subjects"`
	Persons  []CharacterPerson  `json:"persons"`
}

// PersonFull 人物完整信息（详情 + 参与作品 + 配音角色）
type PersonFull struct {
	Detail     PersonDetail       `json:"detail"`
	Subjects   []V0RelatedSubject `json:"subjects"`
	Characters []CharacterPerson  `json:"characters"`
}

// CharacterPerson 角色关联的人物（声优）
type CharacterPerson struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Type          CharacterType `json:"type"`
	Images        PersonImages  `json:"images"`
	SubjectID     int           `json:"subject_id"`
	SubjectType   SubjectType   `json:"subject_type"`
	SubjectName   string        `json:"subject_name"`
	SubjectNameCN string        `json:"subject_name_cn"`
	Staff         string        `json:"staff"`
}

// RelatedCharacter 与条目关联的角色
type RelatedCharacter struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Type     CharacterType `json:"type"`
	Images   PersonImages  `json:"images"`
	Relation string        `json:"relation"`
	Actors   []PersonBase  `json:"actors"`
}

// PersonImages 人物/角色图片
type PersonImages struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Grid   string `json:"grid"`
}

// ===========================================================================
// Person 人物 (v0)
// ===========================================================================

// PersonDetail v0 人物详情
type PersonDetail struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Type      CharacterType `json:"type"`
	Images    PersonImages  `json:"images"`
	Summary   string        `json:"summary"`
	Locked    bool          `json:"locked"`
	Infobox   []interface{} `json:"infobox"`
	Gender    string        `json:"gender"`
	BloodType *BloodType    `json:"blood_type"`
	BirthYear *int          `json:"birth_year"`
	BirthMon  *int          `json:"birth_mon"`
	BirthDay  *int          `json:"birth_day"`
	Career    []string      `json:"career"`
	Stat      Stat          `json:"stat"`
}

// PersonBase 人物简要信息
type PersonBase struct {
	ID     int           `json:"id"`
	Name   string        `json:"name"`
	Type   CharacterType `json:"type"`
	Images PersonImages  `json:"images"`
}

// RelatedPerson 与条目关联的人物
type RelatedPerson struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Type     CharacterType `json:"type"`
	Images   PersonImages  `json:"images"`
	Relation string        `json:"relation"`
	Career   []string      `json:"career"`
}

// ===========================================================================
// Episode 章节 (v0)
// ===========================================================================

// EpType 章节类型
type EpType int

const (
	EpMainStory EpType = 0
	EpSP        EpType = 1
	EpOP        EpType = 2
	EpED        EpType = 3
	EpPV        EpType = 4
	EpMAD       EpType = 5
	EpOtherType EpType = 6
)

// EpisodeDetail v0 章节详情
type EpisodeDetail struct {
	ID              int     `json:"id"`
	Type            EpType  `json:"type"`
	Name            string  `json:"name"`
	NameCN          string  `json:"name_cn"`
	Sort            float64 `json:"sort"`
	Ep              float64 `json:"ep"`
	Airdate         string  `json:"airdate"`
	Comment         int     `json:"comment"`
	Duration        string  `json:"duration"`
	Desc            string  `json:"desc"`
	Disc            int     `json:"disc"`
	SubjectID       int     `json:"subject_id"`
	DurationSeconds int     `json:"duration_seconds"`
}

// ===========================================================================
// User & Collection (v0)
// ===========================================================================

// UserDetail v0 用户详情
type UserDetail struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	UserGroup UserGroup `json:"user_group"`
	Avatar    Avatar    `json:"avatar"`
	Sign      string    `json:"sign"`
}

// Avatar 头像
type Avatar struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

// SubjectCollectionType 条目收藏类型
type SubjectCollectionType int

const (
	CollectionWish    SubjectCollectionType = 1 // 想看
	CollectionDone    SubjectCollectionType = 2 // 看过
	CollectionDoing   SubjectCollectionType = 3 // 在看
	CollectionOnHold  SubjectCollectionType = 4 // 搁置
	CollectionDropped SubjectCollectionType = 5 // 抛弃
)

// EpisodeCollectionType 章节收藏类型
type EpisodeCollectionType int

const (
	EpCollectionNone    EpisodeCollectionType = 0
	EpCollectionWish    EpisodeCollectionType = 1
	EpCollectionDone    EpisodeCollectionType = 2
	EpCollectionDropped EpisodeCollectionType = 3
)

// UserSubjectCollection 用户条目收藏
type UserSubjectCollection struct {
	SubjectID   int                   `json:"subject_id"`
	SubjectType SubjectType           `json:"subject_type"`
	Rate        int                   `json:"rate"`
	Type        SubjectCollectionType `json:"type"`
	Comment     string                `json:"comment"`
	Tags        []string              `json:"tags"`
	EpStatus    int                   `json:"ep_status"`
	VolStatus   int                   `json:"vol_status"`
	UpdatedAt   string                `json:"updated_at"`
	Private     bool                  `json:"private"`
	Subject     SubjectSmall          `json:"subject"`
}

// UserSubjectCollectionUpdate 更新条目收藏请求
type UserSubjectCollectionUpdate struct {
	Type      *SubjectCollectionType `json:"type,omitempty"`
	Rate      *int                   `json:"rate,omitempty"`
	Comment   *string                `json:"comment,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	EpStatus  *int                   `json:"ep_status,omitempty"`
	VolStatus *int                   `json:"vol_status,omitempty"`
	Private   *bool                  `json:"private,omitempty"`
}

// UserEpisodeCollection 用户章节收藏
type UserEpisodeCollection struct {
	EpisodeID int                   `json:"episode_id"`
	Type      EpisodeCollectionType `json:"type"`
}

// UserCharacterCollection 用户角色收藏
type UserCharacterCollection struct {
	ID        int          `json:"id"`
	Type      int          `json:"type"`
	Name      string       `json:"name"`
	Images    PersonImages `json:"images"`
	CreatedAt string       `json:"created_at"`
}

// UserPersonCollection 用户人物收藏
type UserPersonCollection struct {
	ID        int          `json:"id"`
	Type      int          `json:"type"`
	Name      string       `json:"name"`
	Images    PersonImages `json:"images"`
	CreatedAt string       `json:"created_at"`
}

// ===========================================================================
// Revision 编辑历史 (v0)
// ===========================================================================

// Revision 编辑历史基础
type Revision struct {
	ID        int       `json:"id"`
	Type      int       `json:"type"`
	Creator   Creator   `json:"creator"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

// DetailedRevision 详细编辑历史
type DetailedRevision struct {
	ID        int         `json:"id"`
	Type      int         `json:"type"`
	Creator   Creator     `json:"creator"`
	Summary   string      `json:"summary"`
	CreatedAt time.Time   `json:"created_at"`
	Data      interface{} `json:"data"`
}

// Creator 创建者
type Creator struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// Stat 统计
type Stat struct {
	Comments int `json:"comments"`
	Collects int `json:"collects"`
}

// ===========================================================================
// Index 目录 (v0)
// ===========================================================================

// Index 目录
type Index struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Total     int       `json:"total"`
	Stat      Stat      `json:"stat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Creator   Creator   `json:"creator"`
	NSFW      bool      `json:"nsfw"`
}

// IndexSubject 目录条目
type IndexSubject struct {
	ID      int    `json:"id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	AddedAt string `json:"added_at"`
	Images  Images `json:"images"`
}

// NewIndexRequest 创建目录请求
type NewIndexRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	NSFW  bool   `json:"nsfw"`
}

// UpdateIndexRequest 更新目录请求
type UpdateIndexRequest struct {
	Title *string `json:"title,omitempty"`
	Desc  *string `json:"desc,omitempty"`
	NSFW  *bool   `json:"nsfw,omitempty"`
}

// AddIndexSubjectRequest 添加目录条目请求
type AddIndexSubjectRequest struct {
	SubjectID int    `json:"subject_id"`
	Comment   string `json:"comment,omitempty"`
}

// EditIndexSubjectRequest 编辑目录条目请求
type EditIndexSubjectRequest struct {
	Comment string `json:"comment"`
}

// ===========================================================================
// Error
// ===========================================================================

// ErrorDetail API 错误详情
type ErrorDetail struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Details     interface{} `json:"details,omitempty"`
}
