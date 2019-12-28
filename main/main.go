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

func TravelFolder(src string, out string) {
	rd, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return
	}

	luaFile := new(writer.LuaFile)
	err = CompilePackage(luaFile, src, out)
	if err != nil {
		fmt.Println("compile fail:", src)
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := src + "/" + fi.Name()

			luaFile.AppendFile(fi.Name())
			err = CompilePackage(luaFile, fullDir, out)
			if err != nil {
				fmt.Println("compile fail:", src)
				return
			}
			outDir := out + "/" + fi.Name()
			TravelFolder(fullDir, outDir)
		}
	}
}

func CompilePackage(writer writer.LuaWriter,folder string, out string) error {
	writer.Reset()

	rd, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("read dir failed:", err)
		return err
	}
	hasGoFile := false
	for _, fi := range rd {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
			hasGoFile = true

			path := folder + "/" + fi.Name()
			fSet := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fSet, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			//ast.Print(fSet, f)
			writer.AppendFile(strings.TrimSuffix(fi.Name(), ".go"))
			(*translate.AstFile)(f).Translate(writer)
			writer.AddSourceInfo(fSet)
		}
	}

	outPath := out + "/" + writer.GetPackageName() + ".lua"
	os.Remove(outPath)
	fileObj,err := os.OpenFile(outPath,os.O_RDWR|os.O_CREATE,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()


	writer.Write(fileObj)
	if !hasGoFile {
		os.Remove(outPath)
	}

	return nil
}
func compile(src, out string) {
	TravelFolder(src, out)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("param error, try with:compile <src> <out>")
		return
	}
	compile(os.Args[1], os.Args[2])
}