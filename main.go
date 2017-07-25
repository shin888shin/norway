// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample helloworld is a basic App Engine flexible app.
package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
  "html/template"
)

type Tree struct {
  Name string
  Description string
}

func main() {
  http.HandleFunc("/", handle)
  http.HandleFunc("/_ah/health", healthCheckHandler)

  http.HandleFunc("/cities", citiesHandler)
  
  http.HandleFunc("/colors", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "path: %s\n", r.URL.Path)
  })

  http.HandleFunc("/trees", treesHandler)
  log.Print("Listening on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }
  fmt.Fprint(w, "Norway")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "ok")
}

func citiesHandler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprint(w, "Rome, New York, Mexico City, Tokyo")
    tmpl.Execute(w, time.Since(initTime))
}

func treesHandler(w http.ResponseWriter, r *http.Request) {
  trees := []Tree{
    {"Oak", "Quercus robur"},
    {"Maple", "Acer pseudoplatanus"},
  }
  tmpl := template.Must(template.ParseFiles("templates/trees.html"))
  tmpl.Execute(w, struct{ Trees []Tree }{trees})
}


var initTime = time.Now()
var tmpl = template.Must(template.New("front").Parse(`
<html><body>
<p>
Oakland, Moscow, Rome, Florence, New York, Mexico City, Tokyo! 세상아 안녕!
</p>
<p>
This instance has been running for <em>{{.}}</em>.
</p>
</body></html>
`))