package models

import "time"

type FileTracking struct {
	Id            int64     `json:"-"`
	FileId        string    `json:"fileId"`
	FileName      string    `json:"filename"`
	FileSize      int64     `json:"size"`
	UploadDate    time.Time `json:"uploadDate"`
	IsEncrypted   bool      `json:"isEncrypted"`
	Password      string    `json:"-"`
	DownloadCount int64     `json:"downloadCount"`
	DownloadLimit int64     `json:"downloadLimit"`
	ExpireTime    time.Time `json:"expireTime"`
	IsDeleted     bool      `json:"-"`
}
