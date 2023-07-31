package models

type FileTracking struct {
	CommonTracking `gorm:"embedded"`
	FileId         string `json:"id"`
	FileName       string `json:"filename"`
	FileSize       int64  `json:"size"`
}
