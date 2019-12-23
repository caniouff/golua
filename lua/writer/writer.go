package writer

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"sort"
)

type LuaLine struct {
	buffer bytes.Buffer
	astPos token.Pos
}

type LuaFile struct {
	packageName string
	predefine []string
	lines []LuaLine
	lineCount int
	localScopeStack int
}

type LuaWriter interface {
	AppendLine(line int, astPos token.Pos, content string)
	Write(fSet *token.FileSet, fileName string)
	AppendDef(name string)
	SetPackageName(name string)
	GetPackageName()string
	IsGlobalScope() bool
	EnterLocalScope()
	LeaveLocalScope()
}

type LuaReader interface {
	LineCount() int
	Line(line int) string
}
func(f *LuaFile)AppendLine(line int, astPos token.Pos, content string) {
	if line < 0 {
		f.lines = append(f.lines, LuaLine{})
		f.lineCount++
	} else if f.lineCount < line {
		for i:= f.lineCount; i < line; i++{
			f.lines = append(f.lines, LuaLine{})
			f.lineCount++
		}
	}
	if len(f.lines) <= 0 {
		f.lines = append(f.lines, LuaLine{})
		f.lineCount ++
	}
	curLine := &f.lines[len(f.lines) - 1]
	curLine.buffer.WriteString(content)
	if astPos != token.NoPos {
		curLine.astPos = astPos
	}
}

func assertWrite(fileObj *os.File, content string, hasFailed bool, failedOut string)  {
	if hasFailed {
		fmt.Println(failedOut)
		os.Exit(2)
		return
	}
	_,err := io.WriteString(fileObj, content)
	if err != nil {
		fmt.Println(failedOut)
		os.Exit(2)
		return
	}
}

func checkWrite(fileObj *os.File, content string) bool  {
	_,err := io.WriteString(fileObj, content)
	return err != nil
}

func GenBrunchTree(names []string) [] string{
	count := len(names)
	stack := make([]int, 0)
	stack = append(stack, 1)
	stack = append(stack, count)

	tree := make([]string, 0)
	for {
		stackLen := len(stack)
		if stackLen < 2 {
			break
		}

		left := stack[stackLen-1]
		right := stack[stackLen-2]

		stack = stack[0:stackLen-2]

		mid := (left + right) / 2

		if mid == -1 {
			tree = append(tree, "else")
		}else if mid == -2 {
			tree = append(tree, "end")
		}else if left == right {
			tree = append(tree, fmt.Sprintf("return %s ", names[mid - 1]))
		} else {
			tree = append(tree, fmt.Sprintf("if index > %d then", mid))

			stack = append(stack, -2)
			stack = append(stack, -2)

			stack = append(stack, right)
			stack = append(stack, mid)

			stack = append(stack, -1)
			stack = append(stack, -1)

			if mid + 1 <= left {
				stack = append(stack, mid + 1)
				stack = append(stack, left)
			}
		}
	}
	return tree
}

func(f *LuaFile)Write(fSet *token.FileSet, fileName string){
	os.Remove(fileName)
	fileObj,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()

	bSucc := true
	//写预定义内容
	sort.Strings(f.predefine)
	for _, name := range f.predefine{
		if bSucc = checkWrite(fileObj, fmt.Sprintf("local %s = nil \n", name)); bSucc {
			break
		}
	}

	assertWrite(fileObj, "--------------------------", bSucc,
		fmt.Sprintf("failed to write predefine:%s", fileName))
	//写翻译的代码内容
	for _, line := range f.lines {

		if bSucc = checkWrite(fileObj, line.buffer.String()); bSucc {
			break
		}
		if line.astPos > 0 {
			pos := fSet.Position(line.astPos)
			if bSucc = checkWrite(fileObj, fmt.Sprintf("--[[%s:%d]]", pos.Filename, pos.Line));bSucc{
				break
			}
		}

		if bSucc = checkWrite(fileObj, "\n"); bSucc {
			break
		}
	}

	assertWrite(fileObj, "--------------------------", bSucc,
		fmt.Sprintf("Failed to write the code lines:%s",fileName))

	//写入包内容
	/*
	local predefine = { A= 1, B=2, ...}
	local predefineCount = 10
	return setmetatable({}, {
	__newindex = function() error("package is readonly") end,
	__index = function(t, key)
		local index = predefine(key)
		if not index then error(string.format("not found %s in package:packageName", ))return nil end
		if index > 5 {
	        if index > 7 {
	             if index > 9 {

	             }elseif index == 9 {
	             }else{
	             }
	        } elseif index == 7 {

	        } else {
			}
		} elseif index == 5 {
			return F
	    } else {

		}
	end
	})
	*/
	//1.创建查找表
	assertWrite(fileObj, "local predefine = {\n", bSucc,
		fmt.Sprintf("Failed to write the predefine:%s",fileName))
	for index, name := range f.predefine {
		if bSucc = checkWrite(fileObj, fmt.Sprintf("%s = %d\n", name, index)); bSucc {
			break
		}
	}
	assertWrite(fileObj, "}\n", bSucc,
		fmt.Sprintf("Failed to write the predefine:%s",fileName))
	//2.开始package 的meta table
	assertWrite(fileObj, "local index = predefine[key]\n", bSucc,
		fmt.Sprintf("Failed to write the index assignment:%s",fileName))
	//3.实现二分查找算法

	//4.结束package
}

func (f *LuaFile) LineCount() int  {
	return len(f.lines)
}

func (f *LuaFile) Line(line int) string  {
	return f.lines[line].buffer.String()
}

func (f *LuaFile) EnterLocalScope()  {
	f.localScopeStack ++
}

func (f *LuaFile) LeaveLocalScope() {
	f.localScopeStack --
	if f.localScopeStack < 0 {
		fmt.Errorf("localScopeStack error")
	}
}

func (f *LuaFile) IsGlobalScope() bool  {
	return f.localScopeStack == 0
}

func (f *LuaFile) SetPackageName(name string)  {
	f.packageName = name
}

func (f *LuaFile) GetPackageName()string  {
	return f.packageName
}

func (f *LuaFile) AppendDef(name string)  {
	f.predefine = append(f.predefine, name)
}
