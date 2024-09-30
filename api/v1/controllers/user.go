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

func GetUsers(c *fiber.Ctx) error {
	var users []entity.User

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	if err := db.Postgre.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "error query:"+err.Error())
	}

	var total int64
	if err := db.Postgre.Model(&entity.User{}).Count(&total).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "error get count"+err.Error())
	}

	meta := utils.GenerateMetaData(total, limit, offset)

	return utils.SuccessResponseWithMeta(c, users, meta)
}

func CreateUser(c *fiber.Ctx) error {
	user := new(request.CreateUser)

	if err := c.BodyParser(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var userIsEmailExists entity.User
	if result := db.Postgre.Where("email = ?", user.Email).First(&userIsEmailExists); result.Error == nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	newUser := entity.User{
		Name:       user.Name,
		Email:      user.Email,
		Password:   hashedPassword,
		VerifiedAt: "",
	}

	if err := db.Postgre.Create(&newUser).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return utils.SuccessResponse(c, newUser)
}

func GetUserById(c *fiber.Ctx) error {
	userId := c.Params("id")

	var user entity.User

	if err := db.Postgre.Where("id = ?", userId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return utils.SuccessResponse(c, user)
}
