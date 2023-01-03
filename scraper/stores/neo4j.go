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
)

type Neo4JStore struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
}

func (n4j *Neo4JStore) Init(url string) {
	dbUri := url
	neo4jauth := neo4j.NoAuth()
	if viper.InConfig("database.password") {
		viper.SetDefault("database.username", "neo4j")
		neo4jauth = neo4j.BasicAuth(viper.GetString("database.username"), viper.GetString("database.password"), "")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4jauth)
	if err != nil {
		panic(err)
	}
	n4j.driver = driver
	n4j.ctx = context.Background()
}

func (n4j *Neo4JStore) GetNodeCount() (count int64) {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer session.Close(n4j.ctx)
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
	defer session.Close(n4j.ctx)
	entriesUnparsed, err := session.Run(n4j.ctx, "MATCH (n) WHERE NOT (n)-->() RETURN n.ChildrenURLs as children", nil)
	if err != nil {
		return nil, err
	}
	var children []string
	for entriesUnparsed.Next(n4j.ctx) {
		n := entriesUnparsed.Record()
		childrenData, ok := n.Get("children")
		if !ok {
			return nil, errors.New("Cannot get ChildrenURLs")
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
	defer session.Close(n4j.ctx)
	tx, err := session.BeginTransaction(n4j.ctx)
	if err != nil {
		return err
	}
	defer tx.Close(n4j.ctx)
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

	_, err = tx.Run(n4j.ctx, "CREATE (n:Entry { Id: $id,Url: $url, Title: $title, Text: $text, Date: $date, Author: $author}) "+tags+" SET n.ChildrenURLs = $childrenURLs  RETURN n", map[string]interface{}{
		"url":          entry.Url,
		"title":        entry.Title,
		"text":         entry.Text,
		"date":         entry.Date,
		"childrenURLs": childrenData,
		"author":       entry.Author,
		"id":           shortuuid.New(),
	})
	if err != nil {
		fmt.Println("Error when parsing this page: " + entry.Url)
		return err
	}
	err = tx.Commit(n4j.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (n4j *Neo4JStore) ResolveConnections() error {
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	result, err := session.Run(n4j.ctx, "MATCH (n) WHERE NOT (n)-[]->() RETURN n.ChildrenURLs as children, n.Title as title", nil)
	if err != nil {
		return err
	}
	session.Close(n4j.ctx)
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
	return nil
}

func (n4j *Neo4JStore) connectionWorker(wg *sync.WaitGroup, queue <-chan *neo4j.Record) {
	defer wg.Done()
	session := n4j.driver.NewSession(n4j.ctx, neo4j.SessionConfig{})
	defer session.Close(n4j.ctx)
	for entry := range queue {
		parentTitle, ok := entry.Get("title")
		if !ok {
			fmt.Println("Cannot get Title")
			continue
		}
		childrenData, ok := entry.Get("children")
		if !ok {
			fmt.Println("Cannot get ChildrenURLs for entry", parentTitle)
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
				_, err := tx.Run(n4j.ctx, "MATCH (parent), (child) WHERE child.Url = $child AND parent.Title = $title CREATE (parent)-[:Choice {text:$choice}]->(child)", map[string]any{
					"choice": choice,
					"child":  child,
					"title":  parentTitle,
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
				fmt.Println("Title:", parentTitle)
				continue
			}
		}
	}
}

func (n4j *Neo4JStore) Shutdown() {
	n4j.driver.Close(n4j.ctx)
}
