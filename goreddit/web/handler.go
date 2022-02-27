package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/panta82/goreddit"
	"html/template"
	"net/http"
	"os"
)

type Handler struct {
	*chi.Mux
	store goreddit.Store
}

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewRouter(),
		store: store,
	}

	h.Use(middleware.Logger)
	h.Route("/threads", func(router chi.Router) {
		router.Get("/", h.ThreadsList())
		router.Get("/new", h.ThreadsCreate())
		router.Post("/", h.ThreadsStore())
	})

	return h
}

func guard(r *http.Request, err error) {
	if err != nil {
		err2 := fmt.Errorf("Error while handling %s %s: %w\n", r.Method, r.RequestURI, err)
		if err2 != nil {
			fmt.Print(os.Stderr, err2)
		}
	}
}

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.New("").Parse(`
		<h1>Threads</h1>
		<dl>
		{{range .Threads}}
			<dt><strong>{{.Title}}</strong></dt>
			<dd>{{.Description}}</dd>
		{{end}}
		</dl>
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		guard(r, tmpl.Execute(w, data{Threads: threads}))
	}
}

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(`
		<h1>Create Thread</h1>
		<form action="/threads" method="POST">
			<div>
				<label for="title">Title:</label>
				<input type="text" name="title" id="title">
			</div>

			<div>
				<label for="description">Content:</label>
				<input type="text" name="description" id="description">
			</div>
			
			<input type="submit" value="Create" />
		</form>
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		guard(r, tmpl.Execute(w, nil))
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
