package handler

import (
	"html/template"
	"log"
	"net/http"
)

func (h *server) index(w http.ResponseWriter, r *http.Request) {
	//Подключения шаблона
	t, err := template.ParseFiles("templates/index.html", 
									"templates/header.html", 
									"templates/footer.html")
	// Проверка на ошибки
	if err != nil {
		log.Fatalln(w, err.Error())
	}

	Obj := make(map[string][]string)
	Obj["Airtable"], _ = h.services.Airtable.Get()
	Obj["NextCloud"], _ = h.services.NextCloud.Get()

	t.ExecuteTemplate(w, "index", Obj)
}