# Evaluation

> The most obvious and classical choice of what to do with AST is to just interpret it.
Traverse the AST, visit each node and do what the note signifies: 


Interpreters working this way are called `tree-walking interpreters` they are a archetype of interpreters



> Other interpreters also traverse the AST, but instead of interpreting the AST itself
they first convert it to bytecode. Bytecode is another IR of the AST and a really dense one at that. 
The exact format and of which opcodes (the instructions that make up the bytecode) 
it’s composed of varies and depends on the guest and host programming languages.


The opcodes are pretty similar to mnemonics of most assembly languages. it's safe to bet to say that
most bytecode definitions contains opcodes for `push` and `pop` to do stack operations.

But bytecode is no native machine code or assembly code. and it wont be executed by OS and the CPU. 
Insted it's interpreted by a **virtual machine**

The way this vms work is they emulate a machine that understands this particilar bytecode format. 


> A variation of this strategy doesn’t involve an AST at all. insted of building AST the parser emits
bytecode directly.Isnt emitting bytecode that gets interpreted (executed?) form of compilation 
This is where the line between interpreters or compilers become blurly.

> To make things even more fuzzy, some impelementations parse the source code build an ASt and convert 
AST to bytecode. but insted of executing operations in virtual machine it compiles bytecode to to 
native machine code. just in time. THis is called as JIT ( for `just in time`) interpreter/compiler


> Others skip compilation to bytecode. they recursively traverse the AST but before executing a 
branch of it the node is compiled to native machine code. then executed again `just in time`


A tree-walking interpreter that recursively evaluates an AST is probably the slowest of all approaches,
but easy to build, extend, reason about and as portable as the language it's implemented in.


Example: 

1. Ruby is a great example here. Up to and including version 1.8 the interpreter was a tree-walking interpreter,
executing the AST while traversing it. 
But with version 1.9 came the switch to a virtual machine architecture. 
Now the Ruby interpreter parses source code, builds an AST and then compiles this AST into bytecode, 
which gets then executed in a virtual machine. The increase in performance was huge.

2. The WebKit JavaScript engine JavaScriptCore and its interpreter named `Squirrelfish`
also used AST walking and direct execution as its approach. 
Then in 2008 came the switch to a virtual machine and bytecode interpretation. 
Nowadays the engine has four (!) different stages of JIT compilation, which kick in at different times 
in the lifetime of the interpreted program depending on which part of the program needs the best performance.

3. Another example is Lua. The main implementation of the Lua started out as an interpreter that compiles 
to bytecode and executes the bytecode in a register-based virtual machine. 
12 years after its first release another implementation of the language was born: LuaJIT. 
The clear goal of Mike Pall, the creator of LuaJIT, was to create the fastest Lua implementation possible.
And he did. By JIT compiling a dense bytecode format to highly-optimized machine code for
different architectures the LuaJIT implementation beats the original Lua in every benchmark. 
And not just by a tiny bit, no; it’s sometimes 50 times faster.


---


# A Tree-Walking Interpreter

Our interpreter will be a lot like a classic Lisp interpreter. The design we're going to use is heavily 
inspired by the interpreter presented in 'The Structure and Interpretation of Computer Programs' (SICP),
especially its usage of environments.

here's a psudocode of what we try to impl 

```js
function eval(astNode) {
  if (astNode is integerliteral) {
    return astNode.integerValue

  } else if (astNode is booleanLiteral) {
    return astNode.booleanValue

  } else if (astNode is infixExpression) {

    leftEvaluated = eval(astNode.Left)
    rightEvaluated = eval(astNode.Right)

    if astNode.Operator == "+" {
      return leftEvaluated + rightEvaluated
    } else if ast.Operator == "-" {
      return leftEvaluated - rightEvaluated
    }
  }
}
```
as you can see `eval` is recursive. when `astNode` is `infixExpression` is true, `eval` calls itself again
two times to evaluate the left and right operands of the infix expression. 

---

# Representing Objects

its more like a `value system` or `object representation`

the point is we need to define what our `eval` function returns. we need a system that can represent
the values our AST representation or values that we generate when evaluating the AST in memory

> I heartily recommended the [Wren source code](https://github.com/wren-lang/wren), which includes two types of value 
representation, enabled/disabled by using a compiler flag.


```go
package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// when we enconter an int in the source code we first turn it in to
// `ast.IntegerLiteral` and then, when eval the AST node we turn it in to
// `object.Integer` saving the value inside our struct and passing around
// reference to this struct
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
```
---

# Evaluating Expressions

`Eval` will take a `ast.Node` as a input and return a `object.Object`. 

```go
// the reason to do this is. we always traverse the AST, we should start at the top
// of the tree. receiving an `*ast.Program` and then traverse every node in it.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
```

----

## Prefix Expressions

This is where how we decide how our language works.


```go
// evaluator/evaluator_test.go

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}
```

above is a sample test for the bang operator. this is what we expect to evaluated.


```go
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
    // ...
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}
	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
```
above is the implementation for `!true` and other band operator, this is why we say this 
is the place where we decide how the language work. 

`!true` and `!false` are common sense. but in our lagnauge `!5` something others might throw
and error but in monkey this resut acts as "truthy"

---


## Infix Expressions

```go
// the reason to do this is. we always traverse the AST, we should start at the top
// of the tree. receiving an `*ast.Program` and then traverse every node in it.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
    // ...
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	default:
		return NULL
	}
}
```
Just like the `*ast.PrefixExpression` we eval the operand first. and now we have two;
the left and trhe right arm of AST node. 

after eval the operands we return the values and operator to `evalIntegerInfixExpressions`

in above we only check if both are int and pass them but we can extend this later.

the main logic of this is implemented in the `evalIntegerInfixExpression` as you can see
this is where the evaluation happens and we return the result

----

Monkey also supports boolean operands for equality operators `==` and `!=`.

```go
func evalInfixExpression(
    operator string,
    left, right object.Object,
) object.Object {
    switch {
// [...]
    case operator == "==":
        return nativeBoolToBooleanObject(left == right)
    case operator == "!=":
        return nativeBoolToBooleanObject(left != right)
    default:
        return NULL
    }
}
```
yes we can change our existing `evalInfixExpression` to get this to work. we are using 
`pointer comparison` here to check equality between bool. 

The reason why this works is we are always using pointers to our objects and in the case of
booleans we only ever use two `TRUE` or `FALSE`. if something has the same value as `TRUE`
in memory address then its `true`. 

same kind of implementation is for `NULL`

TODO: fix this

> This doesnt work for int or other data types . in the case of `*object.Integer` we are always
allocating new instances of `object.Integer` and thus use new pointers. we cant compare these
pointers in different instances. otherwise `5 == 5` would be false. in this case we explicity
compair the values and not the objects that wrap these values


This is why check of int operands has to be higher up in the switch statement and match earlier
than these newly added `case` branches. as long as we take care of other operand types 
before arriving at the pointer comparisions we are fine and it works.

> This is the reason why integer comparision in Monkey is slower than boolean comparision

Reason: 

> Monkey’s object system doesn’t allow pointer comparison for integer objects


