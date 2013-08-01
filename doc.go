/*

ebnf2y converts EBNF grammars into yacc compatible skeleton .y files.

Installation:

	$ go get github.com/cznic/ebnf2y

Usage:

	ebnf2y [options] [file]

Options:

	-ie number	Inline eligible EBNF productions:
			  0: none (default)
			  1: used once
			  2: all (cannot be used with -m)
	-iy number	Inline eligible BNF (.y) productions:
			  0: none (default)
			  1: used once
			  2: all (cannot be used with -m)
	-m		Magic: Attempt to to minimize yacc conflicts,
			  by finding a minimum of wr*RR+ws*SR
			  (wr*reducereduce+ws*shiftreduce conflicts).
	-M		Like -m and write report to stderr.
	-o name		Output file name. Stdout if left blank (default).
	-oe name	Output pretty printed EBNF to <name>.
	-p string	Prefix for token names, eg. "_". Default blank.
	-start name	Select start production name. Default is "SourceFile".
	-wr		Weight of reduce/reduce conflicts for -m.
	-ws		Weight of shift/reduce conflicts for -m.

File:
	A named EBNF file. If no non-opt args are given, ebnf2y reads stdin.

Notation

The EBNF flavor is the one used by the Go language specification[1]:
_______________________________________________________________________________

The syntax is specified using Extended Backus-Naur Form (EBNF):

	Production  = production_name "=" [ Expression ] "." .
	Expression  = Alternative { "|" Alternative } .
	Alternative = Term { Term } .
	Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
	Group       = "(" Expression ")" .
	Option      = "[" Expression "]" .
	Repetition  = "{" Expression "}" .

Productions are expressions constructed from terms and the following operators,
in increasing precedence:

	|   alternation
	()  grouping
	[]  option (0 or 1 times)
	{}  repetition (0 to n times)

Lower-case production names are used to identify lexical tokens. Non-terminals
are in CamelCase. Lexical tokens are enclosed in double quotes "" or back
quotes ``.

The form a … b represents the set of characters from a through b as
alternatives. The horizontal ellipsis … is also used elsewhere in the spec to
informally denote various enumerations or code snippets that are not further
specified. The character … (as opposed to the three characters ...) is not a
token of the Go language.

_______________________________________________________________________________

Generated code

Many non trivial EBNF grammars will produce shift/reduce conflicts. These must
be resolved manually. The generated yacc (*.y) file is a skeleton parser,
intended only as a starting point of a real parser. However, for some simple
grammars the automatically generated parser might be (almost) useful as it is.

This example EBNF[2]:

	float		= . // http://golang.org/ref/spec#float_lit
	identifier	= . // ASCII letters, digits, "_". No front digit.
	imaginary	= . // http://golang.org/ref/spec#imaginary_lit
	integer		= . // http://golang.org/ref/spec#int_lit
	str		= . // http://golang.org/ref/spec#string_lit
	boolean		= "true" | "false" .

	andnot 	= "&^" .
	lsh 	= "<<" .
	rsh 	= ">>" .

	Expression = Term  { ( "^" | "|" | "-" | "+" ) Term } .
	ExpressionList = Expression { "," Expression } .
	Factor = [ "^" | "!" | "-" | "+" ] Operand .
	Literal = boolean
		| float
		| QualifiedIdent
		| imaginary
		| integer
		| str .
	Term = Factor { ( andnot | "&" | lsh  | rsh | "%" | "/" | "*" ) Factor } .
	Operand = Literal
	        | QualifiedIdent "(" [ ExpressionList ] ")"
	        | "(" Expression ")" .
	QualifiedIdent = identifier [ "." identifier ] .

produces a yacc file[3], seen at the link without any edits in its original
form as generated by ebnf2y. In conjunction with a simple lexer[4] (any other
tokenizer for the particular grammar can be used as well), a simple demo
program can be compiled. Running it without arguments prints:

	$ ./demo
	AST for 'fmt.Printf("%d\012", -1 + 2.3*^3i | 4e2)'.
	[]main.Expression{
	. []main.Term{
	. . []main.Factor{
	. . . []main.Operand{
	. . . . []main.QualifiedIdent{
	. . . . . "fmt",
	. . . . . []main.QualifiedIdent1{
	. . . . . . ".",
	. . . . . . "Printf"
	. . . . . }
	. . . . },
	. . . . "(",
	. . . . []main.ExpressionList{
	. . . . . []main.Expression{
	. . . . . . []main.Term{
	. . . . . . . []main.Factor{
	. . . . . . . . "%d\n"
	. . . . . . . },
	. . . . . . },
	. . . . . },
	. . . . . []main.ExpressionList1{
	. . . . . . ",",
	. . . . . . []main.Expression{
	. . . . . . . []main.Term{
	. . . . . . . . []main.Factor{
	. . . . . . . . . "-",
	. . . . . . . . . 1
	. . . . . . . . },
	. . . . . . . },
	. . . . . . . []main.Expression1{
	. . . . . . . . "+",
	. . . . . . . . []main.Term{
	. . . . . . . . . []main.Factor{
	. . . . . . . . . . 2.3
	. . . . . . . . . },
	. . . . . . . . . []main.Term1{
	. . . . . . . . . . "*",
	. . . . . . . . . . []main.Factor{
	. . . . . . . . . . . "^",
	. . . . . . . . . . . (0+3i)
	. . . . . . . . . . }
	. . . . . . . . . }
	. . . . . . . . },
	. . . . . . . . "|",
	. . . . . . . . []main.Term{
	. . . . . . . . . []main.Factor{
	. . . . . . . . . . 400
	. . . . . . . . . },
	. . . . . . . . }
	. . . . . . . }
	. . . . . . }
	. . . . . }
	. . . . },
	. . . . ")"
	. . . }
	. . },
	. },
	}

	$

Note: In the above output, nil items have been removed from the AST dump before
printing.

Demo

Prerequisites: ebnf2y and golex[4] must be installed.

The demo program can be built by:

	$ make demo

The produced binary 'demo/demo' accepts an expression as the value of its
'-e' flag. Default is `fmt.Printf("%d\012", -1 + 2*^3 | 4)`.

References

Links fom the above godocs:

  [1]: http://golang.org/ref/spec#Notation
  [2]: http://github.com/cznic/ebnf2y/blob/master/demo/demo.ebnf
  [3]: http://github.com/cznic/ebnf2y/blob/master/demo/demo.y
  [3]: http://github.com/cznic/ebnf2y/blob/master/demo/demo.l
  [4]: http://github.com/cznic/golex

*/
package main
