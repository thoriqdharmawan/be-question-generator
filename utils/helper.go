package utils

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func SuccessResponseWithMeta(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"meta":   meta,
	})
}

func GenerateMetaData(total int64, limit int, offset int) fiber.Map {
	hasNextPage := (offset + limit) < int(total)
	hasPrevPage := offset > 0

	return fiber.Map{
		"total":       total,
		"limit":       limit,
		"offset":      offset,
		"hasNextPage": hasNextPage,
		"hasPrevPage": hasPrevPage,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
