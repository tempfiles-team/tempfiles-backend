package file

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"mime/multipart"
	"sort"
)

func GenerateFolderId(files []*multipart.FileHeader) (string, error) {

	sort.Slice(files, func(i, j int) bool {
		return files[i].Filename < files[j].Filename
	})

	var hashes [][]byte

	for _, file := range files {
		fileData, err := file.Open()
		if err != nil {
			return "", err
		}
		defer fileData.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(fileData)
		fileBytes := buf.Bytes()

		nameHash := sha1.Sum([]byte(file.Filename))
		hashes = append(hashes, nameHash[:])
		fileHash := sha1.Sum(fileBytes)
		hashes = append(hashes, fileHash[:])
	}

	combinedHash := sha1.Sum(bytes.Join(hashes, nil))

	return base64.RawURLEncoding.EncodeToString(combinedHash[:]), nil
}
