package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

const ROLE_LIST = "admin,treasury,lawyer,secretary"

func AllowRoles(roles string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			newCtx := context.WithValue(r.Context(), "roles", roles)
			next.ServeHTTP(w, r.WithContext(newCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func AuthorizeRoles(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("Role")
		allowedRoles := r.Context().Value("roles").(string)
		resource := r.URL.RequestURI()
		method := r.Method

		if role == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("ERROR - There is no role assigned\n"))
			return
		}

		if allowedRoles == "all" {
			valid := strings.Contains(ROLE_LIST, role)
			if !valid {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("ERROR - The current role (" + role + ") is not supported\n"))
				return
			}

			next.ServeHTTP(w, r)
		} else {
			allowed := strings.Contains(allowedRoles, role)

			if allowed {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("ERROR - The current role (" + role + ") is not allowed to execute " + resource + " [" + method + "]\n"))
				return
			}
		}
	}

	return http.HandlerFunc(fn)
}

func setRoutes(r *chi.Mux) {
	r.With(AllowRoles("all")).With(AuthorizeRoles).Get("/", func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("Role")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello World! - " + role + "</h1>\n"))
	})

	r.With(AllowRoles("admin")).With(AuthorizeRoles).Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>CREATED! </h1>\n"))
	})

	r.With(AllowRoles("admin,treasury,lawyer")).With(AuthorizeRoles).Get("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("Role")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>IT WORKS! - Using this as role " + role + "</h1>\n"))
	})

	r.With(AllowRoles("admin,lawyer")).With(AuthorizeRoles).Post("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>CREATED!</h1>\n"))
	})

	r.With(AllowRoles("admin,treasury")).With(AuthorizeRoles).Put("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>UPDATED!</h1>\n"))
	})

	r.With(AllowRoles("admin,treasury")).With(AuthorizeRoles).Patch("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>PATCHED!</h1>\n"))
	})

	r.Get("/free-resource", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>No need of roles here!!</h1>\n"))
	})
}

func main() {
	port := "8080"

	log.Printf("Starting up on http://localhost:%s", port)

	router := chi.NewRouter()

	setRoutes(router)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
