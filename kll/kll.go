package kll

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Value interface {
	On_sum(Value) Value
	On_sub(Value) Value
	On_mul(Value) Value
	On_div(Value) Value
	Re_string(prefix string) string
	Re_number() float64
	Re_bool() bool
	On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any)
	On_get_attr(name string) Value
	On_set_attr(name string, value Value) Value
	On_in(name Value) Value
	//On_get_variable(name string) *Value
	VType() string
}
type Number struct {
	Value float64 `json:"value"`
	VTp   string  `json:"value type"`
}
type Node struct {
	Tp    string `json:"type"`
	VTp   string `json:"value type"`
	Line  uint64
	Col   uint64
	Value []Value `json:"value"`
}
type String struct {
	Value string `json:"value"`
	VTp   string `json:"value type"`
}
type Null struct {
	VTp string `json:"value type"`
}
type Bool struct {
	VTp   string `json:"value type"`
	Value bool
}
type Function struct {
	VTp    string `json:"value type"`
	nodes  []Value
	inter  *Interpreter
	locals *Object
	args   []Variable
}
type GoFunction struct {
	VTp      string `json:"value type"`
	function func(args []Value, kwargs map[string]*Variable, pos int) (Value, any)
}
type Pointer struct {
	VTp   string `json:"value type"`
	value *Variable
}
type Object struct {
	VTp   string `json:"value type"`
	value map[string]Variable
}
type Array struct {
	VTp   string `json:"value type"`
	Value []Value
}

var lang = "pt-br"

func (this Number) On_sum(value Value) Value {
	return Create_Number(this.Re_number() + value.Re_number())
}
func (this Number) On_sub(value Value) Value {
	return Create_Number(this.Re_number() - value.Re_number())
}
func (this Number) On_div(value Value) Value {
	return Create_Number(this.Re_number() / value.Re_number())
}
func (this Number) On_mul(value Value) Value {
	return Create_Number(this.Re_number() * value.Re_number())
}
func (this Number) Re_string(prefix string) string {
	re1, re2 := math.Modf(this.Value)
	if (re2) == 0 {
		return fmt.Sprint(int(re1))
	}
	return fmt.Sprint(this.Value)
}
func (this Number) Re_number() float64 {
	return this.Value
}
func (this Number) Re_bool() bool {
	return this.Value > 0
}
func (this Number) VType() string {
	return "Number"
}
func (this Number) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Number) On_get_attr(name string) Value {
	switch name {
	case "string":
		return Create_String(this.Re_string(""))
	case "is_int":
		_, re := math.Modf(this.Value)
		return Create_Bool(re == 0.0)
	case "length":
		return Create_Number(float64(len(strings.Replace(this.Re_string(""), ".", "", 1))))
	case "length1":
		v, _ := math.Modf(this.Value)
		return Create_Number(float64(len(fmt.Sprint(v))))
	case "length2":
		_, v := math.Modf(this.Value)
		return Create_Number(float64(len(fmt.Sprint(v))))
	case "decimal":
		_, v := math.Modf(this.Value)
		return Create_Number(v)
	}
	return Create_Null()
}
func (this Number) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Number) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this Node) On_sum(value Value) Value {
	return Create_Null()
}
func (this Node) On_sub(value Value) Value {
	return Create_Null()
}
func (this Node) On_div(value Value) Value {
	return Create_Null()
}
func (this Node) On_mul(value Value) Value {
	return Create_Null()
}
func (this Node) Re_string(prefix string) string {
	str := "Node " + this.Tp
	var n uint64 = 0
	for n < uint64(len(this.Value)) {
		str += "\n"
		v := this.Value[n].Re_string(" ")
		var lines []string
		//str += "├─"
		//│
		for i, l := range strings.Split(v, "\n") {
			if i <= 1 {
				lines = append(lines, l)
			} else {
				lines = append(lines, l)
			}
		}
		for i, l := range lines {
			if i > 0 {
				str += "\n"
			}
			str += prefix + l
		}
		n++
	}
	return str
}
func (this Node) Re_number() float64 {
	return 0
}
func (this Node) Re_bool() bool {
	return false
}
func (this Node) VType() string {
	return "Node"
}
func (this Node) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Node) On_get_attr(name string) Value {
	return Create_Null()
}
func (this Node) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Node) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this String) On_sum(value Value) Value {
	return Create_String(this.Re_string("") + value.Re_string(""))
}
func (this String) On_sub(value Value) Value {
	return nil
}
func (this String) On_div(value Value) Value {
	return nil
}
func (this String) On_mul(value Value) Value {
	return nil
}
func (this String) Re_string(prefix string) string {
	return this.Value
}
func (this String) Re_number() float64 {
	return 0
}
func (this String) Re_bool() bool {
	return len(this.Value) > 0
}
func (this String) VType() string {
	return "String"
}
func (this String) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this String) On_get_attr(name string) Value {
	switch name {
	case "number":
		return To_int(this.Re_string(""))
	case "length":
		return Create_Number(float64(len(this.Re_string(""))))
	case "replace":
		return Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			old := ""
			if len(args) >= 1 {
				old = args[0].Re_string("")
			} else if _, ok := kwargs["old"]; ok {
				old = kwargs["old"].Value.Re_string("")
			}
			new := ""
			if len(args) >= 2 {
				new = args[1].Re_string("")
			} else if _, ok := kwargs["new"]; ok {
				new = kwargs["new"].Value.Re_string("")
			}
			limit := 0
			if len(args) >= 3 {
				limit = int(args[2].Re_number())
			} else if _, ok := kwargs["limit"]; ok {
				limit = int(kwargs["limit"].Value.Re_number())
			}
			return Create_String(strings.Replace(this.Re_string(""), old, new, limit)), nil
		})
	case "startswith":
		return Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			value := ""
			if len(args) >= 1 {
				value = args[0].Re_string("")
			} else if _, ok := kwargs["value"]; ok {
				value = kwargs["value"].Value.Re_string("")
			}
			return Create_Bool(strings.HasPrefix(this.Re_string(""), value)), nil
		})
	case "endswith":
		return Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			value := ""
			if len(args) >= 1 {
				value = args[0].Re_string("")
			} else if _, ok := kwargs["value"]; ok {
				value = kwargs["value"].Value.Re_string("")
			}
			return Create_Bool(strings.HasSuffix(this.Re_string(""), value)), nil
		})
	}
	return Create_Null()
}
func (this String) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this String) On_in(name Value) Value {
	return Create_Bool(strings.Index(this.Value, name.Re_string("")) != -1)
}

func (this Null) On_sum(value Value) Value {
	return Create_Null()
}
func (this Null) On_sub(value Value) Value {
	return Create_Null()
}
func (this Null) On_div(value Value) Value {
	return Create_Null()
}
func (this Null) On_mul(value Value) Value {
	return Create_Null()
}
func (this Null) Re_string(prefix string) string {
	return "null"
}
func (this Null) Re_number() float64 {
	return -1
}
func (this Null) Re_bool() bool {
	return false
}
func (this Null) VType() string {
	return "Null"
}
func (this Null) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Null) On_get_attr(name string) Value {
	switch name {
	}
	return Create_Null()
}
func (this Null) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Null) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this Function) On_sum(value Value) Value {
	return Create_Null()
}
func (this Function) On_sub(value Value) Value {
	return Create_Null()
}
func (this Function) On_div(value Value) Value {
	return Create_Null()
}
func (this Function) On_mul(value Value) Value {
	return Create_Null()
}
func (this Function) Re_string(prefix string) string {
	str := "<function nodes:" + Create_Array(this.nodes).Re_string("") + ", args:"
	obj := make(map[string]Value)
	for _, a := range this.args {
		obj[a.name] = a.Value
	}
	str += Create_Object(obj).Re_string("")
	return str + ">"
}
func (this Function) Re_number() float64 {
	return -1
}
func (this Function) Re_bool() bool {
	return false
}
func (this Function) VType() string {
	return "Function"
}
func (this Function) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	for i := range this.args {
		if i > len(this.args) {
			break
		}
		v := this.args[i].Value
		if i < len(args) {
			v = args[i]
		}
		v2, ok := kwargs[this.args[i].name]
		if ok {
			v = v2.Value
		}
		this.locals.Create_Var(this.args[i].name, pos+1, v, false)
	}
	re := Create_Null()
	var err any = nil
	i := uint64(0)
	for i < uint64(len(this.nodes)) {
		if this.nodes[i].(Node).Tp == "return" {
			re, err = this.inter.exec_node(this.nodes[i].(Node).Value[0], this.locals, pos+1)
			break
		}
		_, e := this.inter.exec_node(this.nodes[i], this.locals, pos+1)
		if e != nil {
			err = e
		}
		i++
	}
	for v := range this.locals.value {
		if this.locals.value[v].Pos >= pos+1 {
			delete(this.locals.value, v)
		}
	}
	return re, err
}
func (this Function) On_get_attr(name string) Value {
	switch name {
	}
	return Create_Null()
}
func (this Function) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Function) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this GoFunction) On_sum(value Value) Value {
	return Create_Null()
}
func (this GoFunction) On_sub(value Value) Value {
	return Create_Null()
}
func (this GoFunction) On_div(value Value) Value {
	return Create_Null()
}
func (this GoFunction) On_mul(value Value) Value {
	return Create_Null()
}
func (this GoFunction) Re_string(prefix string) string {
	str := "<go_function>"
	return str
}
func (this GoFunction) Re_number() float64 {
	return -1
}
func (this GoFunction) Re_bool() bool {
	return false
}
func (this GoFunction) VType() string {
	return "GoFunction"
}
func (this GoFunction) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return this.function(args, kwargs, pos)
}
func (this GoFunction) On_get_attr(name string) Value {
	switch name {
	}
	return Create_Null()
}
func (this GoFunction) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this GoFunction) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this Bool) On_sum(value Value) Value {
	return Create_Null()
}
func (this Bool) On_sub(value Value) Value {
	return Create_Null()
}
func (this Bool) On_div(value Value) Value {
	return Create_Null()
}
func (this Bool) On_mul(value Value) Value {
	return Create_Null()
}
func (this Bool) Re_string(prefix string) string {
	if this.Value {
		return "true"
	}
	return "false"
}
func (this Bool) Re_number() float64 {
	return -1
}
func (this Bool) Re_bool() bool {
	return this.Value
}
func (this Bool) VType() string {
	return "Bool"
}
func (this Bool) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Bool) On_get_attr(name string) Value {
	switch name {
	}
	return Create_Null()
}
func (this Bool) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Bool) On_in(name Value) Value {
	return Create_Bool(false)
}

func (this Object) On_sum(value Value) Value {
	return Create_Null()
}
func (this Object) On_sub(value Value) Value {
	return Create_Null()
}
func (this Object) On_div(value Value) Value {
	return Create_Null()
}
func (this Object) On_mul(value Value) Value {
	return Create_Null()
}
func (this Object) Re_string(prefix string) string {
	re := "{"
	pos := 0
	for i := range this.value {
		if pos > 0 {
			re += ", "
		}
		re += i + ":"
		switch this.value[i].Value.VType() {
		case "String":
			re += string('"') + this.value[i].Value.Re_string("") + string('"')
		default:
			re += this.value[i].Value.Re_string("")
		}
		pos++
	}
	return re + "}"
}
func (this Object) Re_number() float64 {
	return -1
}
func (this Object) Re_bool() bool {
	return false
}
func (this Object) VType() string {
	return "Object"
}
func (this Object) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Object) On_get_attr(name string) Value {
	switch name {
	default:
		return this.value[name].Value
	}
	//return Create_Null()
}
func (this Object) On_set_attr(name string, value Value) Value {
	if e, ok := this.value[name]; ok {
		e.Value = value
		this.value[name] = e
	}
	return value
}
func (this Object) On_get_Variable(name string) (*Variable, bool) {
	var v Variable
	var ok bool
	v, ok = this.value[name]
	return &v, ok
}
func (this Object) Create_Var(name string, pos int, value Value, is_const bool) Variable {
	this.value[name] = Variable{name: name, Pos: pos, is_const: is_const, Value: value}
	return this.value[name]
}
func (this Object) On_in(name Value) Value {
	if name.VType() == "Node" {
		if name.(Node).Tp == "var" {
			_, ok := this.value[name.(Node).Value[0].Re_string("")]
			return Create_Bool(ok)
		}
	}
	return Create_Bool(false)
}

func (this Array) On_sum(value Value) Value {
	return Create_Null()
}
func (this Array) On_sub(value Value) Value {
	return Create_Null()
}
func (this Array) On_div(value Value) Value {
	return Create_Null()
}
func (this Array) On_mul(value Value) Value {
	return Create_Null()
}
func (this Array) Re_string(prefix string) string {
	str := "["
	for i, e := range this.Value {
		if i > 0 {
			str += ", "
		}
		if e.VType() == "String" {
			str += string('"') + e.Re_string("") + string('"')
		} else {
			str += e.Re_string("")
		}
	}
	return str + "]"
}
func (this Array) Re_number() float64 {
	return -1
}
func (this Array) Re_bool() bool {
	return false
}
func (this Array) VType() string {
	return "Array"
}
func (this Array) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return Create_Null(), nil
}
func (this Array) On_get_attr(name string) Value {
	switch name {
	}
	return Create_Null()
}
func (this Array) On_set_attr(name string, value Value) Value {
	return Create_Null()
}
func (this Array) On_in(name Value) Value {
	for _, v := range this.Value {
		if v == name {
			return Create_Bool(true)
		}
	}
	return Create_Bool(false)
}

func (this Pointer) On_sum(value Value) Value {
	return this.value.Value.On_sum(value)
}
func (this Pointer) On_sub(value Value) Value {
	return this.value.Value.On_sub(value)
}
func (this Pointer) On_div(value Value) Value {
	return this.value.Value.On_div(value)
}
func (this Pointer) On_mul(value Value) Value {
	return this.value.Value.On_mul(value)
}
func (this Pointer) Re_string(prefix string) string {
	return this.value.Value.Re_string("")
}
func (this Pointer) Re_number() float64 {
	return this.value.Value.Re_number()
}
func (this Pointer) Re_bool() bool {
	return false
}
func (this Pointer) VType() string {
	return "Pointer"
}
func (this Pointer) On_call(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return this.value.Value.On_call(args, kwargs, pos)
}
func (this Pointer) On_get_attr(name string) Value {
	return this.value.Value.On_get_attr(name)
}
func (this Pointer) On_set_attr(name string, value Value) Value {
	return this.value.Value.On_set_attr(name, value)
}
func (this Pointer) On_in(name Value) Value {
	return this.value.Value.On_in(name)
}

func Sum(value1 Value, value2 Value) Value {
	return value1.On_sum(value2)
}
func Sub(value1 Value, value2 Value) Value {
	return value1.On_sub(value2)
}
func Mul(value1 Value, value2 Value) Value {
	return value1.On_mul(value2)
}
func Div(value1 Value, value2 Value) Value {
	return value1.On_div(value2)
}

func Call(value Value, args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
	return value.On_call(args, kwargs, pos)
}
func Get_attr(value Value, name Value) Value {
	if name.VType() == "Node" {
		if name.(Node).Tp == "get attr" {
			return Get_attr(value.On_get_attr(name.(Node).Value[0].(Node).Value[0].Re_string("")), name.(Node).Value[1].(Node))
		} else if name.(Node).Tp == "var" {
			return value.On_get_attr(name.(Node).Value[0].Re_string(""))
		}
	} else if name.VType() == "String" {
		return value.On_get_attr(name.Re_string(""))
	}
	return Create_Null()
}
func Set_attr(value Value, name Value, value_set Value) Value {
	if name.VType() == "Node" {
		if name.(Node).Tp == "get attr" {
			value.On_set_attr(name.(Node).Value[0].(Node).Value[0].Re_string(""), Set_attr(value.On_get_attr(name.(Node).Value[0].(Node).Value[0].Re_string("")), name.(Node).Value[1].(Node), value_set))
		} else if name.(Node).Tp == "var" {
			value.On_set_attr(name.(Node).Value[0].Re_string(""), value_set)
		}
	} else {
		value.On_set_attr(name.Re_string(""), value_set)
	}
	return value
}
func In(value Value, name Value) Value {
	return value.On_in(name)
}

/*
	func Get_variable(value Value, name Value) *Value {
		if name.VType() == "Node" {
			if name.(Node).Tp == "get attr" {
				return Get_variable(value.On_get_attr(name.(Node).Value[0].(Node).Value[0].Re_string()), name.(Node).Value[1].(Node))
			} else if name.(Node).Tp == "var" {
				return value.On_get_variable(name.(Node).Value[0].Re_string())
			}
		} else if name.VType() == "String" {
			return value.On_get_variable(name.Re_string())
		}
		return &Create_Null()
	}
*/
func Create_Number(value float64) Value {
	re := Number{Value: value}
	re.VTp = re.VType()
	return re
}
func Create_String(value string) Value {
	re := String{Value: value}
	re.VTp = re.VType()
	return re
}
func Create_Null() Value {
	re := Null{}
	re.VTp = re.VType()
	return re
}
func Create_Bool(value bool) Value {
	re := Bool{Value: value}
	re.VTp = re.VType()
	return re
}
func Create_Node(value []Value, tp string, line uint64, col uint64) Value {
	re := Node{Value: value, Tp: tp, Line: line, Col: col}
	re.VTp = re.VType()
	return re
}
func Create_Function(nodes []Value, inter *Interpreter, locals *Object, args []Variable) Value {

	re := Function{nodes: nodes, inter: inter, locals: locals, args: args}
	re.VTp = re.VType()
	return re
}
func Create_GoFunction(function func(args []Value, kwargs map[string]*Variable, pos int) (Value, any)) Value {
	re := GoFunction{function: function}
	re.VTp = re.VType()
	return re
}
func Create_Pointer(value *Variable) Value {
	re := Pointer{value: value}
	re.VTp = re.VType()
	return re
}
func Create_Object(publics map[string]Value) Value {
	vars := make(map[string]Variable)
	for i := range publics {
		vars[i] = Variable{name: i, Value: publics[i]}
	}
	re := Object{value: vars}
	re.VTp = re.VType()
	return re
}
func Create_Array(values []Value) Value {
	re := Array{Value: values}
	re.VTp = re.VType()
	return re
}

type Token struct {
	value Value
	tp    string
	line  uint64
	col   uint64
}

func (this *Token) Re_string() string {
	if this.value == nil {
		return this.tp
	}
	return this.tp + ":" + this.value.Re_string("")
}

type Error struct {
	msg         string
	other_error any
	line        uint64
	col         uint64
	lines       []string
}

//export lexer
type Lexer struct {
	tok  int
	txt  string
	char string
	line uint64
	col  uint64
}
type Cache struct {
	Txt   string  `json:"txt"`
	Nodes []Value `json:"nodes"`
}

func lang_text(txt string, extras []string) string {
	switch lang {
	case "pt-br":
		switch txt {
		case "erro1":
			return "Erro de Syntaxe: "
		case "erro2":
			return "Erro de Variavel: "
		case "erro msg1":
			return "no numero possui mais de 1 ponto final"
		case "erro msg2":
			return "o simbolo " + extras[0] + " não existe nessa linguagem"
		case "erro msg3":
			return "vc colocol em uma posição errada"
		case "erro msg4":
			return "vc esqueceu de fechar uma string"
		case "erro msg5":
			return "A variavel '" + extras[0] + "' não existe"
		case "erro msg6":
			return "expresão invalida"
		case "erro msg7":
			return "vc esqueceu de fechar as chaves"
		case "erro msg8":
			return "vc esqueceu de fechar os parentses"
		}
	}
	return ""
}

var varsName = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

func (this *Lexer) is_next(txt string) bool {
	ok := false
	walk := 1
	spaces := 0
	v := this.char
	(*Lexer).next(this)
	for this.char != "" {
		if walk-spaces >= len(this.char)+1 {
			break
		}
		if this.char == " " || this.char == "\t" {
			walk++
			spaces++
			(*Lexer).next(this)
			continue
		}
		v += this.char
		if v == txt {
			ok = true
			break
		}
		(*Lexer).next(this)
		walk++
	}
	if !ok {
		this.Return(walk)
	}
	return ok
}
func (this *Lexer) Return(number int) {
	i := 0
	for i < number {
		this.tok -= 1
		this.col -= 1
		(*Lexer).load(this)
		if this.char == "\n" {
			this.line--
			this.col = uint64(len(strings.Split(this.txt, "\n")[this.line-1]))
		}
		i++
	}
}
func To_int(n string) Value {
	f, _ := strconv.ParseFloat(n, 64)
	return Create_Number(f)
}
func (this *Lexer) Tokenizer(txt string) ([]Token, any) {
	this.tok = -1
	this.txt = txt
	this.char = ""
	this.line = 1
	this.col = 0
	(*Lexer).next(this)
	var re []Token
	for this.char != "" {
		if strings.Contains("0123456789.", this.char) {
			var n string = this.char
			ok := false
			if this.char == "." {
				ok = true
			}
			(*Lexer).next(this)
			for this.char != "" && strings.Contains("0123456789.", this.char) {
				if this.char == "." {
					if ok {
						return []Token{}, Error{msg: lang_text("erro1", []string{}) + lang_text("erro msg1", []string{}), lines: strings.Split(this.txt, "\n"), line: this.line, col: this.col}
					}
					ok = true
				}
				n += this.char
				(*Lexer).next(this)
			}
			if n == "." || strings.HasSuffix(n, ".") {
				if strings.HasSuffix(n, ".") && !(n == ".") {
					re = append(re, Token{value: To_int(n), tp: "value", col: this.col, line: this.line})
				}
				re = append(re, Token{tp: ".", col: this.col, line: this.line})
				continue
			}
			re = append(re, Token{value: To_int(n), tp: "value", col: this.col, line: this.line})
			continue
		} else if strings.Contains(varsName, this.char) {
			var n string = this.char
			col, line := this.col, this.line
			(*Lexer).next(this)
			for this.char != "" && strings.Contains(varsName+"0123456789", this.char) {
				n += this.char
				(*Lexer).next(this)
			}
			switch n {
			case "global":
				re = append(re, Token{tp: "create global", col: col, line: line})
				break
			case "return":
				re = append(re, Token{tp: "return", col: col, line: line})
				break
			case "local":
				re = append(re, Token{tp: "create local", col: col, line: line})
				break
			case "and":
				re = append(re, Token{tp: "&&", col: this.col, line: this.line})
				break
			case "or":
				re = append(re, Token{tp: "||", col: this.col, line: this.line})
				break
			case "var":
				re = append(re, Token{tp: "create local", col: col, line: line})
			case "function":
				re = append(re, Token{tp: "function", col: col, line: line})
				break
			case "if":
				re = append(re, Token{tp: "if", col: col, line: line})
				break
			case "exist":
				re = append(re, Token{tp: "exist", col: col, line: line})
				break
			default:
				re = append(re, Token{tp: "var", value: Create_String(n), col: col, line: line})
				break
			}
			continue
		} else {
			switch this.char {
			case "+":
				re = append(re, Token{tp: "+", col: this.col - 1, line: this.line})
				break
			case "-":
				re = append(re, Token{tp: "-", col: this.col - 1, line: this.line})
				break
			case "*":
				re = append(re, Token{tp: "*", col: this.col - 1, line: this.line})
				break
			case "/":
				re = append(re, Token{tp: "/", col: this.col - 1, line: this.line})
				break
			case "&":
				if this.is_next("&&") {
					re = append(re, Token{tp: "&&", col: this.col, line: this.line})
				} else {
					return []Token{}, Error{msg: lang_text("erro1", []string{}) + lang_text("erro msg2", []string{this.char}), line: this.line, col: this.col, lines: strings.Split(this.txt, "\n")}
				}
				break
			case "|":
				if this.is_next("||") {
					re = append(re, Token{tp: "||", col: this.col, line: this.line})
				} else {
					return []Token{}, Error{msg: lang_text("erro1", []string{}) + lang_text("erro msg2", []string{this.char}), line: this.line, col: this.col, lines: strings.Split(this.txt, "\n")}
				}
			case "=":
				if this.is_next("==") {
					re = append(re, Token{tp: "==", col: this.col - 1, line: this.line})
				} else {
					re = append(re, Token{tp: "=", col: this.col - 1, line: this.line})
				}
				break
			case "(":
				re = append(re, Token{tp: "(", col: this.col - 1, line: this.line})
				break
			case ")":
				re = append(re, Token{tp: ")", col: this.col - 1, line: this.line})
				break
			case "{":
				re = append(re, Token{tp: "{", col: this.col - 1, line: this.line})
				break
			case "}":
				re = append(re, Token{tp: "}", col: this.col - 1, line: this.line})
				break
			case ",":
				re = append(re, Token{tp: ",", col: this.col - 1, line: this.line})
				break
			case string('"'):
				col := this.col
				line := this.line
				(*Lexer).next(this)
				v := this.char
				(*Lexer).next(this)
				ok := true
				for this.char != "" {
					if this.char == string('"') {
						ok = false
						break
					}
					v += this.char
					(*Lexer).next(this)
				}
				if ok {
					return []Token{}, Error{msg: lang_text("erro1", nil) + lang_text("erro msg4", nil), line: line, col: col, lines: strings.Split(this.txt, "\n")}
				}
				re = append(re, Token{tp: "value", value: Create_String(v), col: this.col, line: this.line})
				break
			case ";":
				re = append(re, Token{tp: "split", col: this.col, line: this.line})
				break
			case " ":
				break
			case "\t":
				break
			case "\n":
				re = append(re, Token{tp: "new line", col: this.col, line: this.line})
				break
			default:
				return []Token{}, Error{msg: lang_text("erro1", []string{}) + lang_text("erro msg2", []string{this.char}), line: this.line, col: this.col, lines: strings.Split(this.txt, "\n")}
			}
		}
		(*Lexer).next(this)
	}
	return re, nil
}
func (this *Lexer) next() {
	this.tok++
	(*Lexer).load(this)
}
func (this *Lexer) load() {
	if this.tok >= len(this.txt) {
		this.char = ""
	} else {
		this.char = string([]rune(this.txt)[this.tok])
		this.col++
		if this.char == "\n" {
			this.line++
			this.col = 1
		}
	}
}

func splitTokens(toks []Token, st string, vspt string, vept string) ([][]Token, int) {
	var ts []Token = toks
	i := 0
	v := 1
	walk := 0
	var re [][]Token
	for true {
		if i >= len(ts) {
			re = append(re, ts)
			break
		} else if ts[i].tp == vspt {
			v++
		} else if ts[i].tp == vept {
			v--
		}
		if v < 1 {
			re = append(re, ts[:i])
			break
		}
		if ts[i].tp == st {
			re = append(re, ts[:i])
			ts = ts[i+1:]
			i = -1
		}
		i++
		walk++
	}
	return re, walk
}

//extern Parser
type Parser struct {
	codes   [][]Token
	code    uint64
	tok_pos uint64
	tok     Token
	Lexer   Lexer
	txt     string
}

func (this *Parser) next_tok() {
	this.tok_pos++
	this.load_tok()
}
func (this *Parser) load_tok() {
	if this.code >= uint64(len(this.codes)) {
		return
	}
	if this.tok_pos >= uint64(len(this.codes[this.code])) {
		var line uint64 = 1
		var col uint64 = 1
		if len(this.codes[this.code]) >= 1 {
			line = this.codes[this.code][len(this.codes[this.code])-1].line
			col = this.codes[this.code][len(this.codes[this.code])-1].col
		}
		this.tok = Token{tp: "end code", line: line, col: col}
	} else {
		this.tok = this.codes[this.code][this.tok_pos]
	}
}
func is_end_code(tok Token) bool {
	return tok.tp == "end code"
}
func (this *Parser) load_code() {
	if this.code < uint64(len(this.codes)) {
		this.tok_pos = 0
		(*Parser).load_tok(this)
	} else {
		this.tok = Token{tp: "end code", line: 1, col: 1}
	}
}
func (this *Parser) next_code() {
	this.code++
	this.load_code()
}
func (this *Parser) Make_Nodes_Toks(toks [][]Token, end string, ok bool) ([]Value, any) {
	this.code = 0
	var re []Value
	this.codes = toks
	this.load_code()
	var result Value
	var err any
	for this.code < uint64(len(this.codes)) {
		if ok {
			result, err = (*Parser).new_line(this)
		} else {
			result, err = (*Parser).expr(this)
		}
		if err != nil {
			return []Value{}, err
		}
		if result == nil {
			(*Parser).next_code(this)
			continue
		}
		if result.(Node).Tp != "null" {
			re = append(re, result)
		}
		(*Parser).next_code(this)
	}
	return re, nil
}
func (this *Parser) Make_Nodes(txt string) ([]Value, any) {
	this.txt = txt
	tokens, err := this.Lexer.Tokenizer(txt)
	if err != nil {
		return []Value{}, err
	}
	this.code = 0
	this.codes = [][]Token{tokens}
	var re []Value
	(*Parser).load_code(this)
	for this.code < uint64(len(this.codes)) {
		result, err := (*Parser).new_line(this)
		if err != nil {
			return []Value{}, err
		}
		if result == nil {
			(*Parser).next_code(this)
			continue
		}
		if result.(Node).Tp != "null" {
			re = append(re, result)
		}
		(*Parser).next_code(this)
	}
	return re, nil
}
func (this *Parser) WriteCache(txt string, cache any) ([]Value, any) {
	file, _ := os.Create(cache.(string))
	nodes, err := this.Make_Nodes(txt)
	if err != nil {
		return nil, err
	}
	re, _ := json.Marshal(Cache{Nodes: nodes, Txt: txt})
	file.Write(re)
	file.Close()
	return nodes, err
}
func (this *Parser) Parse(txt string) ([]Value, any) {
	//if cache == nil {
	return this.Make_Nodes(txt)
	//}
	/*str, err := ioutil.ReadFile(cache.(string))
	if err != nil {
		return this.WriteCache(txt, cache)
	}
	var j map[string]any
	json.Unmarshal(str, &j)
	if j["txt"].(string) != txt {
		return this.WriteCache(txt, cache)
	}
	var re []Value
	for _, i := range j["nodes"].([]any) {
		v := conv_json_to_value(i.(map[string]any))
		re = append(re, v)
	}
	return re, nil*/
}
func sum_codes(codes [][][]Token) [][]Token {
	re := [][]Token{}
	for _, a1 := range codes {
		for _, a2 := range a1 {
			re = append(re, a2)
		}
	}
	return re
}
func reverse_codes(codes [][][]Token) [][][]Token {
	for i := 0; i < len(codes)/2; i++ {
		j := len(codes) - i - 1
		codes[i], codes[j] = codes[j], codes[i]
	}
	return codes
}
func (this *Parser) new_line() (Value, any) {
	re, err := (*Parser).expr(this)
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case "new line":
		(*Parser).next_tok(this)
		this.codes = sum_codes([][][]Token{{this.codes[this.code][:this.tok_pos]}, {this.codes[this.code][this.tok_pos:]}})
	}
	return re, err
}
func (this *Parser) expr() (Value, any) {
	re, err := (*Parser).booleans(this)
	tok := this.tok
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case "=":
		var n Value
		(*Parser).next_tok(this)
		n, err = (*Parser).expr(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "=", tok.line, tok.col)
		break
	case "&&":
		var n Value
		(*Parser).next_tok(this)
		n, err = (*Parser).expr(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "&&", tok.line, tok.col)
		break
	case "||":
		var n Value
		(*Parser).next_tok(this)
		n, err = (*Parser).expr(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "||", tok.line, tok.col)
		break
	case ")":
		(*Parser).next_tok(this)
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg8", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	}
	return re, err
}
func (this *Parser) booleans() (Value, any) {
	re, err := (*Parser).calc(this)
	tok := this.tok
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case "==":
		var n Value
		(*Parser).next_tok(this)
		n, err = (*Parser).booleans(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "==", tok.line, tok.col)
		break
	}
	return re, err
}
func (this *Parser) call() (Value, any) {
	re, err := (*Parser).term(this)
	tok := this.tok
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case "(":
		v, err := (*Parser).Param(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, v}, "call", tok.line, tok.col)
		break
	}
	return re, err
}
func (this *Parser) calc() (Value, any) {
	re, err := (*Parser).call(this)
	tok := this.tok
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case "+":
		var n Value
		(*Parser).next_tok(this)
		if is_end_code(this.tok) {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		n, err = (*Parser).calc(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "+", tok.line, tok.col)
		break
	case "-":
		var n Value
		(*Parser).next_tok(this)
		if is_end_code(this.tok) {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		n, err = (*Parser).calc(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "-", tok.line, tok.col)
		break
	case "*":
		var n Value
		(*Parser).next_tok(this)
		if is_end_code(this.tok) {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		n, err = (*Parser).calc(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "*", tok.line, tok.col)
		break
	case "/":
		var n Value
		(*Parser).next_tok(this)
		if is_end_code(this.tok) {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		n, err = (*Parser).calc(this)
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, n}, "/", tok.line, tok.col)
		break
	}
	return re, err
}

/*
	func (this *Parser) term1() (Value, any) {
		re, err := (*Parser).term2(this)
		if err != nil {
			return nil, err
		}
		switch this.tok.tp {
		case ".":
			v, err := this.term1()
			break
		}
		return re
	}
*/
func (this *Parser) term() (Value, any) {
	re, err := (*Parser).factor(this)
	tok := this.tok
	if err != nil {
		return nil, err
	}
	switch this.tok.tp {
	case ".":
		(*Parser).next_tok(this)
		v, err := this.term()
		if err != nil {
			return nil, err
		}
		re = Create_Node([]Value{re, v}, "get attr", tok.line, tok.col)
		break
	}
	return re, err
}
func (this *Parser) factor() (Value, any) {
	tok := this.tok
	switch this.tok.tp {
	case "value":
		(*Parser).next_tok(this)
		return Create_Node([]Value{tok.value}, "value", tok.line, tok.col), nil
	case "var":
		(*Parser).next_tok(this)
		return Create_Node([]Value{tok.value}, "var", tok.line, tok.col), nil
	case "-":
		(*Parser).next_tok(this)
		n, err := this.factor()
		if n.(Node).Tp == "null" {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		return Create_Node([]Value{n}, "inverse number", tok.line, tok.col), err
	case "if":
		(*Parser).next_tok(this)
		n, err := this.expr()
		if n.(Node).Tp == "null" {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		var n2 Value
		n2, err = this.Enter_Code()
		if err != nil {
			return nil, err
		}
		return Create_Node([]Value{n, n2}, "if", tok.line, tok.col), err
	case "create global":
		(*Parser).next_tok(this)
		n, err := this.expr()
		return Create_Node([]Value{n}, "create global", tok.line, tok.col), err
	case "create local":
		(*Parser).next_tok(this)
		n, err := this.expr()
		return Create_Node([]Value{n}, "create local", tok.line, tok.col), err
	case "function":
		(*Parser).next_tok(this)
		vn := ""
		if this.tok.tp == "var" {
			vn = this.tok.value.Re_string("")
			(*Parser).next_tok(this)
		}
		parameters, err := (*Parser).Param(this)
		if err != nil {
			return nil, err
		}
		code, err := (*Parser).Enter_Code(this)
		if err != nil {
			return nil, err
		}
		return Create_Node([]Value{Create_String(vn), parameters, code}, "function", tok.line, tok.col), nil
	case "return":
		(*Parser).next_tok(this)
		n, err := this.expr()
		return Create_Node([]Value{n}, "return", tok.line, tok.col), err
	case "exist":
		(*Parser).next_tok(this)
		n, err := this.factor()
		if n.(Node).Tp != "var" {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		return Create_Node([]Value{n}, "exist", tok.line, tok.col), err
	case "(":
		(*Parser).next_tok(this)
		if this.tok.tp == ")" {
			(*Parser).next_tok(this)
			return Create_Node([]Value{Create_Null()}, "value", tok.line, tok.col), nil
		}
		code, err := (*Parser).expr(this)
		if err != nil {
			return nil, err
		}
		if this.tok.tp == ")" {
			(*Parser).next_tok(this)
		} else {
			return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg8", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
		}
		return Create_Node([]Value{code}, "()", tok.line, tok.col), nil
	case "pointer":
		(*Parser).next_tok(this)
		n, err := this.expr()
		return Create_Node([]Value{n}, "pointer", tok.line, tok.col), err
	case "+":
		(*Parser).next_tok(this)
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "*":
		(*Parser).next_tok(this)
		/*
			n, err := this.expr()
			if n.(Node).Tp != "var" {
			}
		return Create_Node([]Value{n}, "*f", tok.line, tok.col), err*/
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "/":
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "&&":
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "||":
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "==":
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg6", nil), line: tok.line, col: tok.col, lines: strings.Split(this.txt, "\n")}
	case "new line":
		(*Parser).next_tok(this)
		return this.factor()
	default:
		return Create_Node([]Value{}, "null", tok.line, tok.col), nil
		//return Create_Node(nil, ""), Error{msg: lang_text("erro1", []string{}) + lang_text("erro msg3", []string{}), line: this.tok.line, col: this.tok.col, lines: strings.Split(this.txt, "\n")}
	}
}
func (this *Parser) Param() (Value, any) {
	line, col := this.tok.line, this.tok.col
	if this.tok.tp == "(" {
		(*Parser).next_tok(this)
	} else {
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg8", nil), line: line, col: col, lines: strings.Split(this.txt, "\n")}
	}
	sp, walk := splitTokens(this.codes[this.code][this.tok_pos:], ",", "(", ")")
	i := 0
	for i < walk {
		(*Parser).next_tok(this)
		i++
	}
	old_codes := this.codes
	old_tok_pos := this.tok_pos
	old_code := this.code
	v, err := this.Make_Nodes_Toks(sp, "", false)
	if err != nil {
		return nil, err
	}
	this.codes = old_codes
	this.code = old_code
	this.tok_pos = old_tok_pos
	this.load_tok()
	if this.tok.tp == ")" {
		(*Parser).next_tok(this)
	} else {
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg8", nil), line: line, col: col - 1, lines: strings.Split(this.txt, "\n")}
	}
	return Create_Node(v, "Parameters", line, col), nil
}
func (this *Parser) Enter_Code() (Value, any) {
	line, col := this.tok.line, this.tok.col
	if this.tok.tp == "{" {
		(*Parser).next_tok(this)
	} else {
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg7", nil), line: line, col: col, lines: strings.Split(this.txt, "\n")}
	}
	i := 0
	vw := 1
	sp := []Token{}
	for vw > 0 {
		if this.tok.tp == "{" {
			vw += 1
		} else if this.tok.tp == "}" {
			vw -= 1
		}
		if vw <= 0 {
			break
		}
		sp = append(sp, this.tok)
		(*Parser).next_tok(this)
		if is_end_code(this.tok) {
			(*Parser).next_code(this)
		}
		if this.code >= uint64(len(this.codes)) {
			break
		}
		i++
	}
	old_codes := this.codes
	old_tok_pos := this.tok_pos
	old_code := this.code
	v, err := this.Make_Nodes_Toks([][]Token{sp}, "", true)
	if err != nil {
		return nil, err
	}
	this.codes = old_codes
	this.code = old_code
	this.tok_pos = old_tok_pos
	this.load_tok()
	if this.tok.tp == "}" {
		(*Parser).next_tok(this)
	} else {
		return nil, Error{msg: lang_text("erro1", nil) + lang_text("erro msg7", nil), line: line, col: col, lines: strings.Split(this.txt, "\n")}
	}
	return Create_Node(v, "{}", line, col), nil
}

type Variable struct {
	Pos      int
	Value    Value
	name     string
	is_const bool
}
type Interpreter struct {
	Globals *Object
	parser  Parser
	Debug   bool
}

func (this *Interpreter) Eval(txt string, locals *Object) Value {
	nodes, err := this.parser.Parse(txt)
	if this.Debug {
		println("{")
		for i := range nodes {
			println(fmt.Sprint(i) + ":" + nodes[i].Re_string(" "))
		}
		println("}")
		fmt.Println(err)
	}
	isTry := false
	if livre(err, nil) {
		i := uint64(0)
		for i < uint64(len(nodes)) {
			re, err := (*Interpreter).exec_node(this, nodes[i], locals, 0)
			if err != nil {
				e := err.(Error)
				e.lines = strings.Split(txt, "\n")
				err = e
			}
			if !isTry && !livre(err, nil) {
				return Create_Null()
			}
			if i == uint64(len(nodes)-1) {
				return re
			}
			i++
		}
	}
	return Create_Null()
}
func (this *Interpreter) Exec(txt string, locals *Object) any {
	nodes, err := this.parser.Parse(txt)
	if this.Debug {
		println("{")
		for i := range nodes {
			println(nodes[i].Re_string("	"))
		}
		println("}")
	}
	if livre(err, nil) {
		i := uint64(0)
		for i < uint64(len(nodes)) {
			_, err := (*Interpreter).exec_node(this, nodes[i], locals, 0)
			if err != nil {
				e := err.(Error)
				e.lines = strings.Split(txt, "\n")
				err = e
				return err
			}
			i++
		}
	}
	return nil
}
func PrintValue(args []Value, line bool) {
	str := ""
	for i, s := range args {
		if i > 0 {
			str += " "
		}
		str += s.Re_string("")
	}
	if line {
		println(str)
	} else {
		print(str)
	}
}
func set_obj_global(obj Object, glob *Object) Object {
	for i := range obj.value {
		if obj.value[i].Value.VType() == "Object" {
			obj.On_set_attr(i, set_obj_global(obj.On_get_attr(i).(Object), glob))
		} else if obj.value[i].Value.VType() == "Function" {
			f := obj.value[i].Value.(Function)
			f.locals = glob
			obj.On_set_attr(i, f)
		}
	}
	return obj
}
func Create_Context(inter *Interpreter) Value {
	return Create_Object(map[string]Value{
		"import": Create_Object(map[string]Value{
			"module": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
				src := ""
				if len(args) >= 1 {
					src = args[0].Re_string("")
				} else if _, ok := kwargs["src"]; ok {
					src = kwargs["src"].Value.Re_string("")
				}
				return inter.Exec_Module(src)
			}),
			"func": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
				src := ""
				if len(args) >= 1 {
					src = args[0].Re_string("")
				} else if _, ok := kwargs["src"]; ok {
					src = kwargs["src"].Value.Re_string("")
				}
				txt, err := os.ReadFile(src)
				if err != nil {
					return Create_Null(), err
				}
				nodes, _ := inter.parser.Parse(string(txt))
				l := Create_Object(make(map[string]Value)).(Object)
				return Create_Function(nodes, inter, &l, []Variable{}), nil
			}),
		}),
		"exec": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			code := ""
			locals := Create_Object(nil)
			if len(args) >= 1 {
				code = args[0].Re_string("")
			} else if _, ok := kwargs["text"]; ok {
				code = kwargs["text"].Value.Re_string("")
			}
			if len(args) >= 2 {
				locals = args[2]
			} else if _, ok := kwargs["locals"]; ok {
				locals = kwargs["text"].Value
			}
			l := locals.(Object)
			err := inter.Exec(code, &l)
			return l, err
		}),
		"globals": inter.Globals,
	})
}
func (this *Interpreter) Init() {
	var g Object = Create_Object(map[string]Value{}).(Object)
	this.Globals = &g
	this.Set_Global("console", Create_Object(map[string]Value{
		"log": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			PrintValue(args, true)
			return Create_Null(), nil
		}),
		"write": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			PrintValue(args, false)
			return Create_Null(), nil
		}),
		"read": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			reader := bufio.NewReader(os.Stdin)
			PrintValue(args, false)
			str, _ := reader.ReadString(byte('\n'))
			str = str[:len(str)-1]
			return Create_String(str), nil
		}),
	}), true)
	this.Set_Global("Math", Create_Object(map[string]Value{
		"pi": Create_Number(3.1415),
		"cos": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var v float64
			if len(args) >= 1 {
				v = args[0].Re_number()
			} else if _, ok := kwargs["x"]; ok {
				v = kwargs["x"].Value.Re_number()
			}
			return Create_Number(math.Cos(v)), nil
		}),
		"sin": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var v float64
			if len(args) >= 1 {
				v = args[0].Re_number()
			} else if _, ok := kwargs["x"]; ok {
				v = kwargs["x"].Value.Re_number()
			}
			return Create_Number(math.Sin(v)), nil
		}),
		"atan": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var x1 float64
			if len(args) >= 1 {
				x1 = args[0].Re_number()
			} else if _, ok := kwargs["x1"]; ok {
				x1 = kwargs["x1"].Value.Re_number()
			}
			var x2 float64
			if len(args) >= 2 {
				x2 = args[1].Re_number()
			} else if _, ok := kwargs["x2"]; ok {
				x2 = kwargs["x2"].Value.Re_number()
			}
			return Create_Number(math.Atan(x1 - x2)), nil
		}),
		"atan2": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var x1 float64
			if len(args) >= 1 {
				x1 = args[0].Re_number()
			} else if _, ok := kwargs["x1"]; ok {
				x1 = kwargs["x1"].Value.Re_number()
			}
			var x2 float64
			if len(args) >= 2 {
				x2 = args[1].Re_number()
			} else if _, ok := kwargs["x2"]; ok {
				x2 = kwargs["x2"].Value.Re_number()
			}

			var y1 float64
			if len(args) >= 3 {
				y1 = args[2].Re_number()
			} else if _, ok := kwargs["y1"]; ok {
				y1 = kwargs["y1"].Value.Re_number()
			}
			var y2 float64
			if len(args) >= 4 {
				y2 = args[3].Re_number()
			} else if _, ok := kwargs["y2"]; ok {
				y2 = kwargs["y2"].Value.Re_number()
			}
			return Create_Number(math.Atan2(y1-y2, x1-x2)), nil
		}),
		"floor": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var v float64
			if len(args) >= 1 {
				v = args[0].Re_number()
			} else if _, ok := kwargs["x"]; ok {
				v = kwargs["x"].Value.Re_number()
			}
			return Create_Number(math.Floor(v)), nil
		}),
		"ceil": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var v float64
			if len(args) >= 1 {
				v = args[0].Re_number()
			} else if _, ok := kwargs["x"]; ok {
				v = kwargs["x"].Value.Re_number()
			}
			return Create_Number(math.Ceil(v)), nil
		}),
		"abs": Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
			var v float64
			if len(args) >= 1 {
				v = args[0].Re_number()
			} else if _, ok := kwargs["x"]; ok {
				v = kwargs["x"].Value.Re_number()
			}
			return Create_Number(math.Abs(v)), nil
		}),
	}), true)
	this.Set_Global("ctx", Create_Context(this), true)
	this.Set_Global("true", Create_Bool(true), true)
	this.Set_Global("false", Create_Bool(false), true)
	this.Set_Global("KllContext", Create_GoFunction(func(args []Value, kwargs map[string]*Variable, pos int) (Value, any) {
		i := Interpreter{}
		return Create_Context(&i), nil
	}), true)
}
func (this *Interpreter) Set_Global(name string, value Value, is_const bool) {
	this.Globals.Create_Var(name, 0, value, is_const)
}
func (this *Interpreter) exec_node(nodeV Value, locals *Object, pos int) (Value, any) {
	node := nodeV.(Node)
	switch node.Tp {
	case "value":
		return node.Value[0], nil
	case "()":
		obj, err := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err != nil {
			return nil, err
		}
		return obj, nil
	case "get attr":
		obj, err := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err != nil {
			return nil, err
		}
		return Get_attr(obj, node.Value[1]), nil
	case "function":
		args := []Variable{}
		for _, v := range node.Value[1].(Node).Value {
			nodeVa := v.(Node)
			if nodeVa.Tp == "var" {
				args = append(args, Variable{name: nodeVa.Value[0].Re_string(""), Value: Create_Null()})
			}
		}
		function := Create_Function(node.Value[2].(Node).Value, this, locals, args)
		if node.Value[0].Re_string("") != "" {
			locals.Create_Var(node.Value[0].Re_string(""), pos, function, false)
		}
		return function, nil
	case "exist":
		ok := In(locals, node.Value[0])
		if ok.Re_bool() {
			return ok, nil
		}
		ok = this.Globals.On_in(node.Value[0])
		if ok.Re_bool() {
			return ok, nil
		}
		return Create_Bool(false), nil
	case "call":
		args := []Value{}
		kwargs := make(map[string]*Variable)
		for _, v := range node.Value[1].(Node).Value {
			nodeVa := v.(Node)
			if nodeVa.Tp == "=" {
				v, err := (*Interpreter).exec_node(this, nodeVa.Value[1], locals, pos)
				if err != nil {
					return Create_Null(), err
				}
				kwargs[nodeVa.Value[0].(Node).Value[0].Re_string("")] = &Variable{Value: v, name: nodeVa.Value[0].(Node).Value[0].Re_string("")}
			} else {
				v, err := (*Interpreter).exec_node(this, nodeVa, locals, pos)
				if err != nil {
					return Create_Null(), err
				}
				args = append(args, v)
			}
		}
		obj, err := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err != nil {
			return Create_Null(), err
		}
		var re Value
		re, err = Call(obj, args, kwargs, pos)
		var re_err any
		if err != nil {
			re_err = Error{line: node.Value[0].(Node).Line, col: node.Value[0].(Node).Col, other_error: err}
		}
		return re, re_err
	case "pointer":
		if node.Value[0].(Node).Tp == "get attr" {

		} else {
			v, ok := locals.On_get_Variable(node.Value[0].(Node).Value[0].Re_string(""))
			if ok {
				return Create_Pointer(v), nil
			}
			v, ok = this.Globals.On_get_Variable(node.Value[0].(Node).Value[0].Re_string(""))
			if ok {
				return Create_Pointer(v), nil
			}
		}
		break
	case "*f":
		if node.Value[0].(Node).Tp == "=" {
			obj, err := (*Interpreter).exec_node(this, node.Value[0].(Node).Value[1], locals, pos)
			if node.Value[0].(Node).Tp == "get attr" {
				v, ok := locals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""))
				ok = ok && !v.is_const
				if ok {
					if node.Value[0].(Node).Value[1].(Node).Tp == "get attr" {
						v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node), obj)
					} else {
						v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node).Value[0], obj)
					}
				}
				v, ok = this.Globals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""))
				ok = ok && !v.is_const
				if ok {
					if node.Value[0].(Node).Value[1].(Node).Tp == "get attr" {
						v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node), obj)
					} else {
						v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node).Value[0], obj)
					}
				}
				return v.Value, err
			} else {
				v, ok := locals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""))
				ok = ok && !v.is_const
				if ok {
					if v.Value.VType() == Pointer.VType(Pointer{}) {
						v.Value.(Pointer).value.Value = obj
					}
					return v.Value, err
				}
				v, ok = this.Globals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""))
				ok = ok && !v.is_const
				if ok {
					if v.Value.VType() == Pointer.VType(Pointer{}) {
						v.Value.(Pointer).value.Value = obj

					}
					return v.Value, err
				}
			}
		}
		break
	case "create global":
		if node.Value[0].(Node).Tp == "var" {
			this.Globals.Create_Var(node.Value[0].(Node).Value[0].Re_string(""), 0, Create_Null(), false)
		} else if node.Value[0].(Node).Tp == "=" {
			obj, err := (*Interpreter).exec_node(this, node.Value[0].(Node).Value[1], locals, pos)
			this.Globals.Create_Var(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""), 0, obj, false)
			return obj, err
		}
		break
	case "create local":
		if node.Value[0].(Node).Tp == "var" {
			locals.Create_Var(node.Value[0].(Node).Value[0].Re_string(""), pos, Create_Null(), false)
		} else if node.Value[0].(Node).Tp == "=" {
			obj, err := (*Interpreter).exec_node(this, node.Value[0].(Node).Value[1], locals, pos)
			locals.Create_Var(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string(""), pos, obj, false)
			return obj, err
		}
		break
	case "var":
		ok := In(locals, node).(Bool).Value
		if ok {
			return locals.On_get_attr(node.Value[0].Re_string("")), nil
		}
		ok = In(this.Globals, node).(Bool).Value
		if ok {
			return this.Globals.On_get_attr(node.Value[0].Re_string("")), nil
		}
		return Create_Null(), Error{msg: lang_text("erro2", []string{}) + lang_text("err msg5", []string{node.Value[0].Re_string("")}), line: node.Line, col: node.Col}
	case "=":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		/*if node.Value[0].(Node).Tp == "get attr" {
			v, ok := locals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string())
			ok = ok && !v.is_const
			if ok {
				if node.Value[0].(Node).Value[1].(Node).Tp == "get attr" {
					v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node), v1)
				} else {
					v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node).Value[0], v1)
				}
			}
			v, ok = this.Globals.On_get_Variable(node.Value[0].(Node).Value[0].(Node).Value[0].Re_string())
			ok = ok && !v.is_const
			if ok {
				if node.Value[0].(Node).Value[1].(Node).Tp == "get attr" {
					v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node), v1)
				} else {
					v.Value = Set_attr(v.Value, node.Value[0].(Node).Value[1].(Node).Value[0], v1)
				}
			}
			return v.Value, err1
		} else {*/
		ok := In(locals, node.Value[0]).(Bool).Value
		if err1 != nil {
			return nil, err1
		}
		if ok {
			if (node.Value[0].(Node).Tp == "get attr" && locals.value[node.Value[0].(Node).Value[0].(Node).Value[0].Re_string("")].is_const) || (node.Value[0].(Node).Tp == "var" && locals.value[node.Value[0].(Node).Value[0].Re_string("")].is_const) {
				return v1, nil
			}
			*locals = Set_attr(locals, node.Value[0], v1).(Object)
			return v1, nil
		}
		ok = In(this.Globals, node.Value[0]).(Bool).Value
		if ok {
			if (node.Value[0].(Node).Tp == "get attr" && this.Globals.value[node.Value[0].(Node).Value[0].(Node).Value[0].Re_string("")].is_const) || (node.Value[0].(Node).Tp == "var" && this.Globals.value[node.Value[0].(Node).Value[0].Re_string("")].is_const) {
				return v1, nil
			}
			*this.Globals = Set_attr(this.Globals, node.Value[0], v1).(Object)
			return v1, nil
		}
		return Create_Null(), Error{msg: lang_text("erro2", []string{}) + lang_text("err msg5", []string{node.Value[0].(Node).Value[0].Re_string("")}), line: node.Line, col: node.Col}
		//}
	case "inverse number":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		return Create_Number(-v1.Re_number()), nil
	case "+":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Sum(v1, v2), nil
	case "-":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Sub(v1, v2), nil
	case "*":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Mul(v1, v2), nil
	case "/":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Div(v1, v2), nil
	case "==":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Create_Bool(v1 == v2), nil
	case "&&":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Create_Bool(v1.Re_bool() && v2.Re_bool()), nil
	case "||":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		v2, err2 := (*Interpreter).exec_node(this, node.Value[1], locals, pos)
		if err2 != nil {
			return nil, err2
		}
		return Create_Bool(v1.Re_bool() || v2.Re_bool()), nil
	case "if":
		v1, err1 := (*Interpreter).exec_node(this, node.Value[0], locals, pos)
		if err1 != nil {
			return nil, err1
		}
		if v1.Re_bool() {
			for _, nv := range node.Value[1].(Node).Value {
				_, err := this.exec_node(nv, locals, pos+1)
				if err != nil {
					return nil, err
				}
			}
			for v := range locals.value {
				if locals.value[v].Pos >= pos+1 {
					delete(locals.value, v)
				}
			}
		}
	}

	return Create_Null(), nil
}
func (this *Interpreter) Exec_Main(src string) Value {
	this.Init()
	locals := Create_Object(make(map[string]Value)).(Object)
	locals.Create_Var("__name__", 0, Create_String("__main__"), true)
	txt, _ := os.ReadFile(src)
	re := (*Interpreter).Eval(this, string(txt), &locals)
	if this.Debug {
		fmt.Println(re.Re_string(""))
	}
	return re
}
func (this *Interpreter) Exec_Module(src string) (Object, any) {
	this.Init()
	locals := Create_Object(make(map[string]Value)).(Object)
	locals.Create_Var("__name__", 0, Create_String("__main__"), true)
	txt, _ := os.ReadFile(src)
	err := (*Interpreter).Exec(this, string(txt), &locals)
	locals = set_obj_global(locals, &locals)
	return locals, err
}

func conv_error_in_str(e Error) string {
	re := e.msg
	re += " line:" + fmt.Sprint(e.line) + ",collum:" + fmt.Sprint(e.col) + "\n"
	re += e.lines[e.line-1] + "\n"

	for i := range e.lines[e.line-1] {
		if uint64(i) == e.col-1 {
			re += "^"
		} else {
			re += " "
		}
	}
	return re
}
func livre(e any, txt []string) bool {
	if e != nil {
		a := e.(Error)
		txt = append(txt, conv_error_in_str(a)+"\n")
		if a.other_error != nil {
			v := a.other_error.(Error)
			v.lines = a.lines
			livre(v, txt)
			return false
		}
		for _, v := range txt {
			print(v)
		}
		return false
	}

	return true
}
