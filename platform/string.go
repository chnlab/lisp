
package platform

import (
    "github.com/chnlab/lisp"
)

func String(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    str := ""
    for i, e := range ls {
        if i != 0 {
            str = str + " "
        }
        str = str + TokenToString(e, p)
    }
    return lisp.Token{lisp.String, str}, nil
}

