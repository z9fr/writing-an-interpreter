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

above is a sample code written in monkey. here we can see 3 statements
```bash
let <identifier> = <expression>;
```
let statements in this language has two changing parts: identifier and an expression

- identifier: `x`,`y` and `add`
- expressions: `10`,`15` and `function literal`