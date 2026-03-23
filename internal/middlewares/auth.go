package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gynScore-backend/internal/config"
	"gynScore-backend/pkg/utils"
)

// AuthMiddleware valida o token JWT presente no cabeçalho Authorization
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "Token de autenticação não fornecido")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.Error(c, fiber.StatusUnauthorized, "Formato do token inválido. Use: Bearer <token>")
		}

		claims, err := utils.ValidarToken(parts[1], cfg.JWTSecret)
		if err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "Token inválido ou expirado")
		}

		// Armazena os dados do usuário autenticado no contexto da requisição
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		return c.Next()
	}
}
