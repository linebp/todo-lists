package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type ToDo struct {
	Title string
	List  []string `json:",omitempty"`
}

var todos map[string][]string = make(map[string][]string)

func exportHandler(w http.ResponseWriter, r *http.Request) {

	result := make([]ToDo, 0, len(todos))

	for title, list := range todos {
		todo := ToDo{Title: title, List: list}
		result = append(result, todo)
	}

	json.NewEncoder(w).Encode(result)
}

func importHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		text := r.FormValue("text")

		var data []ToDo

		if err := json.Unmarshal([]byte(text), &data); err != nil {
			type ImportError struct {
				Text  string
				Error string
			}
			t, _ := template.ParseFiles("import.html")
			t.Execute(w, ImportError{Text: text, Error: err.Error()})
			return
		} else {
			for _, todo := range data {
				todos[todo.Title] = todo.List
			}
		}

		http.Redirect(w, r, "/view/", http.StatusFound)
	default:
		t, _ := template.ParseFiles("import.html")
		t.Execute(w, nil)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, todos)
}

func editHandler(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Query().Get("title")

	t, _ := template.ParseFiles("edit.html")

	if list, ok := todos[title]; ok {
		t.Execute(w, ToDo{title, list})
	} else {
		t.Execute(w, nil)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	list := r.FormValue("list")
	title := r.FormValue("title")

	if len(title) > 0 {
		todos[title] = strings.Fields(list)
		http.Redirect(w, r, "/view/", http.StatusFound)
	} else {
		t, _ := template.ParseFiles("edit.html")
		t.Execute(w, ToDo{title, strings.Fields(list)})
	}

}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("URL: ", r.URL.Path[1:])

	title := r.URL.Query().Get("title")
	log.Println("deleting item", title)
	delete(todos, title)

	http.Redirect(w, r, "/view/", http.StatusFound)
}

func main() {
	http.HandleFunc("/import/", importHandler)
	http.HandleFunc("/export/", exportHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/save/", saveHandler)

	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", fs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
