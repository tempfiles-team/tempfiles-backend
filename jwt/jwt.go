package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minpeter/tempfiles-backend/database"
)

func CreateJWTToken(FileTracking database.FileTracking) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 10).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = FileTracking.Id
	claims["file_id"] = FileTracking.FileId
	claims["file_name"] = FileTracking.FileName
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func IsEncrypted(fileName string) bool {
	FileTracking := new(database.FileTracking)
	has, err := database.Engine.Where("file_name = ?", fileName).Desc("id").Get(FileTracking)
	if err != nil {
		return false
	}
	if !has {
		return false
	}
	return FileTracking.IsEncrypted
}
