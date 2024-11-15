package main

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Note struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteServer struct {
	notes map[string]*Note
	mutex sync.Mutex
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func NewNoteServer() *NoteServer {
	return &NoteServer{
		notes: make(map[string]*Note),
	}
}

func (ns *NoteServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/create.html")
		return
	}
	content := r.FormValue("content")
	id := generateID()

	ns.mutex.Lock()
	ns.notes[id] = &Note{ID: id, Content: content}
	ns.mutex.Unlock()

	http.Redirect(w, r, "/note/"+id, http.StatusSeeOther)
}

func (ns *NoteServer) handleView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/note/"):]

	ns.mutex.Lock()
	note, exists := ns.notes[id]
	ns.mutex.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		ns.mutex.Lock()
		note.Content = content
		ns.mutex.Unlock()
	}

	tmpl := template.Must(template.ParseFiles("templates/note.html"))
	tmpl.Execute(w, note)
}

func main() {
	ns := NewNoteServer()

	http.HandleFunc("/", ns.handleCreate)
	http.HandleFunc("/note/", ns.handleView)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
