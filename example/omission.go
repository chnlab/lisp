package main

import "github.com/chnlab/lisp"

func main() {
	console := lisp.NewLisp()
	console.Eval([]byte(`
		(each
			(define
				(g x y)
				(cons x y)
			)
			(define
				G
				(omission g)
			)
			(println
				(G 3)
			)
			(println
				(G 3 2)
			)
			(println
				(G 3 2 1)
			)
		)
	`))
}
