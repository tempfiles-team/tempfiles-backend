package configs

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	return fiber.Config{
		AppName:   "tempfiles-backend",
		BodyLimit: int(math.Pow(1024, 3)), // 1 == 1byte, = 1GB
	}
}
