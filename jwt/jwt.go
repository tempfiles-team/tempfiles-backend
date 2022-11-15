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
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET") + FileTracking.FileId))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func IsEncrypted(id string) bool {
	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)
	if err != nil {
		return false
	}
	if !has {
		return false
	}

	return !FileTracking.IsEncrypted
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
