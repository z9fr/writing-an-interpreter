# Parsing

> software component which takes input nd builds a data structure.

in interpretters and compilers data structure used for internal representation
of the source called is called `syntax tree` or an `abstract syntax tree (AST)`

"abstract" is based on the fact that certain details on sourcecode is removed in AST
eg: “Semicolons, newlines, whitespace, comments, braces, bracket and parentheses ”

> There is no universal AST format can be used by every parser. 
their implementations are pretty same, concept is the same
but different in details. the concrete implementation depends on 
programing langauge being parsed

Example

```js
if (3 * 5 > 10) {
  return "hello";
} else {
  return  "goodbye";
}
```

if we use `MagicLexer` a `MagicParser` and the AST is build out of js objects,
the parsing step results look like below
```js

> var input = 'if (3 * 5 > 10) { return "hello"; } else { return "goodbye"; }';
> var tokens = MagicLexer.parse(input);
> MagicParser.parse(tokens);
{
  type: "if-statement",
  condition: {
    type: "operator-expression",
    operator: ">",
    left: {
      type: "operator-expression",
      operator: "*",
      left: { type: "integer-literal", value: 3 },
      right: { type: "integer-literal", value: 5 }
    },
    right: { type: "integer-literal", value: 10 }
  },
  consequence: {
    type: "return-statement",
    returnValue: { type: "string-literal", value: "hello" }
  },
  alternative: {
    type: "return-statement",
    returnValue: { type: "string-literal", value: "goodbye" }
  }
```
----

## Parser generators

there are tools like

- yacc
- bison
- ANTLR 

parser generator tools. when we feed formal description of the language. 
it will produce parsers as output

---

## Writing a parser

there are 2 main strategies when parsing a language:

- top down parsing
  - recursive descent parsing
  - early parsing
  - predictive parsing
  - etc (variations of top down parsing)
- bottom up parsing


---

our implmenetation is [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser). in particular
it's a "top down operator precedence" parser. also known as "Praat parser"


> that the difference between top down and bottom up parsers is that 
the former starts with constructing root node of the AST and then 
descends while the latter does it the other way around

---

## Parser first steps: parsing let statements

```js
let x = 10;
let y = 15;

let add = fn(a, b) {
  return a + b;
};
```


> what does it mean to parse let statement correctly ?
> - parser produce an AST that accurately represents the information in orignial let statement


above is a sample code written in monkey. here we can see 3 statements
```bash
let <identifier> = <expression>;
```
let statements in this language has two changing parts: identifier and an expression

- identifier: `x`,`y` and `add`
- expressions: `10`,`15` and `function literal`


Difference between statements and expression
====

- expressions produce values
- statements dont produce values. 


---

```go
// ast/ast.go

package ast

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}
```

- we need two different types of nodes: expression and statements.


> we have 3 interfaces called `Node`, `Statement` and `Expression`. every node in our
  AST has to implment the `Node` interface.

> it has to provide a `TokenLiteral()` method that returns the literal value of
  token its associated with. ( used for debugging and testing)


- The AST we construct will consists solely of Nodes that are connected to each other. 
it's a tree after all.
- Some will implement `Statement` and some will impl  `Expression` interfaces
  - These interfaces only has dummy methords called `statementNode` and `expressionNode`
  - These methords are not necessary  but help to compiler.
  - Can use to throw errors when use a `Statement` where `Expression` should use or vise versa

---

## Parsing let statements

```go
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
```

- We can star with `parseLetStatement` methord. it creates an `*ast.LetStatement` with the token it sitting on `token.LET`.
- Then it advances the tokens while making assertions about the next token. using the calls to `expectPeek`
- First expect the `token.IDENT` token. which then it uses to construct an `*ast.Identifier` node
- Then expects an equal sign and finally jumps over the expression following the equal sign untill encounters a semicolon

