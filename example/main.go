package main

import (
	"flag"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hysios/linetmpl"
	"github.com/hysios/log"
	"github.com/hysios/utils/response"
)

const (
	left  = "["
	right = "]"
)

var (
	funcs = template.FuncMap{}
)

type Map = map[string]interface{}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		var (
			vars    = mux.Vars(r)
			name    = vars["name"]
			content []byte
		)

		f, err := os.OpenFile(filepath.Join("templates", name+".tpl"), os.O_RDONLY, 0)
		if err != nil {
			response.AbortErr(w, http.StatusNotFound, err)
			return
		}

		if content, err = io.ReadAll(f); err != nil {
			response.AbortErr(w, http.StatusInternalServerError, err)
			return
		}

		tree, err := linetmpl.Parse(name, string(content))
		if err != nil {
			response.AbortErr(w, http.StatusInternalServerError, err)
			return
		}

		response.Jsonify(w, &Map{"data": tree})

	}).Methods("GET")
	log.Infof("inline template server at %s", ":8070")
	http.ListenAndServe(":8070", r)
}
