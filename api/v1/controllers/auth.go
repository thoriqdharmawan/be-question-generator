package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/entity"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/request"
	"github.com/thoriqdharmawan/be-question-generator/db"
	"github.com/thoriqdharmawan/be-question-generator/utils"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	loginRequest := new(request.Login)

	if err := c.BodyParser(loginRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var user entity.User
	if err := db.Postgre.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Wrong Email or Password")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	if isPasswordMatch := utils.CheckPasswordHash(loginRequest.Password, user.Password); !isPasswordMatch {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Wrong Email or Password")
	}

	token, errGenerateToken := utils.GenerateJWTToken(user)

	if errGenerateToken != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errGenerateToken.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"token": token,
	})
}
