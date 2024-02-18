Function Literals
-----------------

example of functions in monkey looks like below: we define which params they have and 
what the function do
```js
fn(x, y) {
  return x + y;
}
```

funcion starts with `fn` keyword followed by list of params. followed by block statement.

```js
fn <parameters> <block statement>

// params
(<parameter one>, <parameter two>, <parameter three>, ...)

// params are just a list of identifiers that are comma-seperated
// they can also be empty
```

when parsing function literel we have the below impl. we check if the token is `(` 
and we start parsing function params after its done.

we look for block we check if peek starts with `{` and start parsing the block
statement
```go 
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

```

in the function param parser we try to return array of `*ast.Identifier`. 
we first check if the token is `(` and we build `Identifier` struct. until we have a `,`
we keep on parsing the `Identifier` so we end up building an array of `identifiers`
and return

```go
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

```

-----

Call Expressions
----------------

call expressions means basically is parsing the calling of function: call expressions.
```js

add(2, 3)
add(2 + 2, 3 * 3 * 3) // arguments can be expressions too
callsFunction(2, 3, fn(x, y) { x + y; }); // function literals can be arguments

// the structure 
<expression>(<comma separated expressions>)
```


