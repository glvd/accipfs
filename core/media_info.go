package core

import "encoding/json"

// MediaInfo ...
type MediaInfo struct {
	No           string   `json:"no"`            //编号
	Intro        string   `json:"intro"`         //简介
	Alias        []string `json:"alias"`         //别名，片名
	ThumbHash    string   `json:"thumb_hash"`    //缩略图
	PosterHash   string   `json:"poster_hash"`   //海报地址
	SourceHash   string   `json:"source_hash"`   //原片地址
	M3U8Hash     string   `json:"m3u8_hash"`     //切片地址
	M3U8Index    string   `json:"m3u8_index"`    //M3U8名
	SecurityKey  string   `json:"security_key"`  //秘钥
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

// MediaInformationVersion1 ...
const MediaInformationVersion1 = iota

// MediaInfoV1 ...
type MediaInfoV1 struct {
	No           string   `xorm:"no" json:"no"`                       //编号
	Intro        string   `xorm:"varchar(2048)" json:"intro"`         //简介
	Alias        []string `xorm:"json" json:"alias"`                  //别名，片名
	ThumbHash    string   `xorm:"thumb_hash" json:"thumb_hash"`       //缩略图
	PosterHash   string   `xorm:"poster_hash" json:"poster_hash"`     //海报地址
	SourceHash   string   `xorm:"source_hash" json:"source_hash"`     //原片地址
	M3U8Hash     string   `xorm:"m3u8_hash" json:"m3u8_hash"`         //切片地址
	Key          string   `xorm:"key"  json:"-"`                      //秘钥
	M3U8         string   `xorm:"m3u8" json:"-"`                      //M3U8名
	Role         []string `xorm:"json" json:"role"`                   //主演
	Director     string   `xorm:"director" json:"director"`           //导演
	Systematics  string   `xorm:"systematics" json:"systematics"`     //分级
	Season       string   `xorm:"season" json:"season"`               //季
	TotalEpisode string   `xorm:"total_episode" json:"total_episode"` //总集数
	Episode      string   `xorm:"episode" json:"episode"`             //集数
	Producer     string   `xorm:"producer" json:"producer"`           //生产商
	Publisher    string   `xorm:"publisher" json:"publisher"`         //发行商
	Type         string   `xorm:"type" json:"type"`                   //类型：film，FanDrama
	Format       string   `xorm:"format" json:"format"`               //输出格式：3D，2D,VR(VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽)
	Language     string   `xorm:"language" json:"language"`           //语言
	Caption      string   `xorm:"caption" json:"caption"`             //字幕
	Group        string   `xorm:"group" json:"-"`                     //分组
	Index        string   `xorm:"index" json:"-"`                     //索引
	Date         string   `xorm:"'date'" json:"date"`                 //发行日期
	Sharpness    string   `xorm:"sharpness" json:"sharpness"`         //清晰度
	Series       string   `xorm:"series" json:"series"`               //系列
	Tags         []string `xorm:"json tags" json:"tags"`              //标签
	Length       string   `xorm:"length" json:"length"`               //时长
	Sample       []string `xorm:"json sample" json:"sample"`          //样板图
	Uncensored   bool     `xorm:"uncensored" json:"uncensored"`       //有码,无码
}

// JSON ...
func (v *MediaInfoV1) JSON() ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
