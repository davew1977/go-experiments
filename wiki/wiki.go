package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"strings"
	"path/filepath"
	"regexp"
	"log"
	"fmt"
	"github.com/spf13/viper"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var pages []string

func init() {
	loadConfig()
	initPages()
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := fmt.Sprintf("%s/%s.txt", viper.Get("paths.data"), p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}
func (p *Page) FormatBody() template.HTML {
	html := strings.Replace(string(p.Body), "\n", "<br>", -1)
	html = string(regexp.MustCompile(strings.Join(pages, "|")).ReplaceAllFunc([]byte(html), func(r []byte) []byte {
		if(p.Title == string(r)) {
			return r
		}
		return []byte(fmt.Sprintf("<a href='/view/%s'>%s</a>", r, r))

	}))
	return template.HTML(html);
}


func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func loadPage(title string) (*Page, error) {
	filename := fmt.Sprintf("%s/%s.txt", viper.Get("paths.data"), title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templateName := fmt.Sprintf("%s/%s.go.html", viper.Get("paths.templates"), tmpl)
	log.Print(templateName)
	t, _ := template.ParseFiles(templateName)
	t.Execute(w, data)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	initPages()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {

	 log.Print(pages)
	renderTemplate(w, "list", pages)
}

func initPages() {

	pattern := fmt.Sprintf("%s/*.txt", viper.Get("paths.data"))
	log.Print(pattern)
	files, _ := filepath.Glob(pattern)
	pages = make([]string, len(files))
	for i,f := range(files) {
		pages[i] = strings.Split(strings.Split(f, ".")[0], "/")[1]
		//log.Print(f)
	}
}

func main() {
	log.Print("started wiki")
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/list/", listHandler)
	http.ListenAndServe(":8080", nil)
}
func loadConfig() {
	viper.SetEnvPrefix("wiki")
	viper.SetConfigName("wiki")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not read config file from %s.", viper.ConfigFileUsed())
	}
        log.Print(viper.Get("paths.templates"))
}