package models

import "time"

type CommonTracking struct {
	Id            int64     `json:"-"`
	ItemId        string    `json:"id"`
	UploadDate    time.Time `json:"uploadDate"`
	IsEncrypted   bool      `json:"isEncrypted"`
	Password      string    `json:"-"`
	DownloadCount int64     `json:"downloadCount"`
	DownloadLimit int64     `json:"downloadLimit"`
	ExpireTime    time.Time `json:"expireTime"`
	IsDeleted     bool      `json:"-"`
}
