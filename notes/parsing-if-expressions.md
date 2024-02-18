If Expressions
--------------

Sample usage of if-else 

```js
if (x > y) {
  return x;
} else {
  return y;
}

// the `else` is optional and can be ignored

if(x > y) {
    return x;
}
```

in monkey if-else-conditionals are expressions. means they produce a value 

```js
let foobar = if (x > y) { x } else { y };
```

`if (<condition>) <consequence> else <alternative>`

the impl for the `parseIfExpression` is below. we first parse and get the condition. 
and check untill we enconter `{` when we do will will parse the statement as a block statement


after when we enconter `token.ELSE` we also do the same and parse it as a block.
note that we are allowing `else` to be optional when else is not avaible we just ignore
```go
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		// allow optional `else` but does not add parser error when `token.ELSE`
		// is missing.
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}
```

impl for the `parseBlockStatement` here we return `*ast.BlockStatement` untill we encouner
`}` or `token.EOF` we will parse all the statements avaible and append them to staements array
in the `BlockStatement`

```go 
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// calls `parseStatement` untill it enconters either a `}` or `token.EOF`
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}
```
