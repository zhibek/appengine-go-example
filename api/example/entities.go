package example

import (
	"errors"
	"fmt"
	"math/rand"
)

// Account data model
type Account struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Example data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Example struct {
	ID     string `json:"id"`
	AccountID int64  `json:"account_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

// Example fixture data
var examples = []*Example{
	{ID: "1", AccountID: 100, Title: "Hi", Slug: "hi"},
	{ID: "2", AccountID: 200, Title: "sup", Slug: "sup"},
	{ID: "3", AccountID: 300, Title: "alo", Slug: "alo"},
	{ID: "4", AccountID: 400, Title: "bonjour", Slug: "bonjour"},
	{ID: "5", AccountID: 500, Title: "whats up", Slug: "whats-up"},
}

// Account fixture data
var accounts = []*Account{
	{ID: 100, Name: "London"},
	{ID: 200, Name: "Cairo"},
}

func dbNewExample(example *Example) (string, error) {
	example.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	examples = append(examples, example)
	return example.ID, nil
}

func dbGetExample(id string) (*Example, error) {
	for _, a := range examples {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("example not found.")
}

func dbUpdateExample(id string, example *Example) (*Example, error) {
	for i, a := range examples {
		if a.ID == id {
			examples[i] = example
			return example, nil
		}
	}
	return nil, errors.New("example not found.")
}

func dbRemoveExample(id string) (*Example, error) {
	for i, a := range examples {
		if a.ID == id {
			examples = append((examples)[:i], (examples)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("example not found.")
}

func dbGetAccount(id int64) (*Account, error) {
	for _, u := range accounts {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("account not found.")
}
