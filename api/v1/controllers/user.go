package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/entity"
	"github.com/thoriqdharmawan/be-question-generator/db"
	"github.com/thoriqdharmawan/be-question-generator/utils"
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
