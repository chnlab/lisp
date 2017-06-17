
package platform

import (
    "github.com/chnlab/lisp"
    "time"
)

func Sleep(t []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    if len(t) != 1 {
        return lisp.None, lisp.ErrParaNum
    }

    if t[0].Kind != lisp.Int {
        return lisp.None, lisp.ErrFitType
    }
    time.Sleep(time.Duration(t[0].Text.(int64)) * time.Millisecond)
    return lisp.None, nil
}

