package main

import (
	"fmt"
	"github.com/gocolly/colly"
	. "github.com/gogo199432/bearchivedownloader/src/stores"
	. "github.com/gogo199432/bearchivedownloader/src/types"
	"strings"
	"time"
)

func main() {
	// Initiate DB connection and crawler
	var store Storage = new(Neo4JStore)
	store.Init("neo4j://localhost:7687")
	defer store.Shutdown()
	c := colly.NewCollector(func(c *colly.Collector) {
		c.MaxDepth = 1
		c.Async = true
	})

	// Try to grab any leaf nodes already in the DB
	leafs, err := store.GetLeafs()
	if err != nil {
		panic(err)
	}

	// Define what to do when processing an entry
	// (Creates an Entry object, fills it up and passes it to the DB to store)
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

	// Actual start of the crawling. If there were leafs, we want to continue there, otherwise start at the beginning
	if len(leafs) == 0 {
		c.Visit("https://addventure.bearchive.com/~addventure/game1/docs/000/2.html")
	} else {
		for _, leaf := range leafs {
			c.Visit(leaf)
		}
	}

	// Make sure we have finished crawling on all threads, and then connect new nodes to existing graph
	c.Wait()
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
