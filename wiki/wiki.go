package main

import (
	"net/http"
	"html/template"
	"regexp"
	"log"
	"github.com/davew1977/go-experiments/wiki/pkg/wiki"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var pageService wiki.PageService

func init() {
	loadConfig()
	//pageService = wiki.NewFilePageService()
	pageService = wiki.NewDbPageService()
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pageService.Load(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", map[string]interface{}{
		"Title" : p.Title,
		"Body" : pageFormat(p, pageService.(wiki.RenderContext)), //a bit like a cast - we know the runtime type of page service will satisfy RenderContext
	})
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pageService.Load(title)
	if err != nil {
		p = &wiki.Page{Title: title}
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
	p := &wiki.Page{Title: title, Body: []byte(body)}
	pageService.Save(p)
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {

	pages := pageService.ListPages()
	log.Print(pages)
	renderTemplate(w, "list", pages)
}

func main() {
	log.Print("started wiki")
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/list/", listHandler)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
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

func pageFormat(p *wiki.Page, c wiki.RenderContext) template.HTML {
	html := strings.Replace(string(p.Body), "\n", "<br>", -1)
	html = string(regexp.MustCompile(strings.Join(c.PageNames(), "|")).ReplaceAllFunc([]byte(html), func(r []byte) []byte {
		if (p.Title == string(r)) {
			return r
		}
		return []byte(fmt.Sprintf("<a href='/view/%s'>%s</a>", r, r))

	}))
	return template.HTML(html);
}