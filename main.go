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
  "io/ioutil"
  "regexp"
  "errors"
)

type Tree struct {
  Name string
  Description string
}

type Page struct {
    Title string
    Body  []byte
}

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := "store/" + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func main() {
  http.HandleFunc("/", handle)
  http.HandleFunc("/_ah/health", healthCheckHandler)

  http.HandleFunc("/cities", citiesHandler)
  http.HandleFunc("/colors", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "path: %s\n", r.URL.Path)
  })
  http.HandleFunc("/trees", treesHandler)
  http.HandleFunc("/test", testHandler)

  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  
  log.Print("Listening on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    log.Print("---> url: ", r.URL)
    m := validPath.FindStringSubmatch(r.URL.Path)
    log.Print("---> valid url: ", m)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    // t, err := template.ParseFiles("templates/" + tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func loadPage(title string) (*Page, error) {
    filename := "store/" + title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
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

func testHandler(w http.ResponseWriter, r *http.Request) {
    p1 := &Page{Title: "test", Body: []byte("test test test")}
    p1.save()
    p2, _ := loadPage("test")
    fmt.Println(string(p2.Body))
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