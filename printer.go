package meta

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
)

func Fprint(w io.Writer, node interface{}) (err error) {
	// defer func() {
	// 	err = recover().(error)
	// 	panic(err)
	// }()

	p := &printer{w: w}
	p.printNode(node)
	return nil
}

// TODO: Group methods. Append methods Stmt() Expr() Decl() Type().
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
	if _, err := fmt.Fprint(p.w, a...); err != nil {
		panic(err)
	}
}

func (p *printer) printIndented(a ...interface{}) {
	p.writeIndent()
	p.print(a...)

}

func (p *printer) printf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(p.w, format, a...); err != nil {
		panic(err)
	}
}

func (p *printer) writeLine(a ...interface{}) {
	p.writeIndent()
	if _, err := fmt.Fprintln(p.w, a...); err != nil {
		panic(err)
	}
}

// TODO: rename to Node
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

func (p *printer) paramsList(list *ast.FieldList) {
	p.print("(")
	p.fieldList(list, true)
	p.print(")")
}

func (p *printer) File(file *ast.File) {
	if file.Doc != nil {
		p.CommentGroup(file.Doc)
	}
	p.print("package ", file.Name.Name)
	p.writeEol()

	for _, decl := range file.Decls {
		//p.writeEol()
		p.printNode(decl)
		p.writeEol()
	}
}

func (p *printer) ForStmt(stmt *ast.ForStmt) {
	p.printIndented("for ")
	if stmt.Init != nil {
		p.printNode(stmt.Init)
		p.print("; ")
	}
	if stmt.Cond != nil {
		p.printNode(stmt.Cond)
	}
	if stmt.Post != nil {
		p.print("; ")
		p.printNode(stmt.Post)
	}

	p.BlockStmt(stmt.Body)
}

func (p *printer) FuncDecl(decl *ast.FuncDecl) {
	if decl.Doc != nil {
		p.CommentGroup(decl.Doc)
		p.writeIndent()
	}
	p.print("func ")
	if decl.Recv != nil {
		p.paramsList(decl.Recv)
	}
	p.Ident(decl.Name)
	p.FuncType(decl.Type)
	p.BlockStmt(decl.Body)
}

func (p *printer) FuncLit(lit *ast.FuncLit) {
	p.FuncType(lit.Type)
	p.BlockStmt(lit.Body)
}

func (p *printer) FuncType(t *ast.FuncType) {
	p.paramsList(t.Params)
	p.print(" ")
	if t.Results != nil {
		p.paramsList(t.Results)
	}
}

func (p *printer) GenDecl(decl *ast.GenDecl) {
	if decl.Doc != nil {
		p.CommentGroup(decl.Doc)
	}
	p.printIndented(decl.Tok.String(), " (")
	p.indent()
	for _, spec := range decl.Specs {
		p.writeIndent()
		p.printNode(spec)
	}
	p.unindent()
	p.printIndented(")")
}

func (p *printer) GoStmt(stmt *ast.GoStmt) {
	p.print("go ")
	p.CallExpr(stmt.Call)
}

func (p *printer) Ident(id *ast.Ident) {
	p.print(id.Name)
}

func (p *printer) IfStmt(stmt *ast.IfStmt) {
	p.printIndented("if ")
	if stmt.Init != nil {
		p.printNode(stmt.Init)
		p.print("; ")
	}
	p.printNode(stmt.Cond)
	p.BlockStmt(stmt.Body)
	if stmt.Else != nil {
		p.printIndented(" else ")
		p.printNode(stmt.Else)
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
		p.CommentGroup(spec.Comment)
	}
}

func (p *printer) IncDecStmt(stmt *ast.IncDecStmt) {
	p.printNode(stmt.X)
	p.print(stmt.Tok.String())
}

func (p *printer) IndexExpr(expr *ast.IndexExpr) {
	p.printNode(expr.X)
	p.print("[")
	p.printNode(expr.Index)
	p.print("]")
}

func (p *printer) InterfaceType(t *ast.InterfaceType) {
	p.print("interface {")
	p.indent()
	p.fieldList(t.Methods, false)
	p.unindent()
	p.printIndented("}")
}

func (p *printer) KeyValueExpr(expr *ast.KeyValueExpr) {
	p.printNode(expr.Key)
	p.print(": ")
	p.printNode(expr.Value)
}

func (p *printer) LabeledStmt(stmt *ast.LabeledStmt) {
	p.Ident(stmt.Label)
	p.print(":")
	p.printNode(stmt.Stmt)
}

func (p *printer) MapType(t *ast.MapType) {
	p.print("map[")
	p.printNode(t.Key)
	p.print("]")
	p.printNode(t.Value)
}

func (p *printer) ParenExpr(expr *ast.ParenExpr) {
	p.print("(")
	p.printNode(expr.X)
	p.print(")")
}

func (p *printer) RangeStmt(stmt *ast.RangeStmt) {
	p.print("for ")
	p.printNode(stmt.Key)
	if stmt.Value != nil {
		p.print(", ")
		p.printNode(stmt.Value)
	}
	p.print(" ", stmt.Tok.String(), " range")
	p.printNode(stmt.X)
	p.BlockStmt(stmt.Body)
}

func (p *printer) ReturnStmt(stmt *ast.ReturnStmt) {
	p.print("return ")
	if stmt.Results != nil {
		p.Exprs(stmt.Results)
	}
}

func (p *printer) SelectStmt(stmt *ast.SelectStmt) {
	p.print(" select ")
	p.BlockStmt(stmt.Body)
}

func (p *printer) SelectorExpr(expr *ast.SelectorExpr) {
	p.printNode(expr.X)
	p.print(".")
	p.Ident(expr.Sel)
}

func (p *printer) SendStmt(stmt *ast.SendStmt) {
	p.printNode(stmt.Chan)
	p.print(" <- ")
	p.printNode(stmt.Value)
}

func (p *printer) SliceExpr(expr *ast.SliceExpr) {
	p.printNode(expr.X)
	p.print("[")
	if expr.Low != nil {
		p.printNode(expr.Low)
	}
	p.print(":")
	if expr.High != nil {
		p.printNode(expr.High)
	}
	if expr.Slice3 {
		p.print(":")
		if expr.Max != nil {
			p.printNode(expr.Max)
		}
	}
	p.print("]")
}

func (p *printer) StarExpr(expr *ast.StarExpr) {
	p.print("*")
	p.printNode(expr.X)
}

func (p *printer) StructType(t *ast.StructType) {
	p.print("struct {")
	p.indent()
	p.fieldList(t.Fields, false)
	p.unindent()
	p.printIndented("}")
}

func (p *printer) SwitchStmt(stmt *ast.SwitchStmt) {
	p.print("switch ")
	if stmt.Init != nil {
		p.printNode(stmt.Init)
		p.print("; ")
	}
	if stmt.Tag != nil {
		p.printNode(stmt.Tag)
	}
	p.BlockStmt(stmt.Body)
}

func (p *printer) TypeAssertExpr(expr *ast.TypeAssertExpr) {
	p.printNode(expr.X)
	p.print(".(")
	if expr.Type != nil {
		p.printNode(expr.Type)
	} else {
		p.print("type")
	}
	p.print(")")
}

func (p *printer) TypeSpec(spec *ast.TypeSpec) {
	if spec.Doc != nil {
		p.CommentGroup(spec.Doc)
		p.writeIndent()
	}
	// not need tabs - tabs in block
	p.Ident(spec.Name)
	p.print(" ")
	p.printNode(spec.Type)
	if spec.Comment != nil {
		p.CommentGroup(spec.Comment)
	}
}

func (p *printer) TypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
	p.print("switch ")
	if stmt.Init != nil {
		p.printNode(stmt.Init)
		p.print("; ")
	}
	p.printNode(stmt.Assign)
	p.BlockStmt(stmt.Body)
}

func (p *printer) UnaryExpr(expr *ast.UnaryExpr) {
	p.print(expr.Op.String())
	p.printNode(expr.X)
}

func (p *printer) ValueSpec(spec *ast.ValueSpec) {
	if spec.Doc != nil {
		p.CommentGroup(spec.Doc)
		p.writeIndent()
	}
	// not need tabs - tabs in block
	p.Names(spec.Names)
	if spec.Type != nil {
		p.print(" ")
		p.printNode(spec.Type)
	}
	if spec.Values != nil {
		p.print(" = ")
		p.Exprs(spec.Values)
	}
	if spec.Comment != nil {
		p.CommentGroup(spec.Comment)
	}
}
