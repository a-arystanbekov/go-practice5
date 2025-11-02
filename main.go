package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	TotalOrders int    `json:"total_orders"`
}

func main() {
	db, err := sql.Open("postgres", "user=user password=password dbname=mydatabase sslmode=disable host=localhost port=5430")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		city := r.URL.Query().Get("city")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)
		if limit == 0 {
			limit = 10
		}

		query := `
			SELECT 
				u.id, u.name, u.city, COUNT(o.id) AS total_orders
			FROM users u
			LEFT JOIN orders o ON u.id = o.user_id
			WHERE ($1 = '' OR u.city = $1)
			GROUP BY u.id
			ORDER BY total_orders DESC, u.id DESC
			LIMIT $2 OFFSET $3
		`

		rows, err := db.Query(query, city, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.City, &u.TotalOrders); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		duration := time.Since(start)
		w.Header().Set("X-Query-Time", duration.String())
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(users)
	})

	fmt.Println("🚀 Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
