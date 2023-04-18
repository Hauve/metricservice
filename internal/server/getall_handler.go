package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var htmlTmpl = `
<html>
<p>{{range .gauge}} {{.}} {{end}}</p>
<p>{{range .counter}} {{.}} {{end}}</p>
</html>
`

func (s *MyServer) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Date", time.Now().String())

	data := map[string]map[string]string{
		"gauge":   make(map[string]string),
		"counter": make(map[string]string),
	}

	for _, key := range s.storage.GetGaugeKeys() {
		value, _ := s.storage.GetGauge(key)
		data["gauge"][key] = fmt.Sprintf("%f", value)
	}

	for _, key := range s.storage.GetCounterKeys() {
		value, _ := s.storage.GetCounter(key)
		data["counter"][key] = fmt.Sprintf("%d", value)
	}

	tmpl, err := template.New("test").Parse(htmlTmpl)
	if err != nil {
		log.Printf("cannot parse html template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		log.Printf("cannot write content to the client: %s", err)
	}
}
