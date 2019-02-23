package example

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type AccountPayload struct {
	*Account
	Role string `json:"role"`
}

func NewAccountPayload(account *Account) *AccountPayload {
	return &AccountPayload{Account: account}
}

// Bind on AccountPayload will run after the unmarshalling is complete, its
// a good time to focus some post-processing after a decoding.
func (a *AccountPayload) Bind(r *http.Request) error {
	return nil
}

func (a *AccountPayload) Render(w http.ResponseWriter, r *http.Request) error {
	a.Role = "collaborator"
	return nil
}

// ExamplePayload is the request payload for Example data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type ExamplePayload struct {
	*Example
	Account *AccountPayload `json:"account,omitempty"`
}

func (p *ExamplePayload) Bind(r *http.Request) error {
	// e.Example is nil if no Example fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if p.Example == nil {
		return errors.New("missing required Example fields.")
	}

	// a.Account is nil if no Accountpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.Account or futher nested fields like a.Account.Name are accessed elsewhere.

	// just a post-process after a decode..
	p.Example.Title = strings.ToLower(p.Example.Title) // as an example, we down-case
	return nil
}

func NewExamplePayload(example *Example) *ExamplePayload {
	p := &ExamplePayload{Example: example}

	if p.Account == nil {
		if account, _ := dbGetAccount(p.AccountID); account != nil {
			p.Account = NewAccountPayload(account)
		}
	}

	return p
}

func (rd *ExamplePayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ExampleListPayload []*ExamplePayload

func NewExampleListPayload(examples []*Example) []render.Renderer {
	list := []render.Renderer{}
	for _, example := range examples {
		list = append(list, NewExamplePayload(example))
	}
	return list
}
