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


----

## Conditionals

we have to figure out decide when to evaludate what. that's the point of conditions


```js
if(x > 10){
  puts("all good!");
}else{
  puts(" x is too low!");
  shutdownSystem();
}
```
> when evaluating the if-else expression the important thing is to only evaluate 
the correct branch. if the condition is met. if it isnt met we must only eval
the else branch

in case of monkey, the consequence part of the conditional will be evaluated when the 
condition is 'truthy'. and 'truthy' means: it's not `null` and its not `false`. 
it doesn't necessarily need to be `true`

```js
let x = 10;
if (x) {
  puts("everything okay!");
} else {
  puts("x is too high!");
  shutdownSystem();
}
```

>in this case `"everything okay!"` should be printed. because `x` is bound to `10`, evaluates to `10`
and `10` is not `null` not `false`. 

Also in monkey the conditional does't evaluate to a value it's suppose to return `NULL`


# Return Statements
```go
5 * 5 * 5;
return 10;
9 * 9 * 9;
```
when evaludated this should return 10. and most important thing is the last line `9 * 9 * 9;`
should not get evaluated since its a early return. 

There are a few different ways to implement return statements. In some host languages 
we could use gotos or exceptions. But in Go a `rescue` or `catch` are not easy to 
come by and we don’t really have the option of using gotos in a clean way.

https://en.wikipedia.org/wiki/Goto


```go
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// ...
	case *ast.ReturnStatement:
		// we eval the expression associated with the return statement.
		// if the last eval expression associated with return statement. we then
		// wrap the result of this call to `Eval` in our new `object.ReturnValue` so we can
		// keep track on this
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
		// in case the last eval result is a `object.ReturnValue` if so we stop the
		// evaluation and return the unwrapped value. we dont need to return an `object.ReturlValue`
		// but only the value its wrapping which is what the user expected.

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}
```

the problem with our implementation is we dont keep track of `object.ReturlValues` for longer 
and can't unwrap their values on the first enconter. for example in the below example 


```go
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
```

this should return 10 but our code will return `1`. 

```go
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)
		// when we have block statements we cant unwrap the result in the first sight, because
		// we need to furthure keep track of its so we can stop execution in outermost block
		// statement.

		// here we dont explicity don't unwrap the return value and only check the `Type()`
		// of each evaluation result, if it's `object.RETURN_VALUE_OBJ`
		// we return the value. without unwrapping it's `.Value` so it stops the execution
		// in a possible outer block statement and goes to `evalProgram` where the value gets unwrapped
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}
```

To be able to handle this when evaluating block statement we check for the `result.Type()` if its 
`object.RETURN_VALUE_OBJ` we will return the value. 

Note that we dont unwrap the value because when evalProgram is getting executed the value will 
get unwrapped

---

# Bindings & The Environment

we need to add support to `let` statements. not only we need to support let statements 
we need to support evaludation of identifiers, too

```js
let x = 5 * 5;
```
we need to make sure that x evaluates to 25 after interpreting the line above.

> So we are going to evaludate let statements and identifiers. we evaludate
let statements by evaluating their value-producing expression and keep track
of the produces value under the specified name.

> to eval identifies we check if we already have a value bound to the name.
if we do, the identifies eval to this value, and if we don't we return an error


```go
case *ast.LetStatement:
	val := Eval(node.Value)
	if isError(val) {
		return val
	}
	// how do we keep track on this val here ?
```

As the comment mentions. now what we evaludate the value of the node. but how do we
keep track of the value here?

This is where environment comes to play. 

> Environment is what we use to keep track of value by associating them with a name. 

The name "Environment" is a classic one, used in lot of other interpreters,
especially Lispy ones. but eventho the name might sound involving, at it's heart
the environment is a hash map that associates strings with objects. 

and that's exactly what we're goingh to use for our implementation

```go
// object/environment.go
package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
```
we are going to change our function definition of our eval to pass the env as well

```go
// evaluator/evaluator.go

func Eval(node ast.Node, env *object.Environment) object.Object {
// [...]
```

we can initiate this in our repl

```go
// repl/repl.go

import (
    // [...]
    "monkey/object"
)

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()

    for {
// [...]
        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
        }
    }
}

// evaluator/evaluator_test.go

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    env := object.NewEnvironment()

    return Eval(program, env)
}
```

to memorize the values of the variables we can set the result after eval to the env 

```go
// evaluator/evaluator.go

func Eval(node ast.Node, env *object.Environment) object.Object {
// [...]
    case *ast.LetStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)
// [...]
```
when adding associations to the env let statements we also need to get those values 
to that. we are going to implement evalIdentifier doing this is pretty easy

```go
// evaluator/evaluator.go

func Eval(node ast.Node, env *object.Environment) object.Object {
// [...]
    case *ast.Identifier:
        return evalIdentifier(node, env)
// [...]
}

func evalIdentifier(
    node *ast.Identifier,
    env *object.Environment,
) object.Object {
    val, ok := env.Get(node.Value)
    if !ok {
        return newError("identifier not found: " + node.Value)
    }
    // for now `evalIdentifier` simply check if the value availible if so
    // return the value or throw an error
    return val
}
```
now the variables work as expected
```js
Hello dasith, This is Monkey programming language!
Feel free to type in commands
>> let a = 5;
>> let b = 10;
>> a + b;
15
>> xyx
ERROR: identifier not found: xyx
>>
```
---

# Functions & Function Calls

```go
// object/object.go
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
```

The definition of `object.Function` has the `Parameters` and `Body` fields. But it also 
has a `Env`. this holds a pointer to a `object.Environment` because functions in Monkey
carry their own environment with them. This allows for closures, which "close over"
the environment the're defined and can later access it.


```go
// evaluator/evaluator.go

func Eval(node ast.Node, env *object.Environment) object.Object {
// [...]
    case *ast.FunctionLiteral:
        params := node.Parameters
        body := node.Body
        return &object.Function{Parameters: params, Env: env, Body: body}
// [...]
```
we can build the internal representation of funcions this way.

### Implementation of function appication.

the below are the test cases we are trying to implement

```go 
// evaluator/evaluator_test.go

func TestFunctionApplication(t *testing.T) {
    tests := []struct {
        input    string
        expected int64
    }{
        {"let identity = fn(x) { x; }; identity(5);", 5},
        {"let identity = fn(x) { return x; }; identity(5);", 5},
        {"let double = fn(x) { x * 2; }; double(5);", 10},
        {"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
        {"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
        {"fn(x) { x; }(5)", 5},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}
```
each test here does thesame thing, define function, apply it to args and then make an 
assertion about value.

> We are also testing two possible forms of *ast.CallExpression here. 
One where the function is an identifier that evaluates to a function object, 
and the second one where the function is a function literal. 

The neat thing is that it doesn’t really matter. 
We already know how to evaluate identifiers and function literals:

```go
// evaluator/evaluator.go

func Eval(node ast.Node, env *object.Environment) object.Object {
// [...]
    case *ast.CallExpression:
        function := Eval(node.Function, env)
        if isError(function) {
            return function
        }
// [...]
}
```
we are just using eval to get function we want to call. whether that's a
`ast.Identifier` or an `*ast.FunctionLiteral`: Eval returns `*object.Function`

But how do we call this `*object.Function`? 

first we want to eval the ags of a call expression

```js
let add = fn(x, y) { x + y };
add(2 + 2, 5 + 5);
```


```go
	case *ast.CallExpression:
		// we are just using eval to get function we want to call. whether that's a
		// `ast.Identifier` or an `*ast.FunctionLiteral`
		// Eval returns `*object.Function`
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpression(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
	}
	return nil
}

func evalExpression(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}
```

we can handle this like this, we just iterate over a list of `ast.Expressions` and eval them in
context of current env. if we have an error we return error.

Here we eval the args from left-to-right. Hopefully we wont be writing monkey code 
that makes assertions about the order of argument eval. 

but if we do, we’re on the conservative and safe side of programming language design.

### Now how do we call the function? 

we have to eval our function body. but we have to consider arguments. 

The body of the function contains references to the params of the function and just 
eval the body of function in the current env would result in unknown references. 

- we need to change the env which the function is evaluated. so references to params in
function body resolve to correct args. 

- but we can't just add these args to current env. this can lead to previous bindings being 
over-written. 

```js
let i = 5;
let printNum = fn(i) {
  puts(i); // 10
};

printNum(10);
puts(i); // 5
```
if we over-write the current env the last line will also be 10. 

> So adding args to function call to current env to make them accessable wont work. 
what we need is preserve previous bindings while at the same time make new ones availible

#### Extending the environment

Extending the env means we create new instance of `object.Environment` with pointer to env 
it should extend. 

When new env `Get` methord is called and doesnt have a value associated with the given name.
it call the `Get` of the enclosing env. that's the env it's extending. 

if that enclosing env cant't find the value, it calls its own enclosing env and so untill 
there is  no encolsing env anymore and we can safely say we have an `ERROR: unknown identifier: foobar`
```go
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	// we are adding a new field called `outer` this contains a reference to another
	// `object.Environment` which is the enclosing env, the only one its extending
	outer *Environment
}
```
we have to change the `Get` implementation to check the enclosing environment for given name too

```go
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}
```

> This new behaviour mirrors how we think about variable scopes. There are an inner scope
and outer scope.  if value not found in inner scope, its looked up in outer scope. 
this outer scope *encloses* inner scope. and the inner scope *extends* the outer one

---

## Closures

The reason why we extend the function environment and not the current env is because closures

```go
func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
  fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}
```

```
Hello dasith, This is Monkey programming language!
Feel free to type in commands
>> let newAdder = fn(x) { fn(y) { x + y }; };
>> let addTwo = newAdder(2);
>> addTwo(3);
5
>> let addThree = newAdder(3);
>> addThree(10);
13
>>
```

Closures are functions that "close over" the env they were defined in. they have their own
env around and wheneve they're called they can access it.

```js
let newAdder = fn(x) { fn(y) { x + y }; };
let addTwo = newAdder(2);
```
`newAdder` here is a higher-order function. Higher order functions are functions either return
other functions or receive them as args. in this case `newAdder` returns another function. a closure.

`addTwo` is bound to the closure that's returned when calling newAdder with 2 as the sole arg

when `addTwo` is called it not only has access to arg of the call, (y) it also has
access to the value `x` which was bound at the time of the `newAdder(2)` call

> the closure `addTwo` still has acces sto env that was the current env at the time its definition.
which is when the last time of `newAdder` body was evaludated. this is a funcion literal 
remeber we build an `object.Function` and keep the reference to the current env in it's `.Env` field.



