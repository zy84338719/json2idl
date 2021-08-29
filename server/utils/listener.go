package listener

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"gitlab.com/jsonToIDL/utils/parser"
	"strings"
)

type Node struct {
	Type      string
	Value     string
	ValueType string
	KeyIsNum  bool
}

type Listener struct {
	*parser.BaseJSONListener
	gocodeMap  map[antlr.Tree]Node `json:"gocode_map"`
	JsonStr    string
	Target     Target
	SubStructs []string
}

func NewJsonToGoListener(t Target) *Listener {
	return &Listener{
		BaseJSONListener: &parser.BaseJSONListener{},
		gocodeMap:        make(map[antlr.Tree]Node),
		JsonStr:          "",
		Target:           t,
	}
}

func (l *Listener) PrintGocodeMap() {
	for k, v := range l.gocodeMap {
		if false {
			fmt.Println(fmt.Sprintf("%T", k), v)
		}
	}
}

// VisitTerminal is called when a terminal node is visited.
func (l *Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (l *Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (l *Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (l *Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterJson is called when production json is entered.
func (l *Listener) EnterJson(ctx *parser.JsonContext) {}

// ExitJson is called when production json is exited.
func (l *Listener) ExitJson(ctx *parser.JsonContext) {
	l.JsonStr = l.Target.ExitJson(l.gocodeMap[ctx.Value()].Type, l.gocodeMap[ctx.Value()].Value)
}

// EnterObj is called when production obj is entered.
func (l *Listener) EnterObj(ctx *parser.ObjContext) {}

func (l *Listener) isMap(ctx *parser.ObjContext) bool {
	for _, p := range ctx.AllPair() {
		if l.gocodeMap[p].KeyIsNum {
			return true
		}
	}
	return false
}

// ExitObj is called when production obj is exited.
func (l *Listener) ExitObj(ctx *parser.ObjContext) {
	sb := strings.Builder{}
	sb.WriteString(l.Target.PreExitObj("", ""))
	bIsMap := l.isMap(ctx)
	for i, p := range ctx.AllPair() {
		switch l.Target.(type) {
		case *IdlTarget:
			sb.WriteString(fmt.Sprintf("%d: ", i+1) + l.Target.ExitObj(l.gocodeMap[p].Type, l.gocodeMap[p].Value, i == len(ctx.AllPair())-1, bIsMap))
		default:
			sb.WriteString(l.Target.ExitObj(l.gocodeMap[p].Type, l.gocodeMap[p].Value, i == len(ctx.AllPair())-1, bIsMap))
		}

	}
	sb.WriteString(l.Target.PostExitObj("", ""))
	l.gocodeMap[ctx] = Node{
		Type:      "struct",
		Value:     sb.String(),
		ValueType: "pair",
	}
}

// EnterPair is called when production pair is entered.
func (l *Listener) EnterPair(ctx *parser.PairContext) {}

// ExitPair is called when production pair is exited.
func (l *Listener) ExitPair(ctx *parser.PairContext) {
	subStruct, pair := l.Target.ExitPair(0, ctx.STRING().GetText(), l.gocodeMap[ctx.Value()].Type,
		l.gocodeMap[ctx.Value()].Value, l.gocodeMap[ctx.Value()].ValueType)
	l.gocodeMap[ctx] = Node{
		Type:      "KV",
		Value:     pair,
		ValueType: l.gocodeMap[ctx.Value()].Type,
		KeyIsNum:  IsNumber(ctx.STRING().GetText()),
	}

	//if l.gocodeMap[ctx.Value()].Type=="struct"{
	l.SubStructs = append(l.SubStructs, subStruct)
	//}
	//fmt.Println(ctx.Value().GetText(),"======>",l.gocodeMap[ctx.Value()])
}

// EnterArr is called when production arr is entered.
func (l *Listener) EnterArr(ctx *parser.ArrContext) {}

// ExitArr is called when production arr is exited.
func (l *Listener) ExitArr(ctx *parser.ArrContext) {
	ValueType := ""
	if l.gocodeMap[ctx.Value(0)].Type == "array" {
		ValueType = "list<" + l.gocodeMap[ctx.Value(0)].ValueType + ">"
	} else {
		ValueType = l.gocodeMap[ctx.Value(0)].Type
	}
	l.gocodeMap[ctx] = Node{
		Type:      "array",
		Value:     l.Target.ExitArr(l.gocodeMap[ctx.Value(0)].Type, l.gocodeMap[ctx.Value(0)].Value),
		ValueType: ValueType,
	}
	//fmt.Println(ctx.GetChild(0),ctx.Value(0))
}

// EnterValue is called when production value is entered.
func (l *Listener) EnterValue(ctx *parser.ValueContext) {}

// ExitValue is called when production value is exited.
func (l *Listener) ExitValue(ctx *parser.ValueContext) {
	if ctx.Arr() != nil {
		l.gocodeMap[ctx] = l.gocodeMap[ctx.Arr()]
	} else if ctx.Obj() != nil {
		l.gocodeMap[ctx] = l.gocodeMap[ctx.Obj()]
	} else if ctx.NUMBER() != nil {
		l.gocodeMap[ctx] = Node{
			Type:      "float64",
			Value:     l.Target.ExitValue("float64", ctx.NUMBER().GetText()),
			ValueType: "float64",
		}
		//l.gocodeMap[ctx.NUMBER()]
	} else if ctx.STRING() != nil {
		l.gocodeMap[ctx] = Node{
			Type:      "string",
			Value:     l.Target.ExitValue("string", ctx.STRING().GetText()),
			ValueType: "string",
		}
		//l.gocodeMap[ctx.STRING()]
	} else {
		if ctx.GetText() == "true" || ctx.GetText() == "false" {
			l.gocodeMap[ctx] = Node{
				Type:      "bool",
				Value:     l.Target.ExitValue("bool", ctx.GetText()),
				ValueType: ctx.GetText(),
			}
		} else if ctx.GetText() == "null" {
			l.gocodeMap[ctx] = Node{
				Type:      "null",
				Value:     l.Target.ExitValue("null", ctx.GetText()),
				ValueType: "struct",
			}
		} else {
			l.gocodeMap[ctx] = Node{
				Type:      "string",
				Value:     l.Target.ExitValue("string", ctx.GetText()),
				ValueType: "string",
			}
		}
	}
}
