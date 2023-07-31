package models

type TextTracking struct {
	CommonTracking `gorm:"embedded"`
	TextId         string `json:"id"`
	TextData       string `json:"data"`
	TextCount      int64  `json:"numberOfText"`
}
