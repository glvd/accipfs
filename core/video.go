package core

// Video ...
type Video struct {
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
