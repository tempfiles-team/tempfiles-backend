package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minpeter/tempfiles-backend/database"
)

func CreateJWTToken(fileRow database.FileRow) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 10).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["file id"] = fileRow.Id
	claims["exp"] = exp
	claims["encrypto"] = fileRow.Encrypto
	claims["isAdmin"] = false
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func IsEncrypted(fileName string) bool {
	fileRow := new(database.FileRow)
	has, err := database.Engine.Where("file_name = ?", fileName).Desc("id").Get(fileRow)
	if err != nil {
		return false
	}
	if !has {
		return false
	}
	return fileRow.Encrypto
}
