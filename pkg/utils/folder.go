package utils

import (
	"os"
)

func CheckTmpFolder() error {
	//tmp folder check
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err := os.Mkdir("tmp", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckFileFolder(fileId string) error {
	//tmp 하위 폴더 존재확인

	if _, err := os.Stat("tmp/" + fileId); os.IsNotExist(err) {
		err := os.MkdirAll("tmp/"+fileId, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
