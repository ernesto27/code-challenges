package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/russross/blackfriday/v2"
)

type SiteGenerator struct {
	publicFolder string
}

func NewSiteGenerator() *SiteGenerator {
	publicFolder := "public"
	err := os.Mkdir(publicFolder, 0755)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	return &SiteGenerator{
		publicFolder: publicFolder,
	}
}

func (s *SiteGenerator) Build() error {
	files, err := filepath.Glob("content/*.md")
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(files))

	for _, file := range files {
		fmt.Println(file)
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			mdContent, err := os.ReadFile(file)
			if err != nil {
				errChan <- err
				return
			}

			htmlContent := blackfriday.Run(mdContent)
			htmlFile := filepath.Base(file[:len(file)-3] + ".html")
			f, err := os.Create(filepath.Join(s.publicFolder, htmlFile))
			if err != nil {
				errChan <- err
				return
			}
			defer f.Close()

			_, err = f.Write(htmlContent)
			if err != nil {
				errChan <- err
				return
			}
		}(file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func renderTemplate(tmpl string, data interface{}) error {
	// t, err := template.ParseFiles(tmpl)
	// if err != nil {
	// 	return err
	// }
	// f, err := ioutil.TempFile("", "*.html")
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// return t.Execute(f, data)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// mdContent, err := os.ReadFile("content.md")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// htmlContent := blackfriday.Run(mdContent)

	// data := struct {
	// 	Title   string
	// 	Heading string
	// 	Content template.HTML
	// }{
	// 	Title:   "My Page Title",
	// 	Heading: "Welcome to My Page",
	// 	Content: template.HTML(htmlContent),
	// }
	// renderTemplate(w, "template.html", data)
}

func main() {
	// 	http.HandleFunc("/", handler)
	// 	http.ListenAndServe(":8080", nil)

	sg := NewSiteGenerator()
	err := sg.Build()
	if err != nil {
		panic(err)
	}
}
