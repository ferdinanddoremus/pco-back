package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// Setup cron job
	scheduler := cron.New()
	scheduler.MustAdd("scrape-races", "0 * * * *", func() { // Run every hour
		log.Println("Starting bicycle race scraping job")
		err := ScrapeBicycleRaces(app)
		if err != nil {
			log.Printf("Error scraping bicycle races: %v", err)
		} else {
			log.Println("Bicycle race scraping job completed successfully")
		}
	})

	// Start the scheduler
	scheduler.Start()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/scrap",
			Handler: func(c echo.Context) error {
				err := ScrapeBicycleRaces(app)

				if err != nil {
					log.Printf("Error scraping data: %v", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"status":  "error",
						"message": "Error scraping data: " + err.Error(),
					})
				}

				log.Printf("Scrape completed successfully")

				return c.JSON(http.StatusOK, map[string]string{
					"status":  "success",
					"message": "Scraping completed",
				})
			},
			Middlewares: []echo.MiddlewareFunc{
				// Add any middleware you need here
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
