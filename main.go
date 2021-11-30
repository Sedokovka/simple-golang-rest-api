package main

import (
  "log"
  "io/ioutil"
  "net/http"
  "html/template"
)

type Page struct {
  Title string
  Body []byte
}

func (p *Page) save() error {
  fn := p.Title +".txt"
  return ioutil.WriteFile(fn, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, page *Page, templateName string) {
  t, _ := template.ParseFiles(templateName + ".html")
  t.Execute(w, page)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  page, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/", http.StatusFound)
    page = &Page{Title: title}
  }

  renderTemplate(w, page, "view")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  page, err := loadPage(title)
  if err != nil {
    page = &Page{Title: title}
  }

  renderTemplate(w, page, "edit")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p:= &Page{Title: title, Body: []byte(body)}
  p.save()
  http.Redirect(w, r, "/save/"+title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
