package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	ipdb "github.com/julianfbeck/ip-location-go-server/internal/ip-db"
	"github.com/robfig/cron/v3"
)

func main() {
	db, err := ipdb.NewDB()
	if err != nil {
		log.Println(err)
	}
	c := cron.New()

	// c.AddFunc("0 0 0 3 * *", func() { db.UpdateDB() })
	c.AddFunc("*/1 * * * *", func() { db.UpdateDB() })
	c.Start()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.Header.Get("X-Real-Ip")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Forwarded-For")
		}
		if clientIP == "" {
			clientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		location, err := db.LookUpIP(clientIP)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		pJSON, err := json.Marshal(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// Write the JSON string to the response.
		fmt.Fprint(w, string(pJSON))
	})
	log.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", nil)

}
