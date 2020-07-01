package core

import "encoding/json"

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
	Version    Version   `xorm:"version" json:"version"`         //版本
}

// JSON ...
func (v *DataInfoV1) JSON() ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

// VerifyVersion ...
func (v *DataInfoV1) VerifyVersion() bool {
	return v.Version.Compare(DataInfoVersion1) == 0
}
