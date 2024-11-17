package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	name string
	size int64
	hash string
}

type FileFinder struct {
	dir     string
	files   []File
	minSize int64
}

func NewFileFinder(dir string, minSize int64) *FileFinder {
	return &FileFinder{
		dir:     dir,
		minSize: minSize,
	}
}

func (f *FileFinder) Find() error {
	return filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if info.Size() < f.minSize {
				return nil
			}

			isSymlink := info.Mode()&os.ModeSymlink != 0

			hashString, err := f.generateMD5(path)
			if err != nil {
				return err
			}

			file := File{
				name: path,
				size: info.Size(),
				hash: hashString,
			}

			var symLinkInfo string
			if isSymlink {
				originalPath, err := f.getOriginalPath(path)
				if err != nil {
					return err
				}
				symLinkInfo = " (symlink file) => " + originalPath
			}

			fmt.Printf("%s %s\n", path, symLinkInfo)

			fileDuplicated, found := f.FindSameHash(file)
			if found || isSymlink {
				reader := bufio.NewReader(os.Stdin)
				for {
					if !isSymlink {
						fmt.Println("Which file should be deleted?")
						fmt.Printf("  1) %s\n", file.name)
						fmt.Printf("  2) %s\n", fileDuplicated.name)
						fmt.Print("Enter 1 or 2, or other key to keep: ")
					} else {
						fmt.Println("Symlink file found, do you want to remove it?")
						fmt.Print("Enter 1 to remove or any other key to keep: ")
					}

					input, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					choice := strings.TrimSpace(input)

					switch choice {
					case "1":
						return os.Remove(file.name)
					case "2":
						if isSymlink {
							return nil
						}
						return os.Remove(fileDuplicated.name)
					default:
						return nil
					}
				}
				// fmt.Println("Which file should be deleted?")
				// fmt.Printf("  1) %s\n", file.name)
				// fmt.Printf("  2) %s\n", fileDuplicated.name)
				// fmt.Printf("Probable duplicates: %s - %s\n", file.name, fileDuplicated.name)
				// same, err := compareByteByByte(file.name, fileDuplicated.name)
				// if err != nil {
				// 	return err
				// }
				// if !same {
				// 	fmt.Println("Files are different byte by byte")
				// }
			}
			f.files = append(f.files, file)

		}

		return nil
	})
}

func (f *FileFinder) FindSameHash(file File) (File, bool) {
	for _, f := range f.files {
		if f.hash == file.hash {
			return f, true
		}
	}

	return File{}, false
}

func (f *FileFinder) FindSameSize(file File) (File, bool) {
	for _, f := range f.files {
		if f.size == file.size {
			return f, true
		}
	}

	return File{}, false
}

func (f *FileFinder) generateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func (f *FileFinder) compareByteByByte(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	const chunkSize = 64 * 1024 // 64KB chunks
	b1 := make([]byte, chunkSize)
	b2 := make([]byte, chunkSize)

	for {
		n1, err1 := f1.Read(b1)
		n2, err2 := f2.Read(b2)

		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		}

		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		if n1 != n2 {
			return false, nil
		}

		if !bytes.Equal(b1[:n1], b2[:n2]) {
			return false, nil
		}
	}
}

func (f *FileFinder) getOriginalPath(symlinkPath string) (string, error) {
	originalPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink: %w", err)
	}
	if !filepath.IsAbs(originalPath) {
		dir := filepath.Dir(symlinkPath)
		originalPath = filepath.Join(dir, originalPath)
	}

	_, err = os.Stat(originalPath)
	if err != nil {
		return "", fmt.Errorf("original file not found: %w", err)
	}

	return originalPath, nil
}

func main() {
	dir := flag.String("d", "content", "Directory to search for files")
	minSize := flag.Int64("s", 0, "Minimum file size to consider")
	flag.Parse()

	f := NewFileFinder(*dir, *minSize)
	err := f.Find()
	if err != nil {
		fmt.Println(err)
	}
}
