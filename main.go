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
)

type Tree struct {
  Name string
  Description string
}

type Page struct {
    Title string
    Body  []byte
}

var t *template.Template
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := "store/" + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

// var templates = template.Must(template.ParseFiles("templates/footer.html", "templates/view.html"))
// func init() {
//   templates = template.Must(template.ParseFiles("templates/view.html", "templates/footer.html"))
// }

func main() {
  http.HandleFunc("/", handleRoot)
  http.HandleFunc("/_ah/health", healthCheckHandler)

  http.HandleFunc("/cities", citiesHandler)
  http.HandleFunc("/colors", func(w http.ResponseWriter, r *http.Request) {
    t, _ = template.ParseFiles("templates/colors.html", "templates/layout.html")
    err := t.ExecuteTemplate(w, "layout", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  })
  http.HandleFunc("/trees", treesHandler)
  http.HandleFunc("/test", testHandler)

  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  
  log.Print("Listening on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
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
    // err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    // templates, _ = template.ParseFiles("view.html")
    // if tmpl == "edit" {
    //   t, _ = template.ParseFiles("templates/edit.html", "templates/layout.html")
    // } else {
    //   t, _ = template.ParseFiles("templates/view.html", "templates/layout.html")
    // }
    t, _ = template.ParseFiles("templates/" + tmpl + ".html", "templates/layout.html")
    err := t.ExecuteTemplate(w, "layout", p)
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
      http.NotFound(w, r)
      return
    }
    t, _ = template.ParseFiles("templates/home.html", "templates/layout.html")
    err := t.ExecuteTemplate(w, "layout", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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
    var initTime = time.Now()
    t, _ = template.ParseFiles("templates/cities.html", "templates/layout.html")
    err := t.ExecuteTemplate(w, "layout", initTime)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func treesHandler(w http.ResponseWriter, r *http.Request) {
    trees := []Tree{
      {"Oak", "Quercus robur"},
      {"Maple", "Acer pseudoplatanus"},
    }


    t, _ = template.ParseFiles("templates/trees.html", "templates/layout.html")
    err := t.ExecuteTemplate(w, "layout", struct{ Trees []Tree }{trees})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
