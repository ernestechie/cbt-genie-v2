package middleware

import (
	"fmt"
	"strings"

	"github.com/ernestechie/cbt-genie-v2/utils"
	"github.com/gofiber/fiber/v2"
)

// TODO: Pass a user role enum to protected route, to make it flexible and only restrict certain routes from certain users.
// That is, passing "Admin" as a prop to the ProtectRoute() handler would restrict the route to only allow admins, vice versa
func ProtectRoute () fiber.Handler {
	return func (c *fiber.Ctx) error  {
		authHeader := c.Get("Authorization")
		token := strings.Split(authHeader, "Bearer ")

		if authHeader == "" || len(token) == 0 {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Access denied! You are not authorized to perform this action",
			})
		}

		tokenString := token[1]
		if tokenString == "" || len(tokenString) == 0 {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Access denied! You are not authorized to perform this action",
			})
		}

		tokenClaims, err := utils.VerifyJwt(tokenString)
		if err != nil {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		fmt.Println("Token_claims", tokenClaims)

		return c.Next()
	}
}
