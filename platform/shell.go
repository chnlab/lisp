
package platform

import (
    "github.com/chnlab/lisp"
    "os/exec"

    "strings"
)

func Shell(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {

    cmd := ""
    for i, e := range ls {
        if i != 0 {
            cmd = cmd + " "
        }
        cmd = cmd + TokenToString(e, p)
    }

    args := strings.Fields(cmd)
    out, err := exec.Command(args[0], args[1:]...).Output()
    return lisp.Token{lisp.String, string(out)}, err
}

