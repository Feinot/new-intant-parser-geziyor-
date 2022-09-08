package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://e.intant.ru/catalog/hardware"},
		ParseFunc: parseMovies,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}

func parseMovies(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.wrapper").Each(func(i int, s *goquery.Selection) {
		var sessions = strings.Split(s.Find("span.catalog__title").Text(), " \n ")
		sessions = sessions[:len(sessions)-1]

		for i := 0; i < len(sessions); i++ {
			sessions[i] = strings.Trim(sessions[i], "\n ")
			fmt.Println(sessions[i])
		}

		var description string

		if href, ok := s.Find("a.catalog__link").Attr("href"); ok {
			g.Get(r.JoinURL("https://e.intant.ru"+href), func(_g *geziyor.Geziyor, _r *client.Response) {
				description = _r.HTMLDoc.Find("div.catalog__link span.catalog__name").Text()

				description = strings.ReplaceAll(description, "BOX", "")
				description = strings.ReplaceAll(description, "OEM", "")
				description = strings.ReplaceAll(description, "Товар дня", "")
				description = strings.ReplaceAll(description, "\n", "")
				description = strings.TrimSpace(description)

				g.Exports <- map[string]interface{}{
					"title":       strings.ReplaceAll(_r.HTMLDoc.Find("a.controls__add-basket").Text(), "\u00A0", ""),
					"subtitle":    strings.TrimSpace(_r.HTMLDoc.Find("span.ico_addbasket").Text()),
					"sessions":    sessions,
					"description": description,
				}
			})
		}
	})
}
