package gdgberlin

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

var rootE2ECases = []e2eTestFn{
	{
		"renders-index-html",
		func(db *sql.DB, ctx context.Context) func(t *testing.T) {

			return func(t *testing.T) {
				// time.Sleep(20 * time.Second)

				var res string
				err := chromedp.Run(ctx, chromedp.Text(`#root`, &res, chromedp.NodeVisible, chromedp.ByID))
				if err != nil {
					t.Fatal(err)
				}

				if !strings.Contains(res, "Hello Gophers!") {
					t.Fatalf("did not see 'Hello Gophers!' on the page, saw '%s' instead", res)
				}
			}
		},
	},
}
