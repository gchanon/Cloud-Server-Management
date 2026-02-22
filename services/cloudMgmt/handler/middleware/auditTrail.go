package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golf/cloudmgmt/services/cloudMgmt/behavior"
	"github.com/golf/cloudmgmt/services/cloudMgmt/model"
)

func AuditMiddleware(auditBehavior *behavior.AuditTrail, serverBehavior *behavior.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		var userID int64
		if uid, ok := c.Locals("user_id").(string); !ok || uid == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: user_id not found in context")
		}

		serverId, errGet := getServerIdFromPath(c, serverBehavior)
		if errGet != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Unable to extract server ID from path")
		}

		auditBehavior.Create(&model.AuditTrailModel{
			UserID:     userID,
			Action:     c.Method(),
			ServerID:   serverId,
			Path:       c.Path(),
			IPAddress:  c.IP(),
			ResStatus:  c.Response().StatusCode(),
			ActionTime: time.Now(),
			//OldValue
			//NewValue
		})

		return err

	}

}

func getServerIdFromPath(c *fiber.Ctx, serverBehavior *behavior.Server) (string, error) {
	if isServerGroupPath(c) {
		serverList, errGetAll := serverBehavior.GetAll()

		if errGetAll != nil {
			return "", nil
		}
		if strings.Contains(c.Path(), "power") {
			parts := strings.Split(c.Path()[2:], "/")
			if len(parts) == 2 {
				if parts[1] == "power" {
					return parts[0], nil
				}
			}
		}

		fmt.Println(serverList)
		// for serverId, serverData := range serverList {
		// 	if len(strings.Split(c.Path(), "/")) == 2 && c.Method() == fiber.MethodPost {

		// 		type SkuData struct {
		// 			SkuDetail string `json:"sku"`
		// 		}

		// 		var sku SkuData
		// 		if errParseBody := c.BodyParser(&sku); errParseBody != nil {
		// 			return "", fiber.NewError(fiber.ErrBadRequest.Code, "")
		// 		}

		// 		if sku.SkuDetail == serverData.Sku {
		// 			return
		// 		}
		// 	}
		// }
	}

	return "", nil
}

func isServerGroupPath(c *fiber.Ctx) bool {
	parts := strings.Split(c.Path()[1:], "/")

	if len(parts) >= 0 {
		return parts[0] == "server"
	}

	return false
}
