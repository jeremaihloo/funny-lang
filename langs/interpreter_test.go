package langs

import (
	"fmt"
	"testing"
)

func RunSingle(data interface{}) (*Interpreter, Value) {
	i := NewInterpreterWithScope(make(map[string]Value))
	var d []byte
	switch data.(type) {
	case string:
		d = []byte(data.(string))
	}
	parser := NewParser(d)
	block := parser.Parse()
	r, _ := i.Run(Program{
		Statements: block,
	})
	return i, Value(r)
}

func TestInterpreter_Assign(t *testing.T) {
	i := NewInterpreterWithScope(make(map[string]Value))
	i.Assign("a", Value(1))
	flag := false
	var val interface{}
	for _, scope := range i.Vars {
		for k, v := range scope {
			if k == "a" {
				flag = true
				val = v
			}
		}
	}
	if !flag {
		t.Error("assign error key not in scope")
	} else {
		if val != 1 {
			t.Error("assign error value not equal 1")
		}
	}
	scope := Scope{}
	i.PushScope(scope)
	i.Assign("b", Value(2))
	v := i.Lookup("b")
	if v != 2 {
		t.Errorf("val not eq 2 %s", v)
	}
	i.Assign("a", Value(3))
	a := i.Lookup("a")
	if a != 3 {
		t.Errorf("a not eq 3 %s", a)
	}
	i.PopScope()
	v = i.LookupDefault("b", nil)
	if v != nil {
		t.Error("pop scope error")
	}
}

func TestInterpreter_Lookup(t *testing.T) {
	i := NewInterpreterWithScope(make(map[string]Value))
	i.Assign("a", Value(1))
	val := i.Lookup("a")
	if val != 1 {
		t.Error("lookup error")
	}
}

func TestInterpreter_EvalFunctionCall(t *testing.T) {
	i := NewInterpreterWithScope(make(map[string]Value))
	parser := NewParser([]byte("echo(1)"))
	i.Run(Program{
		parser.Parse(),
	})
}

func TestInterpreter_EvalFunctionCall2(t *testing.T) {
	i := NewInterpreterWithScope(make(map[string]Value))
	parser := NewParser([]byte("echo2(b){echo(b)} \n echo2(1)"))
	i.Run(Program{
		parser.Parse(),
	})
}

func TestInterpreter_EvalPlus(t *testing.T) {
	i := NewInterpreterWithScope(make(map[string]Value))
	parser := NewParser([]byte("  a = 1 + 1"))
	i.Run(Program{
		parser.Parse(),
	})
	a := i.Lookup("a")
	if a != 2 {
		t.Error("eval plus error")
	}
}

func TestInterpreter_Run(t *testing.T) {
	data := `
a = 1
b = 2
c = a + b

echo(c)

p(a, b){
    return a + b
}

d = p(a,b)

return d - 1`

	_, r := RunSingle(data)
	if r != 2 {
		t.Error("RunSingle funny.fun must return 2")
	}
}

func TestInterpreter_Return(t *testing.T) {
	data := `
testReturn(t){
    if t<1{
        return t
    }
    return testReturn(t-1)
}

t = testReturn(10)`
	_, r := RunSingle(data)
	ty := Typing(r)
	t.Log(ty)
	t.Log(r)
}

func TestInterpreter_Fib(t *testing.T) {
	data := `
fib(n) {
    echoln('n: ', n)
    if n < 2 {
      return n
    } else {
      return fib(n - 2) + fib(n - 1)
    }
}

return fib(5)`

	_, r := RunSingle(data)
	ty := Typing(r)
	t.Log(ty)
	t.Log(r)
}

func TestInterpreter_EvalBlock(t *testing.T) {
	data := `
a = 2
b = 1
if a > b {
return a
} else {
return b
}`

	_, r := RunSingle(data)
	if r != 2 {
		t.Error(fmt.Sprintf("RunSingle funny.fun must return 2 but got %s", r))
	}
}
