Parsing Expressions
----

Parsing expressions does has some challengers.

### Operator Precedence

lets say we want to parse the below arithmetic expression:

```js
5 * 5 + 10
```

after parsing this what we want is a AST that represents the expression like below

```js
((5 * 5) + 10)
```

> This is to say 5 * 5 needs to be 'deeper' in the AST and evaluated earlier before addition
  in order to produce AST that looks like this, parser needs to know about operator precedences
  where * is higher than +

but there are other cases where this is important

---
```js 
5 * (5 + 10) ```

here the parenthesis group together the `5 + 10` and give `precedence bump` 
the addition now has to be evaluated before the multiplication. 

---
```js
-5 - 10
```
when expressions tokens of the same type appear in multiple possitions.
here the `-` operator appears at the beginning of the expression as a prefix operator,
then there is an infix operator in the middle. 

below is a variation of the same challenge

```js
5 * (add(2, 3) + 10)
```

---

Expressions in Monkey
---

in monkey programming language everything besides let and return statement is an expression.
these can come in different varities

- prefix operators:
```js
-5 
!true
!false
```
- infix operators (or "binary operators"):
```js
5 + 5
5 - 5
5 / 5
5 * 5
```
- basic arithmetic operators and comparison operators:
```js
foo == bar
foo != bar
foo < bar
foo > bar
```

- parentheses to group expressions and influence order of evaluation:
```js
5 * (5 + 5)
((5 + 5) * 5) * 5
```
- call expressions
```js
add(2, 3)
add(add(2, 3), add(5, 10))
max(5, add(5, (5 * 5)))
```
- Identifiers are expressions too:
```js
foo * bar / foobar
add(foo, bar)
```
- Functions in Monkey are first-class citizens, function literals are expressions:
```js
let add = fn(x, y) { return x + y };

// here we use a function literal in place of an identifier:
fn(x, y) { return x + y }(5, 5)
(fn(x) { return x }(5) + 10 ) * 10
```
- if expressions
```js
let result = if (10 > 5) { true } else { false };
result // => true
```
---

# Top Down Operator Precedence (or: Pratt Parsing)

References 
- https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
- https://crockford.com/javascript/tdop/tdop.html

> The parsing approach described by all three, which is called Top Down 
Operator Precedence Parsing, or Pratt parsing, was invented as an alternative
to parsers based on context-free grammars and the Backus-Naur-Form.

#### This is also the main difference:

insted of associateing parsing functions with rules, Praat associates these functions
with single token type. 

this main part of this idea is that each token type can have two parsing functions 
associated with it. depending on the token position - infix or prefix

---


### Terminology

- A *prefix operator* is an operator "in front of" its operand eg: `--5`
    - here the operator is `--`, the operand is the int `5` operator is the prefix position

- A *postfix operator* is an operator "after" its operand eg: `foobar++`
    - operator is `++`, the operand is `foobar`. the operator is in the postfix possition

- A *infix operators* sits between its operands, eg: `5 * 8`
    - the `*` operator sits in the infix position between two int `5` and `8`. 
    - infix operators appear in *binary expressions* - where the operator has two operands

- *operator precedence* (order of operations), which priority do different operators have. 
eg : `5 + 5 * 10`
    - the result for this is `55` not `100`. that's because `*` operator has higher precedence (rank)
    - its more important than the `+` operator. 
    - `*` gets eval before the `+` operator


> These are all basic terms: prefix, postfix, infix operator and precedence. 
But it’s important that we keep these simple definitions in mind later on,
where we’ll use these terms in other places.
