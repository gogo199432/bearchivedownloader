package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	. "github.com/gogo199432/bearchivedownloader/src/stores"
	. "github.com/gogo199432/bearchivedownloader/src/types"
)

func main() {

	var store Storage = new(Neo4JStore)
	store.Init("neo4j://localhost:7687")
	defer store.Shutdown()
	c := colly.NewCollector(func(c *colly.Collector) {
		c.MaxDepth = 5
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {
		entry := createEntry(h)
		h.ForEach("li", func(i int, h2 *colly.HTMLElement) {
			if strings.Contains(h2.Text, "*") {
				attr := h2.ChildAttr("a[href]", "href")
				fmt.Println("Visiting ", attr)
				entry.ChildrenURLs[h2.Text] = h.Request.AbsoluteURL(attr)
				h.Request.Visit(attr)
			}
		})

		store.Write(&entry)
	})

	c.Visit("https://addventure.bearchive.com/~addventure/game1/docs/000/2.html")

	fmt.Println(time.Now(), ": Finished crawling.")
	fmt.Println("Starting connecting")

	store.ResolveConnections()

	fmt.Println(time.Now(), ": Finished :)")
}

func createEntry(h *colly.HTMLElement) Entry {
	var entry = Entry{
		Url: h.Request.URL.String(),
	}
	h.ForEach("p", func(i int, parag *colly.HTMLElement) {
		txt := strings.ReplaceAll(parag.Text, "\x0a", "")
		published, err := time.Parse("Mon Jan 2 15:04:05 2006", txt)
		if err != nil {
			if strings.Contains(txt, "Edit Tags") ||
				strings.Contains(txt, "Add comment") ||
				strings.Contains(txt, "comments") ||
				strings.Contains(txt, "Last updated") ||
				txt == "" {
				return
			}
			entry.Text += txt
		} else {
			entry.Date = published
		}
	})
	entry.Title = h.DOM.ChildrenFiltered("h1").Text()
	if entry.Title == "" {
		entry.Title = h.DOM.ChildrenFiltered("h2").Text()
	}
	entry.ChildrenURLs = make(map[string]string)
	entry.Children = make(map[string]*Entry)
	return entry
}
