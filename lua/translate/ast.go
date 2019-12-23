package translate

import (
	"fmt"
	"go/ast"
	"go/token"
	"lua/writer"
	"strings"
)

func cast(astNode interface{}) interface{} {
	switch astNode.(type) {
	case *ast.File:
		return (*AstFile)(astNode.(*ast.File))
	case *ast.BadDecl:
		return (*AstBadDecl)(astNode.(*ast.BadDecl))
	case *ast.FuncDecl:
		return (*AstFuncDecl)(astNode.(*ast.FuncDecl))
	case *ast.GenDecl:
		return (*AstGenDecl)(astNode.(*ast.GenDecl))
	case *ast.BadStmt:
		return (*AstBadStmt)(astNode.(*ast.BadStmt))
	case *ast.DeclStmt:
		return (*AstDeclStmt)(astNode.(*ast.DeclStmt))
	case *ast.EmptyStmt:
		return (*AstEmptyStmt)(astNode.(*ast.EmptyStmt))
	case *ast.LabeledStmt:
		return (*AstLabeledStmt)(astNode.(*ast.LabeledStmt))
	case *ast.ExprStmt:
		return (*AstExprStmt)(astNode.(*ast.ExprStmt))
	case *ast.GoStmt:
		return (*AstGoStmt)(astNode.(*ast.GoStmt))
	case *ast.AssignStmt:
		return (*AstAssignStmt)(astNode.(*ast.AssignStmt))
	case *ast.DeferStmt:
		return (*AstDeferStmt)(astNode.(*ast.DeferStmt))
	case *ast.ReturnStmt:
		return (*AstReturnStmt)(astNode.(*ast.ReturnStmt))
	case *ast.BranchStmt:
		return (*AstBranchStmt)(astNode.(*ast.BranchStmt))
	case *ast.BlockStmt:
		return (*AstBlockStmt)(astNode.(*ast.BlockStmt))
	case *ast.IfStmt:
		return (*AstIfStmt)(astNode.(*ast.IfStmt))
	case *ast.CaseClause:
		return (*AstCaseClause)(astNode.(*ast.CaseClause))
	case *ast.SwitchStmt:
		return (*AstSwitchStmt)(astNode.(*ast.SwitchStmt))
	case *ast.IncDecStmt:
		return (*AstIncDecStmt)(astNode.(*ast.IncDecStmt))
	case *ast.TypeSwitchStmt:
		return (*AstTypeSwitchStmt)(astNode.(*ast.TypeSwitchStmt))
	case *ast.CommClause:
		return (*AstCommClause)(astNode.(*ast.CommClause))
	case *ast.SelectStmt:
		return (*AstSelectStmt)(astNode.(*ast.SelectStmt))
	case *ast.ForStmt:
		return (*AstForStmt)(astNode.(*ast.ForStmt))
	case *ast.RangeStmt:
		return (*AstRangeStmt)(astNode.(*ast.RangeStmt))
	case *ast.BadExpr:
		return (*AstBadExpr)(astNode.(*ast.BadExpr))
	case *ast.Ident:
		return (*AstIdent)(astNode.(*ast.Ident))
	case *ast.Ellipsis:
		return (*AstEllipsis)(astNode.(*ast.Ellipsis))
	case *ast.BasicLit:
		return (*AstBasicLit)(astNode.(*ast.BasicLit))
	case *ast.FuncLit:
		return (*AstFuncLit)(astNode.(*ast.FuncLit))
	case *ast.CompositeLit:
		return (*AstCompositeLit)(astNode.(*ast.CompositeLit))
	case *ast.SelectorExpr:
		return (*AstSelectorExpr)(astNode.(*ast.SelectorExpr))
	case *ast.IndexExpr:
		return (*AstIndexExpr)(astNode.(*ast.IndexExpr))
	case *ast.SliceExpr:
		return (*AstSliceExpr)(astNode.(*ast.SliceExpr))
	case *ast.TypeAssertExpr:
		return (*AstTypeAssertExpr)(astNode.(*ast.TypeAssertExpr))
	case *ast.CallExpr:
		return (*AstCallExpr)(astNode.(*ast.CallExpr))
	case *ast.StarExpr:
		return (*AstStarExpr)(astNode.(*ast.StarExpr))
	case *ast.UnaryExpr:
		return (*AstUnaryExpr)(astNode.(*ast.UnaryExpr))
	case *ast.BinaryExpr:
		return (*AstBinaryExpr)(astNode.(*ast.BinaryExpr))
	case *ast.KeyValueExpr:
		return (*AstKeyValueExpr)(astNode.(*ast.KeyValueExpr))
	case *ast.ArrayType:
		return (*AstArrayType)(astNode.(*ast.ArrayType))
	case *ast.Field:
		return (*AstField)(astNode.(*ast.Field))
	case *ast.StructType:
		return (*AstStructType)(astNode.(*ast.StructType))
	case *ast.FuncType:
		return (*AstFuncType)(astNode.(*ast.FuncType))
	case *ast.InterfaceType:
		return (*AstInterfaceType)(astNode.(*ast.InterfaceType))
	case *ast.MapType:
		return (*AstMapType)(astNode.(*ast.MapType))
	case *ast.ChanType:
		return (*AstChanType)(astNode.(*ast.ChanType))
	case *ast.ImportSpec:
		return (*AstImportSpec)(astNode.(*ast.ImportSpec))
	case *ast.ValueSpec:
		return (*AstValueSpec)(astNode.(*ast.ValueSpec))
	case *ast.TypeSpec:
		return (*AstTypeSpec)(astNode.(*ast.TypeSpec))
	}
	return nil
}

type Translator interface {
	Translate(writer writer.LuaWriter)
}

type AstFile ast.File
func (file *AstFile)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, 0, "do")
	writer.SetPackageName(file.Name.Name)

	for _, decl := range file.Decls {
		cast(decl).(Translator).Translate(writer)
	}
	writer.AppendLine(-1, 0,"if init then init() end")
	writer.AppendLine(-1, 0, "end")
	//TODO:全局函数前置声明
	/*TODO:第二版方向：
	1.同一个包名下的文件编译在同一个lua文件中。
	2.前置全局变量声明，其余按文件内顺序排列
	3.全局变量，如var，const，type，func 使用local定义，前置声明。
	4.取包内容时, 通过metatable读取，table的__newindex 设置不可写
	  __index 通过key，计算hash，通过预编译的if语句实现快速查找，如：
	  keyHash := hash(key) //可以所有变量名排序后的索引
	  注意：实现时，如果中间值为小数，去掉==项即可
	  //假设有20个
	  if keyHash > 10 {
	        if keyHash > 15 {
	         } else if keyHash == 15 {
				return var15
	         } else {

	         }
	   } else if keyHash == 10 {
	    	return var10
	   } else {
	   }
	*/
}

//Ast.Decl
type AstBadDecl ast.BadDecl
func (decl *AstBadDecl)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, decl.From, "BadDecl")
}

type AstFuncDecl ast.FuncDecl
func (decl *AstFuncDecl)Translate(w writer.LuaWriter) {
	var paramCheckers []string
	if decl.Recv == nil {
		w.AppendDef(decl.Name.Name)
		w.AppendLine(-1, decl.Name.NamePos, fmt.Sprintf("%s = function(", decl.Name.Name))
	}else {

		w.AppendLine(-1, 0, "")
		recv := decl.Recv
		cast(recv.List[0].Type).(Translator).Translate(w)

		recvName := recv.List[0].Names[0].Name
		w.AppendLine(0, decl.Name.NamePos, fmt.Sprintf(".%s = function(%s", decl.Name.Name, recvName))
		if len(decl.Type.Params.List) > 0 {
			w.AppendLine(0, 0, ",")
		}
		tempWriter := writer.LuaFile{}
		cast(recv.List[0].Type).(Translator).Translate(&tempWriter)
		paramCheckers = append(paramCheckers, fmt.Sprintf("%s = checkType(%s, %s)", recvName, recvName, tempWriter.Line(0)))
	}

	for index, param := range decl.Type.Params.List {
		for iName, name := range param.Names {
			w.AppendLine(0, name.NamePos, name.Name)
			if iName < len(param.Names) - 1 {
				w.AppendLine(0, 0, ",")
			}

			tempWriter := writer.LuaFile{}
			cast(param.Type).(Translator).Translate(&tempWriter)

			paramCheckers = append(paramCheckers, fmt.Sprintf("%s = checkType(%s, %s)", name.Name, name.Name, tempWriter.Line(0)))
		}
		if index < len(decl.Type.Params.List) - 1 {
			w.AppendLine(0, 0, ",")
		}
	}
	w.AppendLine(0, 0, ")")
	for _, checker := range paramCheckers {
		w.AppendLine(-1, 0, checker)
	}
	//检查是否有lua实现
	useLua := false
	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if strings.HasPrefix(comment.Text, "//@Lua:") {
				useLua = true
				continue
			}
			if strings.HasPrefix(comment.Text, "/*\n") &&
				strings.HasSuffix(comment.Text, "\n*/") {
				//去掉换行和注释符
				w.AppendLine(-1, comment.Pos(), comment.Text[3:len(comment.Text) - 3])
				break
			}
		}
	}
	if !useLua {
		(*AstBlockStmt)(decl.Body).Translate(w)
	}
	w.AppendLine(-1, 0, "end")
}

type AstGenDecl ast.GenDecl
func (decl *AstGenDecl)Translate(writer writer.LuaWriter) {
	for _, spec := range decl.Specs {
		switch  {
		case decl.Tok == token.IMPORT:
			packageName := ""
			importName := spec.(*ast.ImportSpec).Name
			if importName != nil {
				packageName = importName.Name
			} else {
				importPath := spec.(*ast.ImportSpec).Path.Value
				names:= strings.Split(importPath[1:len(importPath) - 1], "/")
				packageName = names[len(names) - 1]
			}
			writer.AppendLine(-1, 0, fmt.Sprintf("local %s = import(", packageName))
			cast(spec).(Translator).Translate(writer)
			writer.AppendLine(0, 0, ")")
		case decl.Tok == token.CONST || decl.Tok == token.VAR:
			names := spec.(*ast.ValueSpec).Names
			values := spec.(*ast.ValueSpec).Values
			writer.AppendLine(-1, 0, "")
			for index, name := range names {
				if writer.IsGlobalScope() {
					writer.AppendDef(name.Name)
				}
				writer.AppendLine(0, 0, name.Name)
				if index < len(names) -1 {
					writer.AppendLine(0, 0, ",")
				}
			}
			writer.AppendLine(0, 0, " = ")
			for index, value := range values {
				if writer.IsGlobalScope() {
					cast(value).(Translator).Translate(writer)
				}
				cast(value).(Translator).Translate(writer)
				if index < len(names) -1 {
					writer.AppendLine(0, 0, ",")
				}
			}
		case decl.Tok == token.TYPE:
			typeSpec := spec.(*ast.TypeSpec)
			writer.AppendDef(typeSpec.Name.Name)
			writer.AppendLine(-1, 0, fmt.Sprintf("%s = ", typeSpec.Name.Name))
			cast(typeSpec.Type).(Translator).Translate(writer)
		}
	}
}

//Ast.Stmt
type AstBadStmt ast.BadStmt
func (stmt *AstBadStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.From, "BadStmt")
}

type AstDeclStmt ast.DeclStmt
func (stmt *AstDeclStmt)Translate(writer writer.LuaWriter) {
	cast(stmt.Decl).(Translator).Translate(writer)
}
type AstEmptyStmt ast.EmptyStmt
func (stmt *AstEmptyStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Semicolon, "")
}
type AstLabeledStmt ast.LabeledStmt
func (decl *AstLabeledStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, decl.Colon, "--LabeledStmt TODO:")
}

type AstExprStmt ast.ExprStmt
func (stmt *AstExprStmt)Translate(writer writer.LuaWriter) {
	cast(stmt.X).(Translator).Translate(writer)
}
type AstIncDecStmt ast.IncDecStmt
func (stmt *AstIncDecStmt)Translate(writer writer.LuaWriter) {
	cast(stmt.X).(Translator).Translate(writer)
	writer.AppendLine(0, stmt.TokPos, " = ")
	cast(stmt.X).(Translator).Translate(writer)
	if stmt.Tok == token.DEC {
		writer.AppendLine(0, 0, " - 1")
	} else {
		writer.AppendLine(0, 0, " + 1")
	}
}
type AstAssignStmt ast.AssignStmt
func (stmt *AstAssignStmt)Translate(writer writer.LuaWriter) {
	if stmt.Tok == token.DEFINE {
		writer.AppendLine(-1, stmt.TokPos, "local ")
	}
	for _, expr := range stmt.Lhs {
		cast(expr).(Translator).Translate(writer)
	}
	writer.AppendLine(0, 0, " = ")
	for _, expr := range stmt.Rhs {
		cast(expr).(Translator).Translate(writer)
	}
	writer.AppendLine(0, 0, " ")
}
type AstGoStmt ast.GoStmt
func (stmt *AstGoStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Go, "--GoStmt TODO:")
}
type AstDeferStmt ast.DeferStmt
func (stmt *AstDeferStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Defer, "--DeferStmt TODO:")
}
type AstReturnStmt ast.ReturnStmt
func (stmt *AstReturnStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, stmt.Return, "return ")
	count := len(stmt.Results)
	for index, expr := range stmt.Results {
		cast(expr).(Translator).Translate(writer)
		if index < count - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
}

type AstBranchStmt ast.BranchStmt
func (stmt *AstBranchStmt)Translate(writer writer.LuaWriter) {
	if stmt.Tok == token.BREAK {
		writer.AppendLine(0,stmt.TokPos, " break ")
	} else {
		writer.AppendLine(0,stmt.TokPos, " goto post ")
	}
}

type AstBlockStmt ast.BlockStmt
func (stmt *AstBlockStmt)Translate(writer writer.LuaWriter) {
	writer.EnterLocalScope()
	for index, iStmt := range stmt.List {
		cast(iStmt).(Translator).Translate(writer)
		if index < len(stmt.List) - 1 {
			writer.AppendLine(-1, 0, "")
		}
	}
	writer.LeaveLocalScope()
}

type AstIfStmt ast.IfStmt
func (stmt *AstIfStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, stmt.If, "if ")
	cast(stmt.Cond).(Translator).Translate(writer)
	writer.AppendLine(0, 0, " then ")
	writer.AppendLine(-1, 0, "")
	(*AstBlockStmt)(stmt.Body).Translate(writer)

	needEnd := true //不在当前结束，由Else结束
	if stmt.Else != nil {
		switch stmt.Else.(type) {
		case *ast.IfStmt:
			writer.AppendLine(-1, stmt.Else.Pos(), "else")
			needEnd = false
		default:
			writer.AppendLine(-1, stmt.Else.Pos(), "else ")
			writer.AppendLine(-1, 0, "")
		}

		cast(stmt.Else).(Translator).Translate(writer)
	}
	if needEnd {
		writer.AppendLine(-1, 0, "end")
	}
}
type AstCaseClause ast.CaseClause
func (stmt *AstCaseClause)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Colon, "CaseClause TODO: ")
}
type AstSwitchStmt ast.SwitchStmt
func (stmt *AstSwitchStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Switch, "SwitchStmt TODO: ")
}

type AstTypeSwitchStmt ast.TypeSwitchStmt
func (stmt *AstTypeSwitchStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Switch, "TypeSwitchStmt TODO: ")
}

type AstCommClause ast.CommClause
func (stmt *AstCommClause)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Colon, "TypeSwitchStmt TODO: ")
}

type AstSelectStmt ast.SelectStmt
func (stmt *AstSelectStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.Select, "SelectStmt TODO: ")
}
type AstForStmt ast.ForStmt
func (stmt *AstForStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.For, "do ")
	cast(stmt.Init).(Translator).Translate(writer)
	writer.AppendLine(-1, 0, " while(")
	cast(stmt.Cond).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ") do ")
	writer.AppendLine(-1, 0, "")
	(*AstBlockStmt)(stmt.Body).Translate(writer)
	writer.AppendLine(-1, 0, "::post:: ")
	cast(stmt.Post).(Translator).Translate(writer)
	writer.AppendLine(0, 0, "; end")
	writer.AppendLine(0, 0, " end")
}
type AstRangeStmt ast.RangeStmt
func (stmt *AstRangeStmt)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, stmt.TokPos, "for ")
	cast(stmt.Key).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ",")
	cast(stmt.Value).(Translator).Translate(writer)
	writer.AppendLine(0, 0, " in range(")
	cast(stmt.X).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ") do \n")
	(*AstBlockStmt)(stmt.Body).Translate(writer)
	writer.AppendLine(-1, 0, "::post:: ")
	writer.AppendLine(0, 0, "; end")
}
//ast.Expr
type AstBadExpr ast.BadExpr
func (expr *AstBadExpr)Translate(writer writer.LuaWriter) {
	writer.AppendLine(-1, expr.From, "SelectStmt TODO: error here")
}
type AstIdent ast.Ident
func (ident *AstIdent)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, ident.NamePos, ident.Name)
}

type AstEllipsis ast.Ellipsis
func (ell *AstEllipsis)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, ell.Ellipsis, "...")
}

type AstBasicLit ast.BasicLit
func (basicLit *AstBasicLit)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, basicLit.ValuePos, basicLit.Value)
}

type AstFuncLit ast.FuncLit
func (funcLit *AstFuncLit)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, funcLit.Type.Pos(), "function(")
	(*AstFuncType)(funcLit.Type).Translate(writer)
	(*AstBlockStmt)(funcLit.Body).Translate(writer)
	writer.AppendLine(0, 0, "end")
}


type AstCompositeLit ast.CompositeLit
func (compositeLit *AstCompositeLit)Translate(writer writer.LuaWriter) {
	cast(compositeLit.Type).(Translator).Translate(writer)
	writer.AppendLine(0, compositeLit.Lbrace, "{")
	for _, expr := range compositeLit.Elts {
		cast(expr).(Translator).Translate(writer)
	}
	writer.AppendLine(0, compositeLit.Rbrace, "}")
}
type AstParenExpr ast.ParenExpr
func (parenExpr *AstParenExpr)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, parenExpr.Lparen, "(")
	cast(parenExpr.X).(Translator).Translate(writer)
	writer.AppendLine(0, parenExpr.Rparen, ")")
}

type AstSelectorExpr ast.SelectorExpr
func (selectExpr *AstSelectorExpr)Translate(writer writer.LuaWriter) {
	cast(selectExpr.X).(Translator).Translate(writer)
	writer.AppendLine(0, selectExpr.X.Pos(), ".")
	(*AstIdent)(selectExpr.Sel).Translate(writer)
}


type AstIndexExpr ast.IndexExpr
func (indexExpr *AstIndexExpr)Translate(writer writer.LuaWriter) {
	cast(indexExpr.X).(Translator).Translate(writer)
	writer.AppendLine(0, indexExpr.Lbrack, "[")
	cast(indexExpr.Index).(Translator).Translate(writer)
	writer.AppendLine(0, indexExpr.Rbrack, "]")
}

type AstSliceExpr ast.SliceExpr
func (sliceExpr *AstSliceExpr)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, sliceExpr.Lbrack, "slice(")
	cast(sliceExpr.X).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ",")
	cast(sliceExpr.Low).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ",")
	cast(sliceExpr.High).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ",")
	cast(sliceExpr.Max).(Translator).Translate(writer)
	writer.AppendLine(0, sliceExpr.Rbrack, ")")
}

type AstTypeAssertExpr ast.TypeAssertExpr
func (typeAssertExpr *AstTypeAssertExpr)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, typeAssertExpr.Lparen, "as(")
	cast(typeAssertExpr.X).(Translator).Translate(writer)
	writer.AppendLine(0, 0, ",")
	cast(typeAssertExpr.Type).(Translator).Translate(writer)
	writer.AppendLine(0, typeAssertExpr.Rparen, ")")
}

type AstCallExpr ast.CallExpr
func (callExpr *AstCallExpr)Translate(writer writer.LuaWriter) {
	cast(callExpr.Fun).(Translator).Translate(writer)
	writer.AppendLine(0, callExpr.Lparen, "(")
	if _, ok := callExpr.Fun.(*ast.SelectorExpr); ok{
		cast(callExpr.Fun.(*ast.SelectorExpr).X).(Translator).Translate(writer)
		if len(callExpr.Args) > 0 {
			writer.AppendLine(0, 0, ",")
		}
	}
	for index, arg := range callExpr.Args {
		cast(arg).(Translator).Translate(writer)
		if index < len(callExpr.Args) - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
	writer.AppendLine(0, callExpr.Rparen, ")")
}

type AstStarExpr ast.StarExpr
func (expr *AstStarExpr)Translate(writer writer.LuaWriter) {
	cast(expr.X).(Translator).Translate(writer)
}

type AstUnaryExpr ast.UnaryExpr
func (expr *AstUnaryExpr)Translate(writer writer.LuaWriter) {
	cast(expr.X).(Translator).Translate(writer)
}


type AstBinaryExpr ast.BinaryExpr
func (expr *AstBinaryExpr)Translate(writer writer.LuaWriter) {
	switch expr.Op {
	case token.AND:
		writer.AppendLine(0, expr.OpPos, "bit.band(")
	case token.OR:
		writer.AppendLine(0, expr.OpPos, "bit.bor(")
	case token.XOR:
		writer.AppendLine(0, expr.OpPos, "bit.bxor(")
	case token.SHL:
		writer.AppendLine(0, expr.OpPos, "bit.lshift(")
	case token.SHR:
		writer.AppendLine(0, expr.OpPos, "bit.rshift(")
	case token.AND_NOT:
		writer.AppendLine(0, expr.OpPos, "bit.bnot(bit.band(")
	case token.ADD_ASSIGN:fallthrough
	case token.SUB_ASSIGN:fallthrough
	case token.MUL_ASSIGN:fallthrough
	case token.QUO_ASSIGN:fallthrough
	case token.REM_ASSIGN:
		writer.AppendLine(0, expr.OpPos, " = ")
	}
	cast(expr.X).(Translator).Translate(writer)
	switch {
	case expr.Op == token.ADD || expr.Op == token.ADD_ASSIGN:
		writer.AppendLine(0, expr.OpPos, "+")
	case expr.Op == token.SUB || expr.Op == token.SUB_ASSIGN:
		writer.AppendLine(0, expr.OpPos, "-")
	case expr.Op == token.MUL || expr.Op == token.MUL_ASSIGN:
		writer.AppendLine(0, expr.OpPos, "*")
	case expr.Op == token.QUO || expr.Op == token.QUO_ASSIGN:
		writer.AppendLine(0, expr.OpPos, "/")
	case expr.Op == token.REM || expr.Op == token.REM_ASSIGN:
		writer.AppendLine(0, expr.OpPos, "%")
	case expr.Op == token.LAND:
		writer.AppendLine(0, expr.OpPos, " and ")
	case expr.Op == token.LOR:
		writer.AppendLine(0, expr.OpPos, " or ")
	case expr.Op == token.EQL:
		writer.AppendLine(0, expr.OpPos, " == ")
	case expr.Op == token.LSS:
		writer.AppendLine(0, expr.OpPos, " < ")
	case expr.Op == token.GTR:
		writer.AppendLine(0, expr.OpPos, " > ")
	case expr.Op == token.ASSIGN:
		writer.AppendLine(0, expr.OpPos, " = ")
	case expr.Op == token.NOT:
		writer.AppendLine(0,expr.OpPos, " not ")
	case expr.Op == token.NEQ:
		writer.AppendLine(0, expr.OpPos, " ~= ")
	case expr.Op == token.LEQ:
		writer.AppendLine(0, expr.OpPos, " <= ")
	case expr.Op == token.GEQ:
		writer.AppendLine(0, expr.OpPos, " >= ")
	case expr.Op == token.AND ||  expr.Op == token.OR ||
			expr.Op == token.XOR || expr.Op == token.SHL ||
			expr.Op == token.SHR || expr.Op == token.AND_NOT:
		writer.AppendLine(0, expr.OpPos, " , ")
	}
	cast(expr.Y).(Translator).Translate(writer)

	if expr.Op == token.AND ||  expr.Op == token.OR ||
		expr.Op == token.XOR || expr.Op == token.SHL ||
		expr.Op == token.SHR || expr.Op == token.AND_NOT {

		writer.AppendLine(0, expr.Y.End(), ")")
		if expr.Op == token.AND_NOT {
			writer.AppendLine(0, expr.Y.End(), ")")
		}
	}
}


type AstKeyValueExpr ast.KeyValueExpr
func (expr *AstKeyValueExpr)Translate(writer writer.LuaWriter) {
	cast(expr.Key).(Translator).Translate(writer)
	writer.AppendLine(0, expr.Colon, " = ")
	cast(expr.Value).(Translator).Translate(writer)
}

//ast.Type
type AstArrayType ast.ArrayType
func (arrayType *AstArrayType)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, arrayType.Lbrack, "array")
}

type AstField ast.Field
func (field *AstField)Translate(writer writer.LuaWriter)  {
	if len(field.Names) > 1 {
		writer.AppendLine(0, field.Names[0].NamePos, " error field")
	} else {
		writer.AppendLine(0, field.Names[0].NamePos, field.Names[0].Name)
	}
}
type AstFieldList ast.FieldList
func (fieldList *AstFieldList)Translate(writer writer.LuaWriter)  {
	count := len(fieldList.List)
	for index, field := range fieldList.List {
		(*AstField)(field).Translate(writer)
		if index < count - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
}
type AstStructType ast.StructType
func (structType *AstStructType)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, structType.Struct, "struct{")
	for index, field := range structType.Fields.List {
		for iName, name := range field.Names {
			writer.AppendLine(-1, name.NamePos, name.Name)
			writer.AppendLine(0, 0, "=")
			cast(field.Type).(Translator).Translate(writer)
			if iName < len(field.Names) - 1 {
				writer.AppendLine(0, 0, ",")
			}
		}
		if len(field.Names) == 0 {
			cast(field.Type).(Translator).Translate(writer)
		}

		if index < len(structType.Fields.List) - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
	writer.AppendLine(0, structType.Fields.End(), "}")
}

type AstFuncType ast.FuncType
func (funcType *AstFuncType)Translate(writer writer.LuaWriter) {
	for _, param := range funcType.Params.List {
		cast(param.Type).(Translator).Translate(writer)
		writer.AppendLine(0, 0, ",")
	}

	writer.AppendLine(0, funcType.Func, "\"=>\"")
	if funcType.Results != nil {
		for index, res := range funcType.Results.List {
			cast(res.Type).(Translator).Translate(writer)
			if index < len(funcType.Results.List) {
				writer.AppendLine(0, 0, ",")
			}
		}
	}
}


type AstInterfaceType ast.InterfaceType
func (interfaceType *AstInterfaceType)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, interfaceType.Interface, "interface{")
	for index, field := range interfaceType.Methods.List {
		for iName, name := range field.Names {
			writer.AppendLine(0, name.NamePos, name.Name)
			writer.AppendLine(0, 0, "=")
			writer.AppendLine(0, 0, "method(")
			cast(field.Type).(Translator).Translate(writer)
			writer.AppendLine(0, 0, ")")
			if iName < len(field.Names) - 1 {
				writer.AppendLine(0, 0, ",")
			}
		}
		if index < len(interfaceType.Methods.List) - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
	writer.AppendLine(0, interfaceType.Methods.End(), "}")
}

type AstMapType ast.MapType
func (decl *AstMapType)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, decl.Map, "map")
}

type AstChanType ast.ChanType
func (decl *AstChanType)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, decl.Begin, "--AstChanType todo")
}

//Spec
type AstImportSpec ast.ImportSpec
func (spec *AstImportSpec)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, spec.Path.ValuePos, strings.Replace(spec.Path.Value, "/", ".", -1))
}

type AstValueSpec ast.ValueSpec
func (spec *AstValueSpec)Translate(writer writer.LuaWriter) {
	for _, name := range spec.Names {
		writer.AppendLine(0, name.NamePos, name.Name)
		writer.AppendLine(0, 0, ",")
	}

	for index, value := range spec.Values {
		cast(value).(Translator).Translate(writer)
		if index < len(spec.Values) - 1 {
			writer.AppendLine(0, 0, ",")
		}
	}
}

type AstTypeSpec ast.TypeSpec
func (spec *AstTypeSpec)Translate(writer writer.LuaWriter) {
	writer.AppendLine(0, 0, fmt.Sprintf("\"%s\"", spec.Name.Name))
	writer.AppendLine(0, 0, ", ")
	cast(spec.Type).(Translator).Translate(writer)
}