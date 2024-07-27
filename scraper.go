package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Race struct {
	Slug       string    `json:"slug"`
	Date       time.Time `json:"date"`
	Location   string    `json:"location"`
	Department string    `json:"department"`
	Category   string    `json:"category"`
	Link       string    `json:"link"`
	StartList  string    `json:"start_list,omitempty"`
}

func ScrapeBicycleRaces(app *pocketbase.PocketBase) error {
	log.Println("Starting to scrape bicycle races")
	url := "https://www.cif-ffc.fr/menuCif/dopublic/toutdosite.php"

	// Send HTTP GET request
	log.Printf("Sending GET request to %s", url)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching the webpage: %v", err)
		return err
	}
	defer response.Body.Close()

	log.Printf("Response status: %s", response.Status)

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}

	log.Printf("Response body length: %d bytes", len(body))

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Printf("Error loading HTTP response body: %v", err)
		return err
	}

	var routeRaces []Race
	var departRaces []Race

	document.Find("#ROUTE table").Each(func(index int, tablehtml *goquery.Selection) {
		log.Printf("Processing ROUTE table #%d", index+1)
		routeRaces = append(routeRaces, scrapeRouteTable(tablehtml)...)
	})

	document.Find("#DEPART table").Each(func(index int, tablehtml *goquery.Selection) {
		log.Printf("Processing DEPART table #%d", index+1)
		departRaces = append(departRaces, scrapeDepartTable(tablehtml)...)
	})

	log.Printf("Total ROUTE races found: %d", len(routeRaces))
	log.Printf("Total DEPART races found: %d", len(departRaces))

	races := append(routeRaces, departRaces...)
	return insertOrUpdateRaces(app, races)
}

func scrapeRouteTable(tablehtml *goquery.Selection) []Race {
	var races []Race
	tablehtml.Find("tr").Each(func(rowIndex int, rowhtml *goquery.Selection) {
		var race Race
		rowhtml.Find("td").Each(func(colIndex int, tablecell *goquery.Selection) {
			text := strings.TrimSpace(tablecell.Text())
			switch colIndex {
			case 0:
				parsedDate, err := parseFrenchDate(text)
				if err != nil {
					log.Printf("Error parsing date '%s': %v", text, err)
				} else {
					race.Date = parsedDate
				}
			case 1:
				race.Location = text
			case 2:
				race.Department = text
			case 3:
				race.Category = text
			case 5:
				if link, exists := tablecell.Find("a").Attr("href"); exists {
					race.Link = "https://www.cif-ffc.fr/menuCif/dopublic/" + link
				}
			case 6:
				if text != "" {
					if startList, exists := tablecell.Find("a").Attr("href"); exists {
						race.StartList = "https://www.cif-ffc.fr/menuCif/dopublic/" + startList
					}
				}
			}
		})

		// Generate and set the slug
		race.Slug = generateSlug(race.Date, race.Location, race.Category)

		races = append(races, race)
		log.Printf("Processed ROUTE race: %+v", race)
	})
	return races
}

func scrapeDepartTable(tablehtml *goquery.Selection) []Race {
	var races []Race
	tablehtml.Find("tr").Each(func(rowIndex int, rowhtml *goquery.Selection) {
		var race Race
		rowhtml.Find("td").Each(func(colIndex int, tablecell *goquery.Selection) {
			text := strings.TrimSpace(tablecell.Text())
			switch colIndex {
			case 0:
				parsedDate, err := parseFrenchDate(text)
				if err != nil {
					log.Printf("Error parsing date '%s': %v", text, err)
				} else {
					race.Date = parsedDate
				}
			case 1:
				race.Location = text
			case 2:
				race.Department = text
			case 3:
				race.Category = text
			case 4:
				if link, exists := tablecell.Find("a").Attr("href"); exists {
					race.Link = "https://www.cif-ffc.fr/menuCif/dopublic/" + link
				}
			case 5:
				if text != "" {
					if startList, exists := tablecell.Find("a").Attr("href"); exists {
						race.StartList = "https://www.cif-ffc.fr/menuCif/dopublic/" + startList
					}
				}
			}
		})

		// Generate and set the slug
		race.Slug = generateSlug(race.Date, race.Location, race.Category)

		races = append(races, race)
		log.Printf("Processed DEPART race: %+v", race)
	})
	return races
}

var frenchMonths = map[string]time.Month{
	"Janv.": time.January,
	"Fév.":  time.February,
	"Mars":  time.March,
	"Avr.":  time.April,
	"Mai":   time.May,
	"Juin":  time.June,
	"Juil.": time.July,
	"Août":  time.August,
	"Sept.": time.September,
	"Oct.":  time.October,
	"Nov.":  time.November,
	"Déc.":  time.December,
}

func parseFrenchDate(dateStr string) (time.Time, error) {
	parts := strings.Fields(dateStr)
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", parts[1])
	}

	month, ok := frenchMonths[parts[2]]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month: %s", parts[2])
	}

	year := time.Now().Year()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

func generateSlug(date time.Time, location, category string) string {
	// Normalize the strings (remove accents, lowercase, etc.)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	location, _, _ = transform.String(t, location)
	category, _, _ = transform.String(t, category)

	// Replace spaces with hyphens and remove any non-alphanumeric characters
	location = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, location)

	category = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, category)

	// Create the slug
	slug := fmt.Sprintf("%s-%s-%s", date.Format("2006-01-02"), location, category)

	return slug
}

func insertOrUpdateRaces(app *pocketbase.PocketBase, races []Race) error {
	collection, err := app.Dao().FindCollectionByNameOrId("races")
	if err != nil {
		return nil
	}

	for _, race := range races {
		// Check if a race with this slug already exists
		record, _ := app.Dao().FindFirstRecordByData("races", "slug", race.Slug)

		if record == nil {
			// Create a new record
			record = models.NewRecord(collection)
		}

		// Set or update the record fields
		record.Set("slug", race.Slug)
		record.Set("date", race.Date)
		record.Set("city", race.Location)
		record.Set("area", race.Department)
		record.Set("category", race.Category)
		record.Set("do_link", race.Link)
		record.Set("startlist_link", race.StartList)

		// Save the record
		if err := app.Dao().SaveRecord(record); err != nil {
			log.Printf("Error saving race %s: %v", race.Slug, err)
		}
	}

	return nil
}
