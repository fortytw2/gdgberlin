package gdgberlin

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/fortytw2/dockertest"
)

type testFn struct {
	name string
	do   func(*sql.DB) func(t *testing.T)
}

type e2eTestFn struct {
	name string
	do   func(*sql.DB, context.Context) func(t *testing.T)
}

const pgVersion = "postgres:11-alpine"

func TestDemoApp(t *testing.T) {
	container, err := dockertest.RunContainer(pgVersion, "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			log.Println(err)
			return err
		}

		return db.Ping()
	}, "--tmpfs", "/var/lib/postgresql")
	if err != nil {
		t.Fatal(err)
	}
	defer container.Shutdown()

	db, err := NewDB("postgres://postgres:postgres@" + container.Addr + "?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string][]testFn{
		"fish": fishCases,
	}

	for name, cases := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, c := range cases {
				t.Run(c.name, c.do(db))
				TruncateTables(db, t.Fatal)
			}
		})
	}

	// end to end cases

	e2eCases := map[string][]e2eTestFn{
		// "fish": fishE2ECases,
		"root": rootE2ECases,
	}

	handler := NewHandler(db)
	localPort := freePort()

	go http.ListenAndServe(fmt.Sprintf(":%d", localPort), handler)

	dir, err := ioutil.TempDir("", "test-chromedp")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(t.Logf))
	defer cancel()

	for name, cases := range e2eCases {
		t.Run(name, func(t *testing.T) {
			for _, c := range cases {
				err = chromedp.Run(ctx, chromedp.Navigate(fmt.Sprintf("http://localhost:%d", localPort)))
				if err != nil {
					t.Fatal(err)
				}

				t.Run(c.name, c.do(db, ctx))
				TruncateTables(db, t.Fatal)
			}
		})
	}
}

func TruncateTables(db *sql.DB, f func(...interface{})) {
	// https://stackoverflow.com/a/12082038
	_, err := db.Exec(`
	DO
	$func$
	BEGIN
	EXECUTE (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
		FROM pg_class
		WHERE relkind = 'r'  -- only tables
		AND relnamespace = 'public'::regnamespace);
	END
	$func$;`)
	if err != nil {
		f(err)
	}
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
