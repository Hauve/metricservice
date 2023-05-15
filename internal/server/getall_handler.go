package server

import (
	"github.com/Hauve/metricservice.git/internal/logger"
	"html/template"
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

	for _, m := range s.storage.GetMetrics() {
		data[m.MType][m.ID] = m.GetValue()
	}

	tmpl, err := template.New("metrics").Parse(htmlTmpl)
	if err != nil {
		logger.Log.Errorf("cannot parse html template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		logger.Log.Errorf("cannot write content to the client: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
