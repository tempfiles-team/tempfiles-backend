package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
)

func CreateJWTToken(fileId string) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 10).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["file_id"] = fileId
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET") + fileId))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func IsEncrypted(id string) bool {

	// 오류가 발생한 경우 암호화된 파일로 간주
	FileS := new(queries.FileState)
	has, err := FileS.GetFile(id)
	if err != nil {
		return true
	}

	if has {
		return FileS.Model.IsEncrypted
	}

	TextS := new(queries.TextState)
	has, err = TextS.GetText(id)

	if err != nil {
		return true
	}

	if has {
		return TextS.Model.IsEncrypted
	}

	return true

}
