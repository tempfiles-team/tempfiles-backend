package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"mime/multipart"
	"net/http"
	"sort"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

type MultiFileAndHeader struct {
	File   multipart.File
	Header *multipart.FileHeader
}

// An improved version of net/http.Request.FormFile()
func FormFiles(r *http.Request, key string) ([]MultiFileAndHeader, error) {

	if r.MultipartForm == multipartByReader {
		return nil, errors.New("http: multipart handled by MultipartReader")
	}
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {

		var files []MultiFileAndHeader
		if fhs := r.MultipartForm.File[key]; len(fhs) > 0 {

			for _, fh := range fhs {
				f, err := fh.Open()
				if err != nil {
					return nil, err
				}
				files = append(files, MultiFileAndHeader{
					File:   f,
					Header: fh,
				})
			}
		}

		if len(files) > 0 {
			return files, nil
		}
	}
	return nil, http.ErrMissingFile
}

func GenIdFormMulitpart(MFAH []MultiFileAndHeader) (string, error) {

	sort.Slice(MFAH, func(i, j int) bool {
		return MFAH[i].Header.Filename < MFAH[j].Header.Filename
	})

	var hashes [][]byte

	for _, file := range MFAH {

		defer file.File.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file.File)
		fileBytes := buf.Bytes()

		nameHash := sha1.Sum([]byte(file.Header.Filename))
		hashes = append(hashes, nameHash[:])
		fileHash := sha1.Sum(fileBytes)
		hashes = append(hashes, fileHash[:])
	}

	combinedHash := sha1.Sum(bytes.Join(hashes, nil))

	return base64.RawURLEncoding.EncodeToString(combinedHash[:]), nil
}
