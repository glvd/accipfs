package core

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/basis/hash"
	"strings"
)

// DataHashEncoder ...
type DataHashEncoder interface {
	Hash() string
	Verify(s string) bool
}

// DataJSONer ...
type DataJSONer interface {
	JSON() string
}

// DataRooter ...
type DataRooter interface {
	Root() string
}

// Marshaler ...
type Marshaler interface {
	Marshal() ([]byte, error)
}

// Unmarshaler ...
type Unmarshaler interface {
	Unmarshal([]byte) error
}

// JSONer ...
type JSONer interface {
	Marshaler
	Unmarshaler
}

// Serializable ...
type Serializable interface {
	JSONer
	DataRooter
	DataHashEncoder
	DataJSONer
}

// MediaInfo ...
type MediaInfo struct {
	No           string   `json:"no"`            //编号
	Intro        string   `json:"intro"`         //简介
	Alias        []string `json:"alias"`         //别名，片名
	Role         []string `json:"role"`          //主演
	Director     string   `json:"director"`      //导演
	Systematics  string   `json:"systematics"`   //分级
	Season       string   `json:"season"`        //季
	TotalEpisode string   `json:"total_episode"` //总集数
	Episode      string   `json:"episode"`       //集数
	Producer     string   `json:"producer"`      //生产商
	Publisher    string   `json:"publisher"`     //发行商
	Type         string   `json:"type"`          //类型：film，FanDrama
	Format       string   `json:"format"`        //输出格式：3D，2D,VR(VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽)
	Language     string   `json:"language"`      //语言
	Caption      string   `json:"caption"`       //字幕
	Group        string   `json:"group"`         //分组
	Index        string   `json:"index"`         //索引
	Date         string   `json:"date"`          //发行日期
	Sharpness    string   `json:"sharpness"`     //清晰度
	Series       string   `json:"series"`        //系列
	Tags         []string `json:"tags"`          //标签
	Length       string   `json:"length"`        //时长
	Sample       []string `json:"sample"`        //样板图
	Uncensored   bool     `json:"uncensored"`    //有码,无码
}

// Security ...
type Security struct {
	Key string `xorm:"key" json:"key"`
}

// DataInfoVersion1 ...
var DataInfoVersion1, _ = ParseVersion("v0.0.1")

// Info ...
type Info struct {
	ThumbHash  string `xorm:"thumb_hash" json:"thumb_hash"`   //缩略图
	ThumbURI   string `xorm:"thumb_uri" json:"thumb_uri"`     //缩略图
	PosterHash string `xorm:"poster_hash" json:"poster_hash"` //海报地址
	PosterURI  string `xorm:"poster_uri" json:"poster_uri"`   //缩略图
}

// DataInfoV1 ...
type DataInfoV1 struct {
	RootHash   string    `xorm:"root_hash" json:"root_hash"`     //源信息
	MediaInfo  MediaInfo `xorm:"media_info" json:"media_info"`   //媒体信息
	MediaURI   string    `xorm:"media_uri" json:"media_uri"`     //入口地址
	MediaHash  string    `xorm:"media_hash" json:"media_hash"`   //入口HASH
	MediaIndex string    `xorm:"media_index" json:"media_index"` //入口名称
	Info       Info      `xorm:"info" json:"info"`               //补充信息
	InfoURI    string    `xorm:"info_uri" json:"info_uri"`       //补充信息地址
	Security   Security  `xorm:"security" json:"security"`       //安全验证
	LastUpdate int64     `xorm:"last_update" json:"last_update"` //最后更新时间
	Version    Version   `xorm:"version" json:"version"`         //版本
}

// JSON ...
func (v *DataInfoV1) JSON() string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// Marshal ...
func (v *DataInfoV1) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (v *DataInfoV1) Unmarshal(b []byte) error {
	return json.Unmarshal(b, v)
}

// Root ...
func (v *DataInfoV1) Root() string {
	return v.RootHash
}

// Hash ...
func (v *DataInfoV1) Hash() string {
	sum, err := hash.Sum(v)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", sum)
}

// Verify ...
func (v *DataInfoV1) Verify(hash string) bool {
	return strings.Compare(v.Hash(), hash) == 0
}

// VerifyVersion ...
func (v *DataInfoV1) VerifyVersion() bool {
	return v.Version.Compare(DataInfoVersion1) == 0
}
