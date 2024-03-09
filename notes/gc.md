# Who's taking the trash out?

let's take a example 

```js
let counter = fn(x) {
  if (x > 100) {
    return true;
  } else {
    let foobar = 9999;
    counter(x + 1);
  }
};

counter(0);
```

when eval this code first if-else expression `x > 100` the value is not truthy becuase of this 
the alternative of the if-else expression gets eval. 

in the else block int `9999` get bund to variable `foobar`, this never referenced again.
then `x + 1` is eval the result of that call to `Eval` is then passed to another call to `counter`

And then it all starts again until `x > 100` eval to `TRUE`


In each call to `counter` lot of objects are allocated.


in terms of our implementation `Eval` function and our object system: each evaluation of counter's
body result in lot of `object.Integer`  being allocated and instantiated.

The unused `9999` int and the result of `x + 1` are obvious. but even the literals `100` and `1`
produce new `object.Integers` every time body of `counter` is evaluated.

> this result in we have aorund 400 allocated `object.Integers`


---

Our objects are stored in memory. more objects we use the more memory we need. and eventho 
the number of objects the example is tiny compaired to other programs memory its not infinite.

> with each call to `counter` the memory usage will high and at onepoint os will kill it. 

but if we monitor memory usage we notice that it doesnt take alot. 

> Go’s garbage collector (GC) is the reason why we don’t run out of memory. It manages memory for us

---

## What does GC need to do ?

- Keep track of object allocations and references to objects. 
- Make enough memory availible for future object allocations
- Give memory back when its not needed anymore


There are ways to accomplish all of the above. involving different algorithms and implementations.

For example there's "mark and sweep" algorithm. To implement it one has to decide whether
the GC will be a generational GC or not, or whether it's stop-the-world GC or concurrent GC, 
or how it's organizing memory and handling memory fragmentation having decided all of that 
an effiecent implementation is still a lot of hard work

FYI: https://www.geeksforgeeks.org/mark-and-sweep-garbage-collection-algorithm/


---

### Adding GC to this implementation ?

we need to dissable Go's GC and find a way to take over all its duties. It’s a huge 
undertaking since we would also have to take care of allocating and freeing memory ourselves 
- in a language that per default prohibits exactly that.

