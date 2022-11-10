package jwt

import (
	"fmt"
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

func IsEncrypted(id, fileName string) bool {
	FileTracking := database.FileTracking{
		FileName: fileName,
		FileId:   id,
	}

	// var user = User{ID: 27}
	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return false
	}
	if !has {
		return false
	}

	return !FileTracking.IsEncrypted
}

func IsMatched() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if !token.Claims.(jwt.MapClaims)["isAdmin"].(bool) {
			return nil, fmt.Errorf("not admin")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}
}
