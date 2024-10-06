package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/controllers"
)

func SetupVerificationTokenRoutes(router fiber.Router) {
	router.Post("/verification-email/verify/:token", controllers.VerifyEmailToken)
	router.Post("/verification-email/resend", controllers.ResendVerificationToken)
}
