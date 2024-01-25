# Lexing

## Lexical Analysis

- in order to work with sourcecode, we need to convert it to more accessable form.
- we need to represent our sourcecode in a form that is easier to work with,
  we are going to change the representation of our source code two times before we evaluate it

  |Source code| -> |tokens| -> |AST|

first transformation, from source code to tokens is called "Lexical Analysis", or "lexting" for short
its done by lexter also called tokenizer or scanner – some use one word or the other
to denote subtle differences in behaviour, which we can ignore in this book).

Example:

this is the input given to lexer

```
let  = 5 + 5;
```

and we get a output looks like below

```go
[
    LET,
    IDENTIFIER("x"),
    EQUAL_SIGN,
    INTEGER(5),
    PLUS_SIGN,
    INTEGER(5),
    SEMICOLON
]
```

A production-ready lexer might also attach the line number, column number and filename to
a token. Why? For example, to later output more useful error messages in the parsing stage.
Instead of "error: expected semicolon token" it can output:

```
error: expected semicolon token. line 42, column 23, program.monkey
```

## Defining Our Tokens

```
let five = 5;
let ten = 10;

let add = fn(x, y) {
    x + y;
};

let result = add(five, ten);
```

Let’s break this down: which types of tokens does this example contain? First of all, there are
the numbers like 5 and 10. These are pretty obvious. Then we have the variable names x, y,
add and result. And then there are also these parts of the language that are not numbers, just
words, but no variable names either, like let and fn. Of course, there are also a lot of special
characters: (, ), {, }, =, ,, ;.

- numbers are just int and we are going to treat them as such
- we call the variable names "identifiers"
- the otherones looks like identifiers but arnt really identifiers, since they are part of lang
  we call them keywords. we wont group group them togther since it should make a difference in parsing stage
  wether we encounter a let of fn
- same goes for special characters, we treat them seperately since it is a big difference whether or not we have a ( or a )

in our token struct we need type so we can distinguish between “integers” and “right bracket”
also it needs a field that holds the literal value of the token, so we can reuse ti later

# The Lexer
