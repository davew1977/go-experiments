package wiki

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type PageService interface {
	Save(page *Page) error
	Load(name string) (*Page, error)
	ListPages() []string
}

type filePageService struct {
     pageNames []string
}

func NewFilePageService() PageService {
	service := &filePageService{}
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

