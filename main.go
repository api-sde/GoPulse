package main

import (
	"fmt"
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

	fmt.Println("Launching..")

	defaultRoom := newRoom()
	fmt.Println("new..")

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", defaultRoom)

	// Start the room:
	go defaultRoom.run()
	fmt.Println("Start the room...")

	// Start the webserver
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	fmt.Println("Ready to serve..")

}
