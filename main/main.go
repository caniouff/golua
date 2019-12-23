package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"lua/translate"
	"lua/writer"
	"os"
	"strings"
)

func GetAllFile(src, suffix string, files []string)([]string) {
	rd, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return files
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := src + "/" + fi.Name()
			files = GetAllFile(fullDir, suffix, files)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return files
			}
		} else {
			if strings.HasSuffix(fi.Name(), suffix) {
				fullName := src + "/" + fi.Name()
				files = append(files, fullName)
			}

		}
	}
	return files
}
func compileOne(src, out string)  {
	// Create the AST by parsing src.
	fSet := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fSet, src, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	luaFile := writer.LuaFile{}
	//ast.Print(fSet, f)
	(*translate.AstFile)(f).Translate(&luaFile)
	luaFile.Write(fSet,out)
}

func compile(src, out string) {
	var files []string
	for _, file := range GetAllFile(src, ".go", files) {
		outName := strings.Replace(file, src, out, 1)
		compileOne(file, strings.Replace(outName, ".go", ".lua", 1))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("param error, try with:compile <src> <out>")
		return
	}
	compile(os.Args[1], os.Args[2])
	names := []string{"A", "B", "C", "D", "E"}
	tree := writer.GenBrunchTree(names[0:])
	for _, leaf := range tree {
		fmt.Println(leaf)
	}
}