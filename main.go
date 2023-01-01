package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ipdb "github.com/julianfbeck/ip-location-go-server/internal/ip-db"
)

func main() {
	db, err := ipdb.NewDB()
	if err != nil {
		fmt.Println(err)
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// get correct ip even if behind proxy
		ip := c.IP()
		// X-Real-Ip or X-Forwarded-For headers

		fmt.Println(string(c.Request().Header.Peek("X-Real-Ip")))
		fmt.Println(string(c.Request().Header.Peek("X-Forwarded-For")))

		// get location from ip address
		location, err := db.LookUpIP(ip)
		if err != nil {
			fmt.Println(err)
			// retrun error
			return c.SendString(err.Error())
		}
		// return struct as json
		return c.JSON(location)

	})

	app.Listen(":3000")

}
