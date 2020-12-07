package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

/// Template
type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.template.Execute(w, nil)
}

func main() {

	defaultRoom := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", defaultRoom)

	// Start the room:
	defaultRoom.run()

	// Start the webserver
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
