Infix math expression evaluator.  As of writing, only supports integer math.
Defined operators are + for addition, - for subtraction, \* for multiplication,
/ for division, and ^ for exponent.  Computes integer math in the normal
computer fashion; i.e., 10 / 3 is 3.

To install: `go get github.com/jordanorelli/shunt`

If you have your Go environment set properly, you can now just say `shunt` and
you'll get an infix expression evaluator that reads statements from stdin,
evaluates them, and prints the result on stdout.  You can also pipe your input
into the program, so `echo 10 * (5-1) | shunt` will print `40` and quit.
