package main

import (
	"flag"
	"fmt"
	"go-meta"
	_ "go-meta/cgo"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

const MetaBuildTag = "meta"

func init() {
	build.Default.BuildTags = append(build.Default.BuildTags, MetaBuildTag)
}

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		name := filepath.Base(arg)
		dir := filepath.Dir(arg)

		match, err := build.Default.MatchFile(dir, name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if match {
			if err = transform(arg); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func transform(path string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

	if err != nil {
		return err
	}

	ds := meta.Extract(file)
	if err = meta.Apply(ds); err != nil {
		return err
	}

	return printAst(generateName(path), fset, file)
}

func generateName(path string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + "_generated" + ext
}

func printAst(filepath string, fset *token.FileSet, node interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return printer.Fprint(file, fset, node)
}
