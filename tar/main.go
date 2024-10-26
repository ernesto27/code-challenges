package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"syscall"
)

type Tar struct {
	file       *os.File
	newTarName string
}

func NewTar(tarName string, create bool) (*Tar, error) {
	if create {
		return &Tar{newTarName: tarName}, nil
	}

	file, err := os.Open(tarName)
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

		content, err := io.ReadAll(tarReader)
		if err != nil {
			return err
		}
		file, err := os.Create(header.Name)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tar) CreateTar(args []string) error {
	fmt.Println(args)
	newTarFile, err := os.Create(t.newTarName)
	if err != nil {
		return err
	}
	defer newTarFile.Close()

	tarWriter := tar.NewWriter(newTarFile)

	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		sysStat, ok := stat.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get system-specific file info for %s", arg)
		}

		u, err := user.LookupId(fmt.Sprint(sysStat.Uid))
		if err != nil {
			return err
		}

		group, err := user.LookupGroupId(fmt.Sprint(sysStat.Gid))
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name:    file.Name(),
			Size:    stat.Size(),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
			Uname:   u.Username,
			Gname:   group.Name,
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}

	return nil

}

func main() {
	t := flag.Bool("t", false, "list files for stdin")
	tf := flag.Bool("tf", false, "list files")
	xf := flag.Bool("xf", false, "extract files")
	cf := flag.Bool("cf", false, "create tar file")

	flag.Parse()

	tarName := ""
	args := flag.Args()

	if *tf || *xf || *cf {
		if len(args) == 0 {
			fmt.Println("Usage: tar <tarName>")
			os.Exit(1)
		}
		tarName = args[0]
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
		tarName = tmpfile.Name()
	}

	tar, err := NewTar(tarName, *cf)
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
	} else if *cf {
		err := tar.CreateTar(args[1:])
		if err != nil {
			fmt.Println("Error creating tar file:", err)
			os.Exit(1)
		}
	}
}
