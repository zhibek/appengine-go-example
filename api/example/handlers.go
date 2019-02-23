package example

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/zhibek/appengine-go-example/util/errors"
)

// RouteIdParser middleware is used to load an Example object from
// the URL parameters passed through as the request. In case
// the Example could not be found, we stop here and return a 404.
func RouteIdParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var example *Example
		var err error

		if exampleID := chi.URLParam(r, "exampleId"); exampleID != "" {
			example, err = dbGetExample(exampleID)
		} else {
			render.Render(w, r, errors.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, errors.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "example", example)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewExampleListPayload(examples)); err != nil {
		render.Render(w, r, errors.ErrRender(err))
		return
	}
}

// CreateExample persists the posted Example and returns it
// back to the client as an acknowledgement.
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	data := &ExamplePayload{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}

	example := data.Example
	dbNewExample(example)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewExamplePayload(example))
}

// GetHandler returns the specific Example. You'll notice it just
// fetches the Example right off the context, as its understood that
// if we made it this far, the Example must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the example
	// context because this handler is a child of the RouteIdParser
	// middleware. The worst case, the recoverer middleware will save us.
	example := r.Context().Value("example").(*Example)

	if err := render.Render(w, r, NewExamplePayload(example)); err != nil {
		render.Render(w, r, errors.ErrRender(err))
		return
	}
}

// UpdateHandler updates an existing Example in our persistent store.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	example := r.Context().Value("example").(*Example)

	data := &ExamplePayload{Example: example}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}
	example = data.Example
	dbUpdateExample(example.ID, example)

	render.Render(w, r, NewExamplePayload(example))
}

// DeleteHandler removes an existing Example from our persistent store.
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the example
	// context because this handler is a child of the RouteIdParser
	// middleware. The worst case, the recoverer middleware will save us.
	example := r.Context().Value("example").(*Example)

	example, err = dbRemoveExample(example.ID)
	if err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewExamplePayload(example))
}
