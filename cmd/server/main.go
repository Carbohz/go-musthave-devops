package main

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/handler"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	host = "127.0.0.1"
	port = "8080"
	htmlFile = "D:\\Go\\yandex-praktikum\\Sprint1\\net_http\\increment1\\go-musthave-devops2\\cmd\\server\\index.html"
	//htmlFile = "index.html"
)

//var htmlTemplate *template.Template

func main() {
	//handler.PrepareHtmlFile()

	bytes, err := os.ReadFile(htmlFile)
	if err != nil {
		log.Fatal(err)
	}
	handler.HtmlTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		log.Fatal(err)
	}
	//HtmlTemplate = htmlTemplate
	RunServer()
}

func RunServer() {
	r := chi.NewRouter()
	handler.SetupRouters(r)
	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}


//func AllMetricsHandler(w http.ResponseWriter, r *http.Request) {
//	//htmlFile := "index.html"
//	htmlFile := "D:\\Go\\yandex-praktikum\\Sprint1\\net_http\\increment1\\go-musthave-devops2\\cmd\\server\\index.html"
//	htmlPage, err := os.ReadFile(htmlFile)
//	if err != nil {
//		log.Println("File reading error:", err)
//	}
//
//	renderData := map[string]interface{}{
//		"gaugeMetrics": gaugeMetricsStorage,
//		"counterMetrics": counterMetricsStorage,
//	}
//	tmpl := template.Must(template.New("").Parse(string(htmlPage)))
//	tmpl.Execute(w, renderData)
//}
