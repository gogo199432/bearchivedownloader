package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/gogo199432/bearchivedownloader/src/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4JStore struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
}

func (n4j *Neo4JStore) Init(url string) {
	dbUri := url
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.NoAuth())
	if err != nil {
		panic(err)
	}
	n4j.driver = driver
	n4j.ctx = context.Background()
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
	_, err = tx.Run(n4j.ctx, "CREATE (n:Entry { Url: $url, Title: $title, Text: $text, Date: $date}) SET n.ChildrenURLs = $childrenURLs  RETURN n", map[string]interface{}{
		"url":          entry.Url,
		"title":        entry.Title,
		"text":         entry.Text,
		"date":         entry.Date,
		"childrenURLs": childrenData,
	})
	if err != nil {
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
	defer session.Close(n4j.ctx)
	entries, err := session.Run(n4j.ctx, "MATCH (n) WHERE NOT (n)-[]->() RETURN n.ChildrenURLs as children, n.Title as title", nil)
	if err != nil {
		return err
	}
	tx, err := session.BeginTransaction(n4j.ctx)
	if err != nil {
		return err
	}
	defer tx.Close(n4j.ctx)

	for entries.Next(n4j.ctx) {
		n := entries.Record()
		childrenData, ok := n.Get("children")
		if !ok {
			return errors.New("Cannot get ChildrenURLs")
		}
		var childrenURLs map[string]string
		err = json.Unmarshal(childrenData.([]byte), &childrenURLs)
		if err != nil {
			fmt.Println(err)
			return err
		}
		parentTitle, ok := n.Get("title")
		if !ok {
			return errors.New("Cannot get Title")
		}
		for choice, child := range childrenURLs {
			tx.Run(n4j.ctx, "MATCH (parent), (child) WHERE child.Url = $child AND parent.Title = $title CREATE (parent)-[:Choice {text:$choice}]->(child)", map[string]any{
				"choice": choice,
				"child":  child,
				"title":  parentTitle,
			})
		}
	}
	tx.Commit(n4j.ctx)
	return nil
}

func (n4j *Neo4JStore) Shutdown() {
	n4j.driver.Close(n4j.ctx)
}
