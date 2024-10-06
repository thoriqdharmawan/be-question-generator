package controllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/entity"
	"github.com/thoriqdharmawan/be-question-generator/db"
	"github.com/thoriqdharmawan/be-question-generator/utils"
)

func VerifyEmailToken(c *fiber.Ctx) error {
	token := c.Params("token")

	var verificationToken entity.VerificationToken
	if err := db.Postgre.Where("token = ? AND is_used = false", token).First(&verificationToken).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired token")
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Token has expired")
	}

	verificationToken.IsUsed = true
	db.Postgre.Save(&verificationToken)

	var user entity.User
	if err := db.Postgre.First(&user, verificationToken.UserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "User not found")
	}

	user.VerifiedAt = time.Now().Format("2006-01-02 15:04:05")
	db.Postgre.Save(&user)

	return utils.SuccessResponse(c, "Email verified successfully")
}

func ResendVerificationToken(c *fiber.Ctx) error {
	type ResendRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var resendReq ResendRequest

	if err := c.BodyParser(&resendReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request")
	}

	validate := validator.New()
	if err := validate.Struct(&resendReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var user entity.User
	if err := db.Postgre.Where("email = ?", resendReq.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	if user.VerifiedAt != "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User already verified")
	}

	db.Postgre.Where("user_id = ? AND is_used = false", user.ID).Delete(&entity.VerificationToken{})

	verificationToken := utils.GenerateVerificationToken()
	expiresAt := time.Now().Add(15 * time.Minute)

	newToken := entity.VerificationToken{
		UserID:    user.ID,
		Token:     verificationToken,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}

	if err := db.Postgre.Create(&newToken).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create new verification token")
	}

	utils.SendEmail([]string{user.Email}, "Verification Token for account creation", "Token : "+newToken.Token)

	return utils.SuccessResponse(c, "Verification email has been resent")
}
