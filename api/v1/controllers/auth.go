package controllers

import (
	"time"

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

func RequestResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var resetReq ResetPasswordRequest

	if err := c.BodyParser(&resetReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request")
	}

	validate := validator.New()
	if err := validate.Struct(&resetReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var user entity.User
	if err := db.Postgre.Where("email = ?", resetReq.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	db.Postgre.Where("user_id = ? AND is_used = false", user.ID).Delete(&entity.VerificationToken{})

	resetToken := utils.GenerateVerificationToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	newResetToken := entity.VerificationToken{
		UserID:    user.ID,
		Token:     resetToken,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}

	if err := db.Postgre.Create(&newResetToken).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create reset password token")
	}

	utils.SendEmail([]string{user.Email}, "Verification Token for reset password", "Token : "+newResetToken.Token)

	return utils.SuccessResponse(c, "Password reset email has been sent")
}

func ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordInput struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	var resetInput ResetPasswordInput

	if err := c.BodyParser(&resetInput); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request")
	}

	validate := validator.New()
	if err := validate.Struct(&resetInput); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var resetToken entity.VerificationToken
	if err := db.Postgre.Where("token = ? AND is_used = false", resetInput.Token).First(&resetToken).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired token")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Token has expired")
	}

	hashedPassword, err := utils.HashPassword(resetInput.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash new password")
	}

	var user entity.User
	if err := db.Postgre.First(&user, resetToken.UserID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "User not found")
	}

	user.Password = hashedPassword
	db.Postgre.Save(&user)

	resetToken.IsUsed = true
	db.Postgre.Save(&resetToken)

	return utils.SuccessResponse(c, "Password has been reset successfully")
}
