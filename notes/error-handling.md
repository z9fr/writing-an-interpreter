# Error Handling 

What we are trying to implement is *not* user defined exceptions. It's internal error handling.
Errors for wrong operators, unsupported operations and other user or internal errors that may 
arise during execution

> the error handling is implement in a same way as handling return statements. the reason 
for this is easy to find: errors and return statements both stop the evaluation of a seriese
of statements

we first need to define error object

```go
// object/object.go

const (
// [...]
    ERROR_OBJ = "ERROR"
)

type Error struct {
    Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
```
we can write a helper function to create an error

```go
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
```

in cases where we didnt handled errors we can use this and return error

```go
func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	// ...
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}
```

also to be able to stop the execution when error hanppens we do need to do some changers
to the `evalProgram` and `evalBlockStatement`


```go
func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
```

as in the above example we check if the result is error and just return the error value


There’s still one last thing we need to do. We need to check for errors whenever we 
call `Eval` inside of `Eval`, in order to stop errors from being passed around and then 
bubbling up far away from their origin:


```go
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// ...
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	// ...
	case *ast.ReturnStatement:
		// we eval the expression associated with the return statement.
		// if the last eval expression associated with return statement. we then
		// wrap the result of this call to `Eval` in our new `object.ReturnValue` so we can
		// keep track on this
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}

	// ...
}
```


