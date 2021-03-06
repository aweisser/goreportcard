package handlers

import (
	"log"
	"net/http"
	"text/template"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving about page")

	t := template.Must(template.New("about.html").Delims("[[", "]]").ParseFiles("templates/about.html", "templates/footer.html"))
	t.Execute(w, map[string]interface{}{
		"google_analytics_key": googleAnalyticsKey,
	})
	return
}
