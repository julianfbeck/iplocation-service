package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	ipdb "github.com/julianfbeck/ip-location-go-server/internal/ip-db"
)

func main() {
	db, err := ipdb.NewDB()
	if err != nil {
		fmt.Println(err)
	}

	// app.Listen(":3000")
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
			fmt.Println(err)
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

	http.ListenAndServe(":3000", nil)

}
