package file

import (
	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHandler(c *gin.Context) {
	id := c.Param("id")

	pw := c.Query("pw")

	if id == "" || pw == "" {
		c.JSON(400, gin.H{
			"message": "Please provide a file id and password",
			"error":   nil,
			"unlock":  false,
		})
	}

	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "db query error",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	if !has {
		c.JSON(404, gin.H{
			"message": "file not found",
			"error":   nil,
			"unlock":  false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(FileTracking.Password), []byte(pw)); err != nil {
		c.JSON(401, gin.H{
			"message": "password incorrect",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	token, _, err := jwt.CreateJWTToken(FileTracking)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "jwt token creation error",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	c.JSON(200, gin.H{
		"message": "password correct",
		"token":   token,
		"unlock":  true,
	})
}
