package meta

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"text/tabwriter"
)

func Fprint(w io.Writer, node interface{}) (err error) {
	// defer func() {
	// 	err = recover().(error)
	// 	panic(err)
	// }()

	p := &printer{w, w, tabwriter.NewWriter(w, 1, 4, 1, ' ', tabwriter.TabIndent), 0, " "}
	p.Node(node)
	return nil
}

type printer struct {
	w     io.Writer
	cpw   io.Writer
	tw    *tabwriter.Writer
	ind   int
	space string
}

func (p *printer) indent() {
	p.ind++
}

func (p *printer) unindent() {
	p.ind--
}

func (p *printer) twMode() {
	p.w = p.tw
	p.space = "\t"
}

func (p *printer) normalMode() {
	p.w = p.cpw
	p.space = " "
	if err := p.tw.Flush(); err != nil {
		panic(err)
	}
}

func (p *printer) tapMode() {
	p.normalMode()
	p.twMode()
}

// TODO: change this
var indents = []byte("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t")

func (p *printer) writeIndent() {
	p.writeEol()
	if p.ind > 0 {
		if _, err := p.w.Write(indents[:p.ind]); err != nil {
			panic(err)
		}
	}
}

var eol = []byte{'\n'}

func (p *printer) writeEol() {
	if _, err := p.w.Write(eol); err != nil {
		panic(err)
	}
}

func (p *printer) print(a ...interface{}) {
	if _, err := fmt.Fprint(p.w, a...); err != nil {
		panic(err)
	}
}

func (p *printer) printIndented(a ...interface{}) {
	p.writeIndent()
	p.print(a...)

}

func (p *printer) unexpected(s string, node interface{}) {
	p.print(" /* ", s, " ", reflect.TypeOf(node), " */ ")
}

func (p *printer) Node(node interface{}) {
	switch n := node.(type) {
	case *ast.ArrayType:
		p.ArrayType(n)
	case *ast.AssignStmt:
		p.AssignStmt(n)
	case *ast.BasicLit:
		p.BasicLit(n)
	case *ast.BinaryExpr:
		p.BinaryExpr(n)
	case *ast.BlockStmt:
		p.BlockStmt(n)
	case *ast.BranchStmt:
		p.BranchStmt(n)
	case *ast.CallExpr:
		p.CallExpr(n)
	case *ast.CaseClause:
		p.CaseClause(n)
	case *ast.ChanType:
		p.ChanType(n)
	case *ast.CommClause:
		p.CommClause(n)
	case *ast.Comment:
		p.Comment(n)
	case *ast.CommentGroup:
		p.CommentGroup(n)
	case *ast.CompositeLit:
		p.CompositeLit(n)
	case *ast.DeclStmt:
		p.DeclStmt(n)
	case *ast.DeferStmt:
		p.DeferStmt(n)
	case *ast.Ellipsis:
		p.Ellipsis(n)
	case *ast.EmptyStmt:
		p.EmptyStmt(n)
	case *ast.ExprStmt:
		p.ExprStmt(n)
	case *ast.Field:
		p.Field(n)
	case *ast.FieldList:
		p.FieldList(n)
	case *ast.File:
		p.File(n)
	case *ast.ForStmt:
		p.ForStmt(n)
	case *ast.FuncDecl:
		p.FuncDecl(n)
	case *ast.FuncLit:
		p.FuncLit(n)
	case *ast.FuncType:
		p.FuncType(n)
	case *ast.GenDecl:
		p.GenDecl(n)
	case *ast.GoStmt:
		p.GoStmt(n)
	case *ast.Ident:
		p.Ident(n)
	case *ast.IfStmt:
		p.IfStmt(n)
	case *ast.ImportSpec:
		p.ImportSpec(n)
	case *ast.IncDecStmt:
		p.IncDecStmt(n)
	case *ast.IndexExpr:
		p.IndexExpr(n)
	case *ast.InterfaceType:
		p.InterfaceType(n)
	case *ast.KeyValueExpr:
		p.KeyValueExpr(n)
	case *ast.LabeledStmt:
		p.LabeledStmt(n)
	case *ast.MapType:
		p.MapType(n)
	case *ast.ParenExpr:
		p.ParenExpr(n)
	case *ast.RangeStmt:
		p.RangeStmt(n)
	case *ast.ReturnStmt:
		p.ReturnStmt(n)
	case *ast.SelectStmt:
		p.SelectStmt(n)
	case *ast.SelectorExpr:
		p.SelectorExpr(n)
	case *ast.SendStmt:
		p.SendStmt(n)
	case *ast.SliceExpr:
		p.SliceExpr(n)
	case *ast.StarExpr:
		p.StarExpr(n)
	case *ast.StructType:
		p.StructType(n)
	case *ast.SwitchStmt:
		p.SwitchStmt(n)
	case *ast.TypeAssertExpr:
		p.TypeAssertExpr(n)
	case *ast.TypeSpec:
		p.TypeSpec(n)
	case *ast.TypeSwitchStmt:
		p.TypeSwitchStmt(n)
	case *ast.UnaryExpr:
		p.UnaryExpr(n)
	case *ast.ValueSpec:
		p.ValueSpec(n)
	default:
		p.unexpected("Node", n)
	}
}

//****************************************************************************

func (p *printer) File(file *ast.File) {
	if file.Doc != nil {
		p.CommentGroup(file.Doc)
		p.writeEol()
	}
	p.print("package ", file.Name.Name)
	p.writeEol()

	for _, decl := range file.Decls {
		p.writeEol()
		p.Decl(decl)
	}
}

func (p *printer) Comment(comment *ast.Comment) {
	p.print(comment.Text)
}

func (p *printer) CommentGroup(group *ast.CommentGroup) {
	for i, comment := range group.List {
		if i != 0 {
			p.writeIndent()
		}
		p.Comment(comment)
	}
}

func (p *printer) Field(fld *ast.Field) {
	if fld.Doc != nil {
		p.CommentGroup(fld.Doc)
		p.tapMode()
		p.writeIndent()
	}

	if fld.Names != nil {
		p.Names(fld.Names)
		if _, ok := fld.Type.(*ast.FuncType); !ok {
			p.print(p.space)
		}
	}

	p.Expr(fld.Type)

	if fld.Tag != nil {
		p.print(p.space)
		p.Expr(fld.Tag)
	}

	if fld.Comment != nil {
		p.print(p.space)
		p.CommentGroup(fld.Comment)
	}
}

func fieldHaveComments(fld *ast.Field) bool {
	return fld.Doc != nil || fld.Comment != nil
}

func (p *printer) FieldList(list *ast.FieldList) {
	p.fieldList(list)
}

func (p *printer) fieldList(list *ast.FieldList) {
	if list.List == nil {
		p.print("{}")
		return
	}

	p.print(" { ")

	if len(list.List) < 2 && !fieldHaveComments(list.List[0]) {
		p.Field(list.List[0])
		p.print(" }")
		return
	}

	p.indent()
	p.twMode()
	for i, item := range list.List {
		if i != 0 && item.Doc != nil {
			p.writeIndent()
		}
		p.writeIndent()
		p.Field(item)
	}
	p.normalMode()
	p.unindent()
	p.printIndented("}")
}

func (p *printer) paramsList(list *ast.FieldList) {
	p.print("(")
	if list.List != nil {
		for i, item := range list.List {
			if i != 0 {
				p.print(", ")
			}
			p.Field(item)
		}
	}
	p.print(")")
}

//****************************************************************************

func (p *printer) Exprs(exprs []ast.Expr) {
	for i, exp := range exprs {
		if i != 0 {
			p.print(", ")
		}
		p.Expr(exp)
	}
}

func (p *printer) Expr(x ast.Expr) {
	switch expr := x.(type) {
	case *ast.Ident:
		p.Ident(expr)
	case *ast.Ellipsis:
		p.Ellipsis(expr)
	case *ast.BasicLit:
		p.BasicLit(expr)
	case *ast.FuncLit:
		p.FuncLit(expr)
	case *ast.CompositeLit:
		p.CompositeLit(expr)
	case *ast.CallExpr:
		p.CallExpr(expr)
	case *ast.StarExpr:
		p.StarExpr(expr)
	case *ast.IndexExpr:
		p.IndexExpr(expr)
	case *ast.ParenExpr:
		p.ParenExpr(expr)
	case *ast.SliceExpr:
		p.SliceExpr(expr)
	case *ast.UnaryExpr:
		p.UnaryExpr(expr)
	case *ast.BinaryExpr:
		p.BinaryExpr(expr)
	case *ast.SelectorExpr:
		p.SelectorExpr(expr)
	case *ast.KeyValueExpr:
		p.KeyValueExpr(expr)
	case *ast.TypeAssertExpr:
		p.TypeAssertExpr(expr)
	case *ast.ArrayType:
		p.ArrayType(expr)
	case *ast.StructType:
		p.StructType(expr)
	case *ast.FuncType:
		p.FuncType(expr)
	case *ast.InterfaceType:
		p.InterfaceType(expr)
	case *ast.MapType:
		p.MapType(expr)
	case *ast.ChanType:
		p.ChanType(expr)
	default:
		p.unexpected("Expr", x)
	}
}

// TODO: append BadExpr

func (p *printer) Names(names []*ast.Ident) {
	for i, id := range names {
		if i != 0 {
			p.print(", ")
		}
		p.Ident(id)
	}
}

func (p *printer) Ident(id *ast.Ident) {
	p.print(id.Name)
}

func (p *printer) Ellipsis(el *ast.Ellipsis) {
	p.print("...")
	p.Node(el.Elt)
}

func (p *printer) BasicLit(lit *ast.BasicLit) {
	p.print(lit.Value)
}

func (p *printer) FuncLit(lit *ast.FuncLit) {
	p.FuncType(lit.Type)
	p.BlockStmt(lit.Body)
}

func (p *printer) CompositeLit(lit *ast.CompositeLit) {
	p.Node(lit.Type)
	p.print("{")
	p.Exprs(lit.Elts)
	p.print("}")
}

func (p *printer) CallExpr(expr *ast.CallExpr) {
	p.Expr(expr.Fun)
	p.print("(")
	p.Exprs(expr.Args)
	if expr.Ellipsis.IsValid() {
		p.print("...")
	}
	p.print(")")
}

func (p *printer) ParenExpr(expr *ast.ParenExpr) {
	p.print("(")
	p.Expr(expr.X)
	p.print(")")
}

func (p *printer) StarExpr(expr *ast.StarExpr) {
	p.print("*")
	p.Expr(expr.X)
}

func (p *printer) BinaryExpr(expr *ast.BinaryExpr) {
	p.Expr(expr.X)
	p.print(" ", expr.Op.String(), " ")
	p.Expr(expr.Y)
}

func (p *printer) IndexExpr(expr *ast.IndexExpr) {
	p.Expr(expr.X)
	p.print("[")
	p.Expr(expr.Index)
	p.print("]")
}

func (p *printer) KeyValueExpr(expr *ast.KeyValueExpr) {
	p.Expr(expr.Key)
	p.print(": ")
	p.Expr(expr.Value)
}

func (p *printer) SelectorExpr(expr *ast.SelectorExpr) {
	p.Expr(expr.X)
	p.print(".")
	p.Ident(expr.Sel)
}

func (p *printer) SliceExpr(expr *ast.SliceExpr) {
	p.Expr(expr.X)
	p.print("[")
	if expr.Low != nil {
		p.Expr(expr.Low)
	}
	p.print(":")
	if expr.High != nil {
		p.Expr(expr.High)
	}
	if expr.Slice3 {
		p.print(":")
		if expr.Max != nil {
			p.Expr(expr.Max)
		}
	}
	p.print("]")
}

func (p *printer) TypeAssertExpr(expr *ast.TypeAssertExpr) {
	p.Expr(expr.X)
	p.print(".(")
	if expr.Type != nil {
		p.Expr(expr.Type)
	} else {
		p.print("type")
	}
	p.print(")")
}

func (p *printer) UnaryExpr(expr *ast.UnaryExpr) {
	p.print(expr.Op.String())
	p.Expr(expr.X)
}

//****************************************************************************

func (p *printer) ArrayType(t *ast.ArrayType) {
	p.print("[")
	if t.Len != nil {
		p.Expr(t.Len)
	}
	p.print("]")
	p.Expr(t.Elt)
}

func (p *printer) ChanType(t *ast.ChanType) {
	if t.Arrow != token.NoPos && t.Dir == ast.RECV {
		p.print("<-")
	}
	p.print("chan")
	if t.Arrow != token.NoPos && t.Dir == ast.SEND {
		p.print("<-")
	}
	p.Expr(t.Value)
}

func (p *printer) FuncType(t *ast.FuncType) {
	p.paramsList(t.Params)
	if t.Results != nil {
		p.print(" ")
		p.paramsList(t.Results)
	}
}

func (p *printer) MapType(t *ast.MapType) {
	p.print("map[")
	p.Expr(t.Key)
	p.print("]")
	p.Expr(t.Value)
}

func (p *printer) InterfaceType(t *ast.InterfaceType) {
	p.print("interface")
	p.fieldList(t.Methods)
}

func (p *printer) StructType(t *ast.StructType) {
	p.print("struct")
	p.fieldList(t.Fields)
}

//****************************************************************************

func (p *printer) Stmts(stmts []ast.Stmt) {
	p.indent()
	for _, item := range stmts {
		if _, ok := item.(*ast.LabeledStmt); ok {
			p.writeEol()
		} else {
			p.writeIndent()
		}
		p.Stmt(item)
	}
	p.unindent()
}

func (p *printer) Stmt(s ast.Stmt) {
	switch stmt := s.(type) {
	case *ast.DeclStmt:
		p.DeclStmt(stmt)
	case *ast.EmptyStmt:
		p.EmptyStmt(stmt)
	case *ast.LabeledStmt:
		p.LabeledStmt(stmt)
	case *ast.ExprStmt:
		p.ExprStmt(stmt)
	case *ast.SendStmt:
		p.SendStmt(stmt)
	case *ast.IncDecStmt:
		p.IncDecStmt(stmt)
	case *ast.AssignStmt:
		p.AssignStmt(stmt)
	case *ast.GoStmt:
		p.GoStmt(stmt)
	case *ast.DeferStmt:
		p.DeferStmt(stmt)
	case *ast.ReturnStmt:
		p.ReturnStmt(stmt)
	case *ast.BranchStmt:
		p.BranchStmt(stmt)
	case *ast.BlockStmt:
		p.BlockStmt(stmt)
	case *ast.IfStmt:
		p.IfStmt(stmt)
	case *ast.CaseClause:
		p.CaseClause(stmt)
	case *ast.SwitchStmt:
		p.SwitchStmt(stmt)
	case *ast.TypeSwitchStmt:
		p.TypeSwitchStmt(stmt)
	case *ast.CommClause:
		p.CommClause(stmt)
	case *ast.SelectStmt:
		p.SelectStmt(stmt)
	case *ast.ForStmt:
		p.ForStmt(stmt)
	case *ast.RangeStmt:
		p.RangeStmt(stmt)
	default:
		p.unexpected("Stmt", s)
	}
}

// TODO: BadStmt

func (p *printer) AssignStmt(stmt *ast.AssignStmt) {
	p.Exprs(stmt.Lhs)
	p.print(" ", stmt.Tok.String(), " ")
	p.Exprs(stmt.Rhs)
}

func (p *printer) BlockStmt(stmt *ast.BlockStmt) {
	p.print("{")
	p.Stmts(stmt.List)
	p.printIndented("}")
}

func (p *printer) BranchStmt(stmt *ast.BranchStmt) {
	p.print(stmt.Tok.String())
	if stmt.Label != nil {
		p.print(stmt.Label.Name)
	}
}

func (p *printer) CaseClause(clause *ast.CaseClause) {
	if clause.List != nil {
		p.print("case ")
		p.Exprs(clause.List)
		p.print(": ")
	} else {
		p.print("default: ")
	}
	p.Stmts(clause.Body)
}

func (p *printer) CommClause(clause *ast.CommClause) {
	if clause.Comm != nil {
		p.print("case ")
		p.Stmt(clause.Comm)
		p.print(": ")
	} else {
		p.print("default: ")
	}
	p.Stmts(clause.Body)
}

func (p *printer) DeclStmt(stmt *ast.DeclStmt) {
	p.Decl(stmt.Decl)
}

func (p *printer) DeferStmt(stmt *ast.DeferStmt) {
	p.print("defer ")
	p.CallExpr(stmt.Call)
}

func (p *printer) EmptyStmt(stmt *ast.EmptyStmt) {
	p.print(";")
}

func (p *printer) ExprStmt(stmt *ast.ExprStmt) {
	p.Expr(stmt.X)
}

func (p *printer) ForStmt(stmt *ast.ForStmt) {
	p.print("for")
	if stmt.Init != nil {
		p.print(" ")
		p.Stmt(stmt.Init)
		p.print(";")
	}
	if stmt.Cond != nil {
		p.print(" ")
		p.Expr(stmt.Cond)
	}
	if stmt.Post != nil {
		p.print("; ")
		p.Stmt(stmt.Post)
	}

	p.print(" ")
	p.BlockStmt(stmt.Body)
}

func (p *printer) GoStmt(stmt *ast.GoStmt) {
	p.print("go ")
	p.CallExpr(stmt.Call)
}

func (p *printer) IfStmt(stmt *ast.IfStmt) {
	p.print("if ")
	if stmt.Init != nil {
		p.Stmt(stmt.Init)
		p.print("; ")
	}
	p.Expr(stmt.Cond)
	p.print(" ")
	p.BlockStmt(stmt.Body)
	if stmt.Else != nil {
		p.print(" else ")
		p.Stmt(stmt.Else)
	}
}

func (p *printer) IncDecStmt(stmt *ast.IncDecStmt) {
	p.Expr(stmt.X)
	p.print(stmt.Tok.String())
}

func (p *printer) LabeledStmt(stmt *ast.LabeledStmt) {
	p.Ident(stmt.Label)
	p.print(":")
	p.writeIndent()
	p.Stmt(stmt.Stmt)
}

func (p *printer) RangeStmt(stmt *ast.RangeStmt) {
	p.print("for ")
	p.Expr(stmt.Key)
	if stmt.Value != nil {
		p.print(", ")
		p.Expr(stmt.Value)
	}
	p.print(" ", stmt.Tok.String(), " range ")
	p.Expr(stmt.X)
	p.print(" ")
	p.BlockStmt(stmt.Body)
}

func (p *printer) ReturnStmt(stmt *ast.ReturnStmt) {
	p.print("return ")
	if stmt.Results != nil {
		p.Exprs(stmt.Results)
	}
}

func (p *printer) SelectStmt(stmt *ast.SelectStmt) {
	p.print("select ")
	p.BlockStmt(stmt.Body)
}

func (p *printer) SendStmt(stmt *ast.SendStmt) {
	p.Expr(stmt.Chan)
	p.print(" <- ")
	p.Expr(stmt.Value)
}

func (p *printer) SwitchStmt(stmt *ast.SwitchStmt) {
	p.print("switch ")
	if stmt.Init != nil {
		p.Stmt(stmt.Init)
		p.print("; ")
	}
	if stmt.Tag != nil {
		p.Expr(stmt.Tag)
	}
	p.print(" ")
	p.BlockStmt(stmt.Body)
}

func (p *printer) TypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
	p.print("switch ")
	if stmt.Init != nil {
		p.Stmt(stmt.Init)
		p.print("; ")
	}
	p.Stmt(stmt.Assign)
	p.print(" ")
	p.BlockStmt(stmt.Body)
}

//****************************************************************************

func (p *printer) Spec(s ast.Spec) {
	switch spec := s.(type) {
	case *ast.ValueSpec:
		p.ValueSpec(spec)
	case *ast.ImportSpec:
		p.ImportSpec(spec)
	case *ast.TypeSpec:
		p.TypeSpec(spec)
	default:
		p.unexpected("Stmt", s)
	}
}

func specsDoc(s ast.Spec) *ast.CommentGroup {
	switch spec := s.(type) {
	case *ast.ValueSpec:
		return spec.Doc
	case *ast.ImportSpec:
		return spec.Doc
	case *ast.TypeSpec:
		return spec.Doc
	}
	return nil
}

func (p *printer) ValueSpec(spec *ast.ValueSpec) {
	if spec.Doc != nil {
		p.CommentGroup(spec.Doc)
		p.tapMode()
		p.writeIndent()
	}

	p.Names(spec.Names)
	if spec.Type != nil {
		p.print(p.space)
		p.Expr(spec.Type)
	}

	if spec.Values != nil {
		p.print(" = ")
		p.Exprs(spec.Values)
	}

	if spec.Comment != nil {
		p.print(p.space)
		p.CommentGroup(spec.Comment)
	}
}

func (p *printer) ImportSpec(spec *ast.ImportSpec) {
	if spec.Doc != nil {
		p.CommentGroup(spec.Doc)
		p.writeIndent()
	}

	if spec.Name != nil {
		p.Ident(spec.Name)
		p.print(" ")
	}
	p.BasicLit(spec.Path)

	if spec.Comment != nil {
		p.print(" ")
		p.CommentGroup(spec.Comment)
	}
}

func (p *printer) TypeSpec(spec *ast.TypeSpec) {
	if spec.Doc != nil {
		p.CommentGroup(spec.Doc)
		p.tapMode()
		p.writeIndent()
	}

	p.Ident(spec.Name)
	p.print(p.space)
	p.Expr(spec.Type)

	if spec.Comment != nil {
		p.print(p.space)
		p.CommentGroup(spec.Comment)
	}
}

//****************************************************************************

func (p *printer) Decl(d ast.Decl) {
	switch decl := d.(type) {
	case *ast.FuncDecl:
		p.FuncDecl(decl)
	case *ast.GenDecl:
		p.GenDecl(decl)
	default:
		p.unexpected("Decl", d)
	}
}

// TODO: BadDecl

func (p *printer) FuncDecl(decl *ast.FuncDecl) {
	if decl.Doc != nil {
		p.writeIndent()
		p.CommentGroup(decl.Doc)
		p.writeIndent()
	}

	p.print("func ")
	if decl.Recv != nil {
		p.paramsList(decl.Recv)
	}
	p.Ident(decl.Name)
	p.FuncType(decl.Type)
	p.print(" ")
	p.BlockStmt(decl.Body)
}

func (p *printer) GenDecl(decl *ast.GenDecl) {
	if decl.Doc != nil {
		p.writeIndent()
		p.CommentGroup(decl.Doc)
		p.writeIndent()
	}

	p.print(decl.Tok.String(), " ")
	if len(decl.Specs) > 1 {
		p.print("(")
		p.indent()
		p.twMode()
		for i, spec := range decl.Specs {
			if i != 0 && specsDoc(spec) != nil {
				p.writeIndent()
			}
			p.writeIndent()
			p.Spec(spec)
		}
		p.normalMode()
		p.unindent()
		p.printIndented(")")
	} else {
		p.Spec(decl.Specs[0])
	}
}
