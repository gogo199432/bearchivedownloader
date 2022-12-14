package main

import (
	"fmt"
	"github.com/gocolly/colly"
	. "github.com/gogo199432/bearchivedownloader/stores"
	. "github.com/gogo199432/bearchivedownloader/types"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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
	done := make(chan bool, 1)
	store.Init(viper.GetString("database.connectionString"), done)
	defer store.Shutdown()

	// Make sure we can handle docker stop and kill commands gracefully
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-signalChan
		log.Printf("Captured SIG: %v. Shutting down...", sig)
		store.Shutdown()
		done <- true
	}()
	var depth = 0
	if viper.InConfig("scraper.depth") {
		depth = viper.GetInt("scraper.depth")
	}
	if depth > 0 {
		fmt.Println("Depth is", strconv.Itoa(depth))
		crawl(store, depth)
		fmt.Println(time.Now(), ": Finished crawling. Starting connecting now...")
	} else {
		fmt.Println(time.Now(), ": Skipped crawling because depth was 0 or less. Starting connecting now...")
	}

	store.ResolveConnections()
	// Wait for done notification
	<-done
	fmt.Println(time.Now(), ": Finished :)")
}

func crawl(store Storage, depth int) {
	c := colly.NewCollector(func(c *colly.Collector) {
		c.Async = true
		c.DetectCharset = true
		c.MaxDepth = depth
	})

	// Limit the maximum parallelism
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	viper.SetDefault("scraper.parallelism", 2)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: viper.GetInt("scraper.parallelism")})

	// Try to grab any leaf nodes already in the DB
	leafs, err := store.GetLeafs()
	if err != nil {
		panic(err)
	}
	var addedEntries = 0
	// Define what to do when processing an entry
	// (Creates an Entry object, fills it up and passes it to the DB to store)
	c.OnHTML("body", func(h *colly.HTMLElement) {
		entry := createEntry(h)
		h.ForEach("li", func(i int, h2 *colly.HTMLElement) {
			if strings.Contains(h2.Text, "*") ||
				(strings.Contains(h2.Text, ">") && strings.Index(h2.Text, ">") == 0) {
				attr := h2.ChildAttr("a[href]", "href")
				choiceText := h2.Text
				// Since there might be duplicate choice texts (why wouldnt there be...)
				_, exists := entry.ChildrenURLs[choiceText]
				if exists {
					choiceText = choiceText + " - " + strconv.Itoa(i)
				}
				entry.ChildrenURLs[choiceText] = h.Request.AbsoluteURL(attr)
				h.Request.Visit(attr)
			}
		})

		err := store.Write(&entry)
		if err != nil {
			fmt.Println(err)
			return
		}
		addedEntries++
	})
	fmt.Println(time.Now(), ": Starting crawl...")
	// Actual start of the crawling. If there were leafs, we want to continue there, otherwise start at the beginning
	if store.GetNodeCount() == 0 {
		fmt.Println("No nodes present. Using root")
		c.Visit("https://addventure.bearchive.com/~addventure/game1/docs/000/2.html")
	} else {
		fmt.Println("Found leafs:", len(leafs))
		for _, leaf := range leafs {
			c.Visit(leaf)
		}
	}

	// Make sure we have finished crawling on all threads, and then connect new nodes to existing graph
	c.Wait()
	fmt.Println("Added entries count:", addedEntries)
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
				strings.Contains(txt, "Linking Enabled") ||
				txt == "" {
				return
			}
			entry.Text += txt + "<br>"
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
	tagText := tagHolder.SiblingsFiltered("b")
	if tagHolder.Index() != -1 && tagText.Text() == "Tags:" && tagHolder.Index() == tagText.Index()+1 {
		red, exists := tagHolder.Attr("color")
		if exists && red == "red" && len(entry.Tags) == 0 {
			tags := strings.Fields(tagHolder.Text())
			for _, tag := range tags {
				// Neo4J can't have dash (amongst others) in the label name
				tag = strings.ReplaceAll(tag, "-", "_")
				entry.Tags = append(entry.Tags, tag)
			}
		}
	}
	return entry
}
