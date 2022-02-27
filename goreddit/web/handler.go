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
		router.Post("/{id}/delete", h.ThreadsDelete())
	})

	h.Get("/html", func(writer http.ResponseWriter, request *http.Request) {
		t := template.Must(template.ParseFiles("web/templates/layout.gohtml"))
		type params struct {
			Title   string
			Text    string
			Lines   []string
			Number1 int
			Number2 int
		}
		t.Execute(writer, params{
			Title:   "Hello",
			Text:    "World",
			Lines:   []string{"a", "b", "c"},
			Number1: 1,
			Number2: 4,
		})
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
			<dt><strong>{{.Title}}</strong> <form method="post" action="/threads/{{.ID}}/delete"><button type="submit">Delete</button></form> </dt>
			<dd>
				<div>{{.Description}}</div>
			</dd>
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

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.store.DeleteThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
