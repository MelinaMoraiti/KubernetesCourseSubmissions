package routes

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"todo-app/internal/health"
	"todo-app/internal/models"
)

func RegisterRoutes(todoBackendURL, storedImagesPath string) (http.Handler, error) {
	breakSwitch := &health.Switch{}

	// Liveness: only the process itself. A backend outage must not restart this pod.
	healthz, err := health.NewHandler(breakSwitch.Check())
	if err != nil {
		return nil, err
	}

	// Readiness: can this pod serve traffic? Add new dependencies here.
	readyz, err := health.NewHandler(
		breakSwitch.Check(),
		health.Backend(todoBackendURL),
	)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	// Probes stay outside the logging middleware.
	r.Handle("/healthz", healthz)
	r.Handle("/readyz", readyz)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Handle("/images/*",
			http.StripPrefix("/images/",
				http.FileServer(http.Dir(storedImagesPath)),
			),
		)
		r.Handle("/static/*",
			http.StripPrefix("/static/",
				http.FileServer(http.Dir("web/static")),
			),
		)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if breakSwitch.Broken() {
				renderFailure(w)
				return
			}
			indexHandler(w, r, todoBackendURL)
		})

		r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
			createTodoHandler(w, r, todoBackendURL)
		})

		r.Post("/break", func(w http.ResponseWriter, r *http.Request) {
			breakHandler(w, r, breakSwitch)
		})
		r.Patch("/todos/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
			completeTodoHandler(w, r, todoBackendURL)
		})
	})

	return r, nil
}
func completeTodoHandler(w http.ResponseWriter, r *http.Request, todoBackendURL string) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPatch,
		todoBackendURL+"/todos/"+strconv.Itoa(id)+"/complete",
		nil,
	)
	if err != nil {
		log.Printf("Could not build complete request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Could not contact todo backend: %v", err)
		http.Error(w, "Could not contact todo backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Todo API returned status: %s", resp.Status)
		http.Error(w, "Could not complete todo", resp.StatusCode)
		return
	}

	// Reload the page so the todo shows as completed, same as createTodoHandler.
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
// renderFailure renders the "system failure" page.
func renderFailure(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("web/templates/failure.html")
	if err != nil {
		log.Printf("Could not parse failure template: %v", err)
		http.Error(w, "System Failure", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Could not execute failure template: %v", err)
	}
}

func breakHandler(w http.ResponseWriter, r *http.Request, s *health.Switch) {
	s.Break()
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
func indexHandler(w http.ResponseWriter, r *http.Request, todoBackendURL string) {
	resp, err := http.Get(todoBackendURL + "/todos")
	if err != nil {
		log.Printf("Could not get all todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Todo API returned status: %s", resp.Status)
		http.Error(w, "Could not get todos", http.StatusInternalServerError)
		return
	}

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		log.Printf("Could not decode todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Printf("Could not parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, todos); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}

func createTodoHandler(w http.ResponseWriter, r *http.Request, todoBackendURL string) {
	form := url.Values{}
	form.Set("todo", r.FormValue("todo"))

	resp, err := http.PostForm(todoBackendURL+"/todos", form)
	if err != nil {
		http.Error(w, "Could not contact todo backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}