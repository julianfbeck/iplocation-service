package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	// db, err := ipdb.NewDB()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	// get correct ip even if behind proxy
	// 	ip := c.IP()
	// 	// X-Real-Ip or X-Forwarded-For headers
	// 	text := string(c.Request().Header.Header())
	// 	fmt.Println(text)
	// 	fmt.Println(string(c.Request().Header.Peek("X-Real-Ip")))
	// 	fmt.Println(string(c.Request().Header.Peek("X-Forwarded-For")))

	// 	// get location from ip address
	// 	location, err := db.LookUpIP(ip)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		// retrun error
	// 		return c.SendString(err.Error())
	// 	}
	// 	// return struct as json
	// 	return c.JSON(location)

	// })

	// app.Listen(":3000")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.Header.Get("X-Real-Ip")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Forwarded-For")
		}
		if clientIP == "" {
			clientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		fmt.Fprintf(w, "Client IP: %s", clientIP)
	})

	http.ListenAndServe(":3000", nil)

}
