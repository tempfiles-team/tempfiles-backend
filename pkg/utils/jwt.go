package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
)

func CreateJWTToken(FileTracking models.FileTracking) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 10).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = FileTracking.Id
	claims["file_id"] = FileTracking.FileId
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET") + FileTracking.FileId))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func IsEncrypted(id string) bool {
	FileS := new(queries.FileState)
	has, err := FileS.GetFile(id)
	if err != nil {
		return false
	}

	if !has {
		return false
	}
	return !FileS.Model.IsEncrypted
}

var FileId string

func IsMatched() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, nil
		}
		fileId, ok := claims["file_id"].(string)
		if !ok {
			return nil, nil
		}

		if fileId != FileId {
			return nil, nil
		}

		return []byte(os.Getenv("JWT_SECRET") + fileId), nil
	}
}
