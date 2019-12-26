package main

import (
	"flag"
	"fmt"
	"github.com/south-rit/copy-generator/bootstrap"
	"github.com/south-rit/copy-generator/parser"
	"os"
	"path/filepath"
	"strings"
)

var allStructs = flag.Bool("all", false, "generate marshaler/unmarshalers for all structs in a file")
var specifiedName = flag.String("output_filename", "", "specify the filename of the output")

func generate(fname string) (err error) {
	fInfo, err := os.Stat(fname)
	if err != nil {
		return err
	}

	p := parser.Parser{AllStructs: *allStructs}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return fmt.Errorf("error parsing %v: %v", fname, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fname, p.PkgName+"_copy.go")
	} else {
		if s := strings.TrimSuffix(fname, ".go"); s == fname {
			return fmt.Errorf("filename must end in '.go'")
		} else {
			outName = s + "_copy.go"
		}
	}

	if *specifiedName != "" {
		outName = *specifiedName
	}

	g := bootstrap.Generator{
		PkgPath: p.PkgPath,
		PkgName: p.PkgName,
		Types:   p.StructNames,
		OutName: outName,
	}

	_ = g

	if err := g.Run(); err != nil {
		return fmt.Errorf("bootstrap failed: %v", err)
	}
	return nil
}

func main() {
	flag.Parse()

	files := flag.Args()

	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, fname := range files {
		if err := generate(fname); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
