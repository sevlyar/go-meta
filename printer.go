package meta

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
)

func Fprint(w io.Writer, node interface{}) (err error) {
	defer func() {
		err = recover().(error)
	}()

	p := &printer{w: w}
	p.printNode(node)
	return
}

type printer struct {
	w   io.Writer
	ind int
}

func (p *printer) indent() {
	p.ind++
}

func (p *printer) unindent() {
	p.ind--
}

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

func (p *printer) writeString(line string) {
	p.writeIndent()
	if _, err := p.w.Write([]byte(line)); err != nil {
		panic(err)
	}
	p.writeEol()
}

func (p *printer) print(a ...interface{}) {
	if _, err := fmt.Fprint(p.w, a); err != nil {
		panic(err)
	}
}

func (p *printer) printIndented(a ...interface{}) {
	p.writeIndent()
	p.print(a)

}

func (p *printer) printf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(p.w, format, a); err != nil {
		panic(err)
	}
}

func (p *printer) writeLine(a ...interface{}) {
	p.writeIndent()
	if _, err := fmt.Fprintln(p.w, a); err != nil {
		panic(err)
	}
}

func (p *printer) printNode(node interface{}) {
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
		p.writeLine("/*", reflect.TypeOf(n), "*/")
	}
}

func (p *printer) Exprs(exprs []ast.Expr) {
	// a, b, c
	for i, exp := range exprs {
		if i != 0 {
			p.print(", ")
		}
		p.printNode(exp)
	}
}

func (p *printer) Stmts(stmts []ast.Stmt) {
	p.indent()
	for _, item := range stmts {
		p.writeIndent()
		p.printNode(item)
	}
	p.unindent()
}

func (p *printer) Names(names []*ast.Ident) {
	for i, id := range names {
		if i != 0 {
			p.print(", ")
		}
		p.Ident(id)
	}
}

func (p *printer) ArrayType(t *ast.ArrayType) {
	// []Elt
	// [...]Elt
	p.print("[")
	if t.Len != nil {
		p.printNode(t.Len)
	}
	p.print("]")
	p.printNode(t.Elt)
}

func (p *printer) AssignStmt(stmt *ast.AssignStmt) {
	// a, b := x, y
	p.Exprs(stmt.Lhs)
	p.printf(" %s ", stmt.Tok.String())
	p.Exprs(stmt.Rhs)
}

func (p *printer) BasicLit(lit *ast.BasicLit) {
	p.print(lit.Value)
}

func (p *printer) BinaryExpr(expr *ast.BinaryExpr) {
	p.printNode(expr.X)
	p.print(expr.Op.String())
	p.printNode(expr.Y)
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

func (p *printer) CallExpr(expr *ast.CallExpr) {
	p.printNode(expr.Fun)
	p.print("(")
	p.Exprs(expr.Args)
	// TODO: expr.Ellipsis?
	p.print(")")
}

func (p *printer) CaseClause(clause *ast.CaseClause) {
	if clause.List != nil {
		p.Exprs(clause.List)
		p.print(": ")
	} else {
		p.print("default: ")
	}
	p.Stmts(clause.Body)
}

func (p *printer) ChanType(t *ast.ChanType) {
	if t.Arrow != token.NoPos && t.Dir == ast.RECV {
		p.print("<-")
	}
	p.print("chan")
	if t.Arrow != token.NoPos && t.Dir == ast.SEND {
		p.print("<-")
	}
	p.printNode(t.Value)
}

func (p *printer) CommClause(clause *ast.CommClause) {
	if clause.Comm != nil {
		p.printNode(clause.Comm)
		p.print(": ")
	} else {
		p.print("default: ")
	}
	p.Stmts(clause.Body)
}

func (p *printer) Comment(comment *ast.Comment) {
	p.print(comment.Text)
}

func (p *printer) CommentGroup(group *ast.CommentGroup) {
	for _, comment := range group.List {
		p.writeIndent()
		p.Comment(comment)
	}
}

func (p *printer) CompositeLit(lit *ast.CompositeLit) {
	p.printNode(lit.Type)
	p.print(" {")
	p.Exprs(lit.Elts)
	p.printIndented("}")
}

func (p *printer) DeclStmt(stmt *ast.DeclStmt) {
	p.printNode(stmt.Decl)
}

func (p *printer) DeferStmt(stmt *ast.DeferStmt) {
	p.print("defer ")
	p.printNode(stmt.Call)
}

func (p *printer) Ellipsis(el *ast.Ellipsis) {
	p.print("...")
	p.printNode(el.Elt)
}

func (p *printer) EmptyStmt(stmt *ast.EmptyStmt) {
	p.print(";")
}

func (p *printer) ExprStmt(stmt *ast.ExprStmt) {
	p.printNode(stmt.X)
}

func (p *printer) Field(fld *ast.Field) {
	// TODO: append printing of comments
	if fld.Names != nil {
		p.Names(fld.Names)
		p.print(" ")
	}

	p.printNode(fld.Type)

	if fld.Tag != nil {
		p.printNode(fld.Tag)
	}
}

func (p *printer) FieldList(list *ast.FieldList) {
	p.fieldList(list, false)
}

func (p *printer) fieldList(list *ast.FieldList, oneline bool) {
	if list.List == nil {
		return
	}
	if !oneline {
		p.indent()
		defer p.unindent()
	}
	for i, item := range list.List {
		if oneline {
			if i != 0 {
				p.print(", ")
			}
		} else {
			p.writeIndent()
		}
		p.Field(item)
	}
}

func (p *printer) File(file *ast.File) {
	p.CommentGroup(file.Doc)
	p.print("package", file.Name.Name)
	p.writeEol()

	// TODO: extract in func
	p.print("import (")
	for _, imp := range file.Imports {
		p.writeIndent()
		p.ImportSpec(imp)
	}
	p.printIndented(")")
	p.writeEol()

	for _, decl := range file.Decls {
		p.writeEol()
		p.printNode(decl)
		p.writeEol()
	}
}

func (p *printer) ForStmt(v *ast.ForStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) FuncDecl(v *ast.FuncDecl) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) FuncLit(v *ast.FuncLit) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) FuncType(v *ast.FuncType) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) GenDecl(v *ast.GenDecl) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) GoStmt(v *ast.GoStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) Ident(v *ast.Ident) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) IfStmt(v *ast.IfStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) ImportSpec(v *ast.ImportSpec) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) IncDecStmt(v *ast.IncDecStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) IndexExpr(v *ast.IndexExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) InterfaceType(v *ast.InterfaceType) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) KeyValueExpr(v *ast.KeyValueExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) LabeledStmt(v *ast.LabeledStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) MapType(v *ast.MapType) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) ParenExpr(v *ast.ParenExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) RangeStmt(v *ast.RangeStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) ReturnStmt(v *ast.ReturnStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) SelectStmt(v *ast.SelectStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) SelectorExpr(v *ast.SelectorExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) SendStmt(v *ast.SendStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) SliceExpr(v *ast.SliceExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) StarExpr(v *ast.StarExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) StructType(v *ast.StructType) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) SwitchStmt(v *ast.SwitchStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) TypeAssertExpr(v *ast.TypeAssertExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) TypeSpec(v *ast.TypeSpec) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) TypeSwitchStmt(v *ast.TypeSwitchStmt) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) UnaryExpr(v *ast.UnaryExpr) {
	p.print("/***", reflect.TypeOf(v), "*/")
}
func (p *printer) ValueSpec(v *ast.ValueSpec) {
	p.print("/***", reflect.TypeOf(v), "*/")
}

/*
func (p *printer) printFile(file *ast.File) {
	p.printCommentGroup(file.Doc)

	p.writeLine("package", file.Name.Name)

	p.writeString("import (")
	p.indent()
	for _, imp := range file.Imports {
		p.writeIndent()
		if imp.Name != nil {
			p.write(imp.Name.Name)
		}
		p.write(imp.Path.Value)
		p.writeEol()
	}
	p.unindent()
	p.writeString(")")

	if file.Decls != nil {
		for _, decl := range file.Decls {
			p.fprint(decl)
		}
	}
}
*/
