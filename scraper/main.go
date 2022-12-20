package main

import (
	"fmt"
	"github.com/gocolly/colly"
	. "github.com/gogo199432/bearchivedownloader/stores"
	. "github.com/gogo199432/bearchivedownloader/types"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./scraper")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s\n", err)
	}
	// Initiate DB connection and crawler
	var store Storage
	switch viper.GetString("database.type") {
	case "neo4j":
		store = new(Neo4JStore)
	default:
		store = new(Neo4JStore)
	}
	store.Init(viper.GetString("database.connectionString"))
	defer store.Shutdown()
	c := colly.NewCollector(func(c *colly.Collector) {
		c.Async = true
	})
	if viper.InConfig("scraper.depth") {
		c.MaxDepth = viper.GetInt("scraper.depth")
	}

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
				entry.ChildrenURLs[h2.Text] = h.Request.AbsoluteURL(attr)
				h.Request.Visit(attr)
			}
		})

		err := store.Write(&entry)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
	fmt.Println(time.Now(), ": Starting crawl...")
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
	entry.Author = h.DOM.ChildrenFiltered("address").Text()

	tagHolder := h.DOM.Find("font").First()
	red, exists := tagHolder.Attr("color")
	if exists && red == "red" && len(entry.Tags) == 0 {
		tags := strings.Fields(tagHolder.Text())
		for _, tag := range tags {
			// Neo4J can't have dash (amongst others) in the label name
			tag = strings.ReplaceAll(tag, "-", "_")
			entry.Tags = append(entry.Tags, tag)
		}
	}

	return entry
}
