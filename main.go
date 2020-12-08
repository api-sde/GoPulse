package main

import (
	"flag"
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

	t.template.Execute(w, r) // pass the request as data to the template
}

func main() {

	fmt.Println("Launching.")

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	defaultRoom := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", defaultRoom)

	fmt.Println("Start the room..")
	go defaultRoom.run()

	fmt.Println("Start the web server on...", *addr)

	// Start the webserver
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
