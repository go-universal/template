package main

import (
	"fmt"
	"net/http"

	"github.com/go-universal/fs"
	"github.com/go-universal/template"
)

func main() {
	// Initialize tempalate
	fs := fs.NewDir("./assets")
	tpl := template.New(
		fs,
		template.WithRoot("views"),
		template.WithPartials("views/partials"),
	)

	if err := tpl.Load(); err != nil {
		fmt.Println(err)
		return
	}

	// Handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.Render(w, "pages/home", nil, "layout"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.Render(w, "pages/contacts", nil, "layout", "pages/contact/form", "pages/contact/social"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.Render(w, "errors", nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}

}
