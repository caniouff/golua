package writer

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"sort"
	"strings"
)

type LuaLine struct {
	buffer bytes.Buffer
	astPos token.Pos
	setSource bool
}

type LuaFile struct {
	packageName     string
	fileNames       []string
	fileSets        []*token.FileSet
	predefine       []string
	lines           []LuaLine
	lineCount       int
	lastDefinedLine int
	localScopeStack int
	needSetIota     bool
	iotaExpr        string
}

type LuaWriter interface {
	AppendLine(line int, astPos token.Pos, content string)
	AppendFile(name string, fSet *token.FileSet)
	Write(fileObj *os.File)
	AddSourceInfo(fSet *token.FileSet)
	AppendDef(name string)
	SetPackageName(name string)
	GetPackageName()string
	GetInitFuncName()string
	IsGlobalScope() bool
	EnterLocalScope()
	LeaveLocalScope()
	Reset()
	MakeEmptyWriter() LuaWriter

	SetIota(set bool)
	HasSetIota() bool
	SetIotaExpr(expr string)
	GetIotaExpr() string
	ClearIota()
	String() string
}

type LuaReader interface {
	LineCount() int
	Line(line int) string
}
func(f *LuaFile)AppendLine(line int, astPos token.Pos, content string) {
	definedLine := f.DefinedLine(astPos)
	if line == 0 {
		if f.lastDefinedLine > 0 && definedLine > f.lastDefinedLine {
			line = f.lineCount + definedLine - f.lastDefinedLine
		}
		f.lastDefinedLine = definedLine
	} else {
		f.lastDefinedLine = 0
	}
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

func assertWrite(fileObj *os.File, content string, bSucc bool, failedOut string)  {
	if !bSucc {
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
	return err == nil
}

func(f *LuaFile)Write(fileObj *os.File){
	bSucc := true
	fileName := fileObj.Name()
	//写预定义内容
	sort.Strings(f.predefine)
	for _, name := range f.predefine{
		if bSucc = checkWrite(fileObj, fmt.Sprintf("local %s = nil \n", name)); !bSucc {
			break
		}
	}

	assertWrite(fileObj, "--------------------------\n", bSucc,
		fmt.Sprintf("failed to write predefine:%s", fileName))
	//写翻译的代码内容
	for _, line := range f.lines {
		if bSucc = checkWrite(fileObj, line.buffer.String() + "\n"); !bSucc {
			break
		}
	}

	//调用初始化函数
	for _, name := range f.fileNames {
		if bSucc = checkWrite(fileObj, fmt.Sprintf("if %s_init then %s_init() end\n", name, name)); !bSucc {
			break
		}
	}
	assertWrite(fileObj, "--------------------------\n", bSucc,
		fmt.Sprintf("Failed to write the code lines:%s",fileName))

	//写入包内容
	//1.创建查找表
	//1.1 预定义表
	assertWrite(fileObj, "local predefine = {\n", bSucc,
		fmt.Sprintf("Failed to write the predefine:%s",fileName))
	for _, name := range f.predefine {
		if bSucc = checkWrite(fileObj, fmt.Sprintf("%s = %s,\n", name, name)); !bSucc {
			break
		}
	}
	assertWrite(fileObj, "}\n", bSucc,
		fmt.Sprintf("Failed to write the predefine:%s",fileName))

	//1.2定义表头
	assertWrite(fileObj, `return setmetatable(predefine, {
__newindex = function() error("package is readonly") end,
`, bSucc,
    fmt.Sprintf("Failed to write the metatable header:%s", fileName))

	//结束package
	assertWrite(fileObj, "})", bSucc,
		fmt.Sprintf("Failed to write the ending:%s", fileName))
}

func (f *LuaFile) AddSourceInfo(fSet *token.FileSet) {
	//写翻译的代码内容
	for index := range f.lines {
		line := &f.lines[index]
		if !line.setSource && line.astPos > 0 {
			pos := fSet.Position(line.astPos)
			line.buffer.WriteString(fmt.Sprintf("--[[%s:%d]]", pos.Filename, pos.Line))
			line.setSource = true
		}
	}
}

func (f *LuaFile) AppendFile(name string, fSet *token.FileSet) {
	f.fileNames = append(f.fileNames, name)
	f.fileSets = append(f.fileSets, fSet)
}

func (f *LuaFile) DefinedLine(astPos token.Pos) int{
	if len(f.fileSets) < 0 {
		return 0
	}
	fSet := f.fileSets[len(f.fileSets) - 1]
	if astPos > 0 {
		pos := fSet.Position(astPos)
		return pos.Line
	} else {
		return 0
	}
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

func (f *LuaFile) GetInitFuncName() string  {
	return f.fileNames[len(f.fileNames) - 1] + "_init"
}

func (f *LuaFile) AppendDef(name string)  {
	f.predefine = append(f.predefine, name)
}

func (f *LuaFile) Reset()  {
	f.predefine = []string{}
	f.fileNames = []string{}
	f.packageName = ""
	f.localScopeStack = 0
	f.lineCount = 0
	f.lastDefinedLine = 0
	f.fileSets = []*token.FileSet{}
	f.lines = []LuaLine{}
}

func (f *LuaFile) SetIota(set bool) {
	f.needSetIota = set
}

func (f *LuaFile) HasSetIota() bool {
	return f.needSetIota
}

func (f *LuaFile) SetIotaExpr(expr string)  {
	f.iotaExpr = expr
}

func (f *LuaFile) GetIotaExpr() string  {
	return f.iotaExpr
}

func (f *LuaFile) ClearIota() {
	f.iotaExpr = ""
	f.needSetIota = false
}

func (f *LuaFile) MakeEmptyWriter() LuaWriter {
	emptyFile := &LuaFile{}
	emptyFile.fileSets = f.fileSets
	emptyFile.fileNames = f.fileNames
	emptyFile.packageName = f.packageName

	return emptyFile
}

func (f *LuaFile)String() string  {
	sb := strings.Builder{}
	for index, line := range f.lines {
		sb.WriteString(line.buffer.String())
		if index > 0 && index < len(f.lines) - 1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}