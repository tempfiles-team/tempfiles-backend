package middleware

import (
	"os"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

func JWTProtected() func(*fiber.Ctx) error {
	var id string
	config := jwtware.Config{
		TokenLookup:  "query:token",
		ErrorHandler: jwtError,

		Filter: func(c *fiber.Ctx) bool {

			id = strings.Split(c.OriginalURL(), "/")[2]
			if strings.Contains(id, "?") {
				id = strings.Split(id, "?")[0]
			}

			// true인 경우 스킵
			return !utils.IsEncrypted(id)
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, nil
			}
			fileId, ok := claims["file_id"].(string)
			if !ok {
				return nil, nil
			}
			if fileId != id {
				return nil, nil
			}
			return []byte(os.Getenv("JWT_SECRET") + fileId), nil
		},
	}

	return jwtware.New(config)
}
func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Missing or malformed JWT"))
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(utils.NewFailMessageResponse("Invalid or expired JWT"))
}
