package gdgberlin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var fishCases = []testFn{
	{
		"basic-no-chromedp",
		func(db *sql.DB) func(t *testing.T) {
			return func(t *testing.T) {
				_, err := db.Exec(`
				INSERT INTO fish (name, fin_count, water_type)
				VALUES
				('shark', 4, 'SALT'::water_type),
				('tuna', 3, 'SALT'::water_type),
				('fishy fish', 2, 'FRESH'::water_type)`)
				if err != nil {
					t.Fatal(err)
				}

				handler := fishByFinCountHandler(db)

				r := httptest.NewRequest("GET", "/fish?fin_count=2", nil)
				w := httptest.NewRecorder()

				handler(w, r)

				statusCode := w.Result().StatusCode
				if statusCode != http.StatusOK {
					t.Fatal("did not get a 200, got", statusCode)
				}

				var jsonResp []interface{}
				err = json.NewDecoder(w.Result().Body).Decode(&jsonResp)
				if err != nil {
					t.Fatal(err)
				}

				if len(jsonResp) != 3 {
					t.Fatal("did not get back one fish")
				}
			}
		},
	},
}
