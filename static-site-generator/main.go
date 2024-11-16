package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/russross/blackfriday/v2"
)

type SiteGenerator struct {
	port         string
	publicFolder string
}

func NewSiteGenerator(port string) *SiteGenerator {
	publicFolder := "public"
	err := os.Mkdir(publicFolder, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	return &SiteGenerator{
		port:         port,
		publicFolder: publicFolder,
	}
}

func (s *SiteGenerator) Build() error {
	files, err := filepath.Glob("content/*.md")
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, file := range files {
		fmt.Println(file)
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			mdContent, err := os.ReadFile(file)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				cancel()
				return
			}

			htmlContent := blackfriday.Run(mdContent)
			htmlFile := filepath.Base(file[:len(file)-3] + ".html")
			f, err := os.Create(filepath.Join(s.publicFolder, htmlFile))
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				cancel()
				return
			}
			defer f.Close()

			_, err = f.Write(htmlContent)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				cancel()
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

func (s *SiteGenerator) Server() error {
	http.Handle("/", http.FileServer(http.Dir(s.publicFolder)))
	fmt.Printf("Server starting on http://localhost:%s\n", s.port)

	return http.ListenAndServe(":"+s.port, nil)
}

func main() {

	build := flag.Bool("b", false, "Build the site")
	server := flag.Bool("s", false, "Start the server")
	port := flag.String("p", "8080", "Port to listen on")
	flag.Parse()

	siteGenerator := NewSiteGenerator(*port)

	if *build {
		err := siteGenerator.Build()
		if err != nil {
			log.Fatal("Build error:", err)
		}
		return
	}

	if *server {
		if err := siteGenerator.Server(); err != nil {
			log.Fatal("Server error:", err)
		}
	}

	flag.PrintDefaults()
}
