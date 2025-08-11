package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func isExecutable(mode os.FileMode) bool {
	return mode.IsRegular() && (mode&0o111 != 0)
}

func listExecutables(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if isExecutable(info.Mode()) {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: which <param1> <param2> ...")
		os.Exit(1)
	}

	params := os.Args[1:]
	for i, p := range params {
		fmt.Printf("param %d: %s\n", i+1, p)
	}

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		execs, err := listExecutables(dir)
		if err != nil {
			continue
		}

		fmt.Printf("Executables in %s:\n", dir)

		for _, exec := range execs {

			fmt.Println(exec)
			base := filepath.Base(exec)
			fmt.Println(base)
		}
	}
}
