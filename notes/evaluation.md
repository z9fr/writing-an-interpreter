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
