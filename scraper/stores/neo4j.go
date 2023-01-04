package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/gogo199432/bearchivedownloader/types"
	"github.com/lithammer/shortuuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
	"sync"
	"time"
)

type Neo4JStore struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
	doneCh chan bool
}

func (n4j *Neo4JStore) Init(url string, doneCh chan bool) {
	dbUri := url
	neo4jauth := neo4j.NoAuth()
	if viper.InConfig("database.password") {
		viper.SetDefault("database.username", "neo4j")
		neo4jauth = neo4j.BasicAuth(viper.GetString("database.username"), viper.GetString("database.password"), "")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4jauth, func(config *neo4j.Config) {
		config.ConnectionAcquisitionTimeout = 10 * time.Minute
		config.SocketConnectTimeout = 1 * time.Minute
	})
	if err != nil {
		panic(err)
	}
	n4j.driver = driver
	n4j.ctx = context.Background()
	n4j.doneCh = doneCh
}

func (n4j *Neo4JStore) CloseSession(session neo4j.SessionWithContext) {
	err := session.Close(n4j.ctx)
	if err != nil {
		fmt.Println("Unable to close session properly")
		return
	}
}

func (n4j *Neo4JStore) GetNodeCount() (count int64) {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer n4j.CloseSession(session)
	c, err := session.ExecuteRead(n4j.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(n4j.ctx, "Match (n) Return count(n) as count", nil)
		if err != nil {
			return 0, err
		}
		record, _ := result.Single(n4j.ctx)
		c, _ := record.Get("count")
		return c, nil
	})
	if err != nil {
		return 0
	}
	return c.(int64)
}

func (n4j *Neo4JStore) GetLeafs() (es []string, e error) {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer n4j.CloseSession(session)
	entriesUnparsed, err := session.Run(n4j.ctx, "MATCH (n) WHERE NOT (n)-->() RETURN n.ChildrenURLs as children", nil)
	if err != nil {
		return nil, err
	}
	var children []string
	for entriesUnparsed.Next(n4j.ctx) {
		n := entriesUnparsed.Record()
		childrenData, ok := n.Get("children")
		if !ok {
			return nil, errors.New("cannot get ChildrenURLs")
		}
		var localChildrenURLs map[string]string
		err = json.Unmarshal(childrenData.([]byte), &localChildrenURLs)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		children = append(children, maps.Values(localChildrenURLs)...)
	}
	return children, nil
}

func (n4j *Neo4JStore) Write(entry *Entry) error {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer n4j.CloseSession(session)
	childrenData, err := json.Marshal(entry.ChildrenURLs)
	if err != nil {
		return err
	}

	var tags string
	if len(entry.Tags) > 0 {
		tags = "SET n "
		for _, tag := range entry.Tags {
			tags += ":" + tag
		}
	}
	props := map[string]interface{}{
		"Url":          entry.Url,
		"Title":        entry.Title,
		"Text":         entry.Text,
		"Date":         entry.Date,
		"ChildrenURLs": childrenData,
		"Author":       entry.Author,
		"Id":           shortuuid.New(),
	}
	_, err = session.ExecuteWrite(n4j.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(n4j.ctx, "MERGE (n { Url: $Url}) ON CREATE "+tags+" SET n = $props  RETURN n", map[string]interface{}{
			"props": props,
			"Url":   entry.Url,
		})
	}, neo4j.WithTxTimeout(30*time.Second))
	if err != nil {
		fmt.Println("Error when creating entry for this page: " + entry.Url)
		return err
	}
	return nil
}

func (n4j *Neo4JStore) ResolveConnections() error {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	result, err := session.Run(n4j.ctx, "MATCH (n) WHERE NOT (n)-[]->() RETURN n.ChildrenURLs as children, n.Url as url", nil)
	if err != nil {
		return err
	}
	defer n4j.CloseSession(session)
	// Fingers crossed this isn't too big...
	entries, err := result.Collect(n4j.ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	queue := make(chan *neo4j.Record, len(entries))

	for _, entry := range entries {
		queue <- entry
	}
	close(queue)

	viper.SetDefault("scraper.connectionworkers", 10)
	for i := 0; i < viper.GetInt("scraper.connectionworkers"); i++ {
		wg.Add(1)
		go n4j.connectionWorker(&wg, queue)
	}

	wg.Wait()
	n4j.doneCh <- true
	return nil
}

func (n4j *Neo4JStore) connectionWorker(wg *sync.WaitGroup, queue <-chan *neo4j.Record) {
	defer wg.Done()
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer n4j.CloseSession(session)
	for entry := range queue {
		parentUrl, ok := entry.Get("url")
		if !ok {
			fmt.Println("Cannot get Url of parent")
			continue
		}
		childrenData, ok := entry.Get("children")
		if !ok {
			fmt.Println("Cannot get ChildrenURLs for entry", parentUrl)
			continue
		}
		var childrenURLs map[string]string
		err := json.Unmarshal(childrenData.([]byte), &childrenURLs)
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}

		for choice, child := range childrenURLs {
			_, err := session.ExecuteWrite(n4j.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
				_, err := tx.Run(n4j.ctx, "MATCH (parent), (child) WHERE child.Url = $child AND parent.Url = $url CREATE (parent)-[:Choice {text:$choice}]->(child)", map[string]any{
					"choice": choice,
					"child":  child,
					"url":    parentUrl,
				})
				if err != nil {
					return nil, err
				}
				return nil, nil
			})
			if err != nil {
				fmt.Println("Transaction error:", err)
				fmt.Println("Choice:", choice)
				fmt.Println("Child:", child)
				fmt.Println("Url:", parentUrl)
				continue
			}
		}
	}
}

func (n4j *Neo4JStore) Shutdown() {
	err := n4j.driver.Close(n4j.ctx)
	if err != nil {
		fmt.Println("ERROR WHILE SHUTTING DOWN DB DRIVER!")
		fmt.Println(err)
		return
	}
}
