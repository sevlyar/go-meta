package meta

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

var DirectiveMark = "@"

type Rule interface {
	ApplyTo(node ast.Node) (ast.Node, error)
}

var registered = make(map[string]Rule)

func Register(rule Rule) {
	t := reflect.TypeOf(rule).Elem()
	name := t.Name()

	if rm, ok := registered[name]; ok {
		rt := reflect.TypeOf(rm)
		if rt == t {
			panic(fmt.Sprintf("meta-rule %s registered twice", t))
		}
		panic(fmt.Sprintf("meta-rule %s ovirride %s", rt, t))
	}

	registered[t.Name()] = rule
}

func FindMeta(name string) (rule Rule, err error) {
	ok := false
	if rule, ok = registered[name]; ok {
		return
	}

	for key, val := range registered {
		fmt.Printf("'%s':%+v\n", key, val)
	}
	err = fmt.Errorf("meta-rule %s not found", name)
	return
}

// args: json-encoded fields of rule
func SetArgs(rule Rule, args string) error {
	return json.Unmarshal([]byte(args), rule)
}

func Extract(node ast.Node) []*Directive {
	list := make([]*Directive, 0, 512)

	ast.Walk(&dirExtr{nil, node, &list}, node)

	return list
}

func Apply(ds []*Directive) (err error) {
	for _, d := range ds {
		var rule Rule
		if rule, err = FindMeta(d.RuleName); err != nil {
			// TODO: extend error information
			break
		}

		if err = SetArgs(rule, d.RuleArgs); err != nil {
			// TODO: extend error information
			break
		}

		//var node ast.Node
		if _, err = rule.ApplyTo(d.Node); err != nil {
			// TODO: extend error information
			break
		}
	}

	return
}

type dirExtr struct {
	preparent ast.Node
	parent    ast.Node
	list      *[]*Directive
}

func (extr *dirExtr) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.CommentGroup:
		ds := ExtractDirectives(x.Text(), extr.parent, extr.preparent)
		if ds != nil {
			extr.appendDirectives(ds)
		}
	}
	return &dirExtr{extr.parent, node, extr.list}
}

func (extr *dirExtr) appendDirectives(ds []*Directive) {
	*extr.list = append(*extr.list, ds...)
}

type Directive struct {
	RuleName string
	RuleArgs string
	Node     ast.Node
	Parent   ast.Node
}

func ExtractDirectives(comment string, node, parent ast.Node) (ds []*Directive) {
	var name string
	for {
		name, comment = extractDirectiveName(comment)
		if len(name) == 0 {
			break
		}
		d := &Directive{RuleName: name, Node: node, Parent: parent}
		d.RuleArgs, comment = extractDirectiveArgs(comment)

		if ds == nil {
			ds = make([]*Directive, 0, 2)
		}
		ds = append(ds, d)
	}

	return
}

func extractDirectiveName(comment string) (name, tail string) {
	apos := strings.Index(comment, DirectiveMark)
	if apos < 0 {
		return
	}
	name = comment[apos+1:]

	aend := strings.IndexAny(name, " {\r\n")
	if aend > -1 {
		tail = name[aend:]
		name = name[:aend]
	}

	return
}

func extractDirectiveArgs(comment string) (args, tail string) {
	pos := strings.Index(comment, "{") // TODO: append json array support []
	if pos < 0 {
		tail = comment
		return
	}
	args = comment[pos:]

	openBrcs := 1
	for pos := 1; pos < len(args); pos++ {
		switch args[pos] {
		case '{':
			openBrcs++
		case '}':
			openBrcs--
		default:
		}

		if openBrcs == 0 {
			tail = args[pos+1:]
			args = args[:pos+1]
			break
		}
	}

	return
}
