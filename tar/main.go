package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"os"
)

type Tar struct {
	file *os.File
}

func NewTar(filename string) (*Tar, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &Tar{file: file}, nil
}

func (t *Tar) Close() error {
	return t.file.Close()
}

func (t *Tar) ListFiles() error {
	tarReader := tar.NewReader(t.file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println(header.Name)
	}

	return nil
}

func (t *Tar) ExtractFiles() error {
	_, err := t.file.Seek(0, 0)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(t.file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg { // Regular file
			fmt.Printf("File: %s\n", header.Name)
			content, err := io.ReadAll(tarReader)
			if err != nil {
				return err
			}
			fmt.Printf("Content:\n%s\n", string(content))
			fmt.Println("------------------------")
		}
	}

	return nil
}

func main() {
	t := flag.Bool("t", false, "list files for stdin")
	tf := flag.Bool("tf", false, "list files")
	xf := flag.Bool("xf", false, "extract files")

	flag.Parse()

	filename := ""
	if *tf || *xf {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("Usage: tar <filename>")
			os.Exit(1)
		}
		filename = args[0]
	} else if *t {
		stdContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading input: ", err)
			os.Exit(1)
		}

		tmpfile, err := os.CreateTemp("", "tar-")
		if err != nil {
			fmt.Println("Error creating temporary file:", err)
			os.Exit(1)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(stdContent); err != nil {
			fmt.Println("Error writing to temporary file:", err)
			os.Exit(1)
		}
		filename = tmpfile.Name()
	}

	tar, err := NewTar(filename)
	if err != nil {
		fmt.Println("Error opening tar file: ", err)
		os.Exit(1)
	}

	defer tar.Close()

	if *xf {
		err := tar.ExtractFiles()
		if err != nil {
			fmt.Println("Error extracting files:", err)
			os.Exit(1)
		}
	} else if *tf || *t {
		err := tar.ListFiles()
		if err != nil {
			fmt.Println("Error listing files:", err)
			os.Exit(1)
		}
	}
}
