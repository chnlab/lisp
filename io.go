package lisp

import (
	"bufio"
	"fmt"
	"os"
)

func init() {

	Add("scan", func(t []Token, p *Lisp) (Token, error) {
		if len(t) != 0 {
			return None, ErrParaNum
		}

		buf := bufio.NewReader(os.Stdin)
		one := section{}
		for {
			data, err := buf.ReadBytes('\n')
			if err != nil {
				return None, err
			}
			err = one.feed(data)
			if err != nil {
				return None, err
			}
			if one.over() {
				break
			}
		}
		return p.Eval([]byte(one.total))
	})

	Add("load", func(t []Token, p *Lisp) (Token, error) {
		if len(t) != 1 {
			return None, ErrParaNum
		}

		if t[0].Kind != String {
			return None, ErrFitType
		}

		return p.Load(t[0].Text.(string))
	})


    display := func(t []Token, p *Lisp) (Token, error) {
        var ret Token

		for i, e := range t {
            ret, err := p.Exec(e)
			if err != nil {
				return None, err
			}

            if i != 0 {
                fmt.Print(" ")
            }
            fmt.Print(ret)
        }
        return ret, nil
    }

	Add("print", display)

	Add("println", func(t []Token, p *Lisp) (Token, error) {

        ret, err := display(t, p)
        if err != nil {
            return None, err
        }

        fmt.Print("\n")
        return ret, nil
	})
}

