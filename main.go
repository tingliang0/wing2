package main

import (
	"flag"
	"go/build"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func defaultAssetPath() string {
	p, err := build.Default.Import(".", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func benchmarkHandler(c http.ResponseWriter, req *http.Request) {
	h.broadcast <- []byte("test message")
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	go h.run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)
	Info.Println("listen on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		Error.Fatal("ListenAndServe:", err)
	}
}
