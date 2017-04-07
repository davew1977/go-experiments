package wiki

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
)

type PageService interface {
	Save(page *Page) error
	Load(name string) (*Page, error)
	ListPages() []string
}

type filePageService struct {
     pageNames []string //cache of existing page names
}


var _ PageService = (*filePageService)(nil) //asserts that filePageService actually implements PageService
var _ PageService = (*dbPageService)(nil) //asserts that filePageService actually implements PageService

func NewFilePageService() PageService {
	service := &filePageService{}
	service.init()
	return service
}
func NewDbPageService() PageService {
	service := &dbPageService{}
	service.init()
	return service
}

func (f *filePageService) init() {
	pattern := fmt.Sprintf("%s/*.txt", viper.Get("paths.data"))
	log.Print(pattern)
	files, _ := filepath.Glob(pattern)
	f.pageNames = make([]string, len(files))
	for i,fName := range(files) {
		f.pageNames[i] = strings.Split(strings.Split(fName, ".")[0], "/")[1]
		//log.Print(f)
	}
}

func (f *filePageService) Save(p *Page) error {
	filename := fmt.Sprintf("%s/%s.txt", viper.Get("paths.data"), p.Title)
	err := ioutil.WriteFile(filename, p.Body, 0600)
	f.init()
	return err
}
func (f *filePageService) Load(title string) (*Page, error) {
	filename := fmt.Sprintf("%s/%s.txt", viper.Get("paths.data"), title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}
func (f *filePageService)  ListPages() []string {
	return f.PageNames()
}
func (f *filePageService)  PageNames() []string {
	return f.pageNames
}

//DB page service:

type dbPageService struct {
	db *sql.DB
}

func (f *dbPageService) init() {
	log.Print("initializing db page service")
	db, _ := sql.Open("mysql", "dev:dev@tcp(localhost:8889)/go_wiki")
	f.db = db
	if err := f.db.Ping(); err != nil {
		log.Fatal(err)
	}
}
func (f *dbPageService) Save(p *Page) error {
	res, err :=f.db.Exec("Update Page set Body=? where Name=?", p.Body, p.Title)
	if err != nil {
		log.Fatal(err)
	}
	rowCount, _ := res.RowsAffected()
	if rowCount == 0 {
		_, err :=f.db.Exec("insert into Page values (?,?)", p.Title, p.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
func (f *dbPageService) Load(title string) (*Page, error) {
	rows, _ := f.db.Query("select * from Page where Name=?", title)
	defer rows.Close()
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.Title, &page.Body); err != nil {
			log.Fatal(err)
		}
		return &page, nil
	}
	return nil, errors.New("page not found: " + title)
}
func (f *dbPageService)  ListPages() []string {
	return f.PageNames()
}
func (f *dbPageService)  PageNames() []string {
	rows, _ := f.db.Query("select Name from Page")
	defer rows.Close()
	var result []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Fatal(err)
		}
		result = append(result, title)

	}
	return result
}

