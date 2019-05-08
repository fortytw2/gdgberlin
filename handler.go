package gdgberlin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func NewHandler(db *sql.DB) http.Handler {

	m := http.NewServeMux()

	m.Handle("/", rootHandler(db))
	m.Handle("/fish", fishByFinCountHandler(db))

	return m
}

func rootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
		<html>
			<head>
				<title>E2E Testin'</title>
			</head>
			<body>
				<div id="root">
					<h1>Hello Gophers!</h1>
				</div>
			</body>
		</html>`))
	}
}

func fishByFinCountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		finCountStr := r.URL.Query().Get("fin_count")
		finCount, err := strconv.Atoi(finCountStr)
		if err != nil {
			http.Error(w, "fin_count query param is required", http.StatusBadRequest)
			return
		}

		rows, err := db.QueryContext(r.Context(), `
		SELECT 
			  id
			, name
			, water_type
		FROM fish 
		WHERE fin_count = $1;`, finCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type fish struct {
			ID        uuid.UUID
			Name      string
			WaterType string
		}

		var fishes []fish
		for rows.Next() {
			var f fish

			err = rows.Scan(&f.ID, &f.Name, &f.WaterType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fishes = append(fishes, f)
		}

		err = json.NewEncoder(w).Encode(fishes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
