package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/ernestechie/cbt-genie-v2/utils"
	"github.com/gofiber/fiber/v2"
)

// TODO: Pass a user role enum to protected route, to make it flexible and only restrict certain routes from certain users.
// That is, passing "Admin" as a prop to the ProtectRoute() handler would restrict the route to only allow admins, vice versa
func ProtectRoute () fiber.Handler {
	return func (c *fiber.Ctx) error  {
		authHeader := c.Get("Authorization")
		token := strings.Split(authHeader, "Bearer ")

		if len(token) > 2 || authHeader == "" || len(token) == 0 {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		tokenString := token[1]
		if tokenString == "" || len(tokenString) == 0 {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		fmt.Println("Token ", tokenString)

		var operationErrors []fiber.Map;

		jwtUser, err := utils.VerifyJwt(tokenString)
		if err != nil {
			operationErrors = append(operationErrors, fiber.Map{
				"message": err.Error(),
			})
		}	
		if time.Now().After(jwtUser.TokenExpiry) {
			operationErrors = append(operationErrors, fiber.Map{
				"message": "access denied! session expired",
			})
		}

		if len(operationErrors) > 0 {
			return  c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"errors": operationErrors,
			"message": "Invalid or expired token",
		})
		}
		
		fmt.Printf("Jwt_user %+v\n", jwtUser)

		return c.Next()
	}
}
