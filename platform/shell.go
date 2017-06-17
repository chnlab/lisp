
package platform

import (
    "github.com/chnlab/lisp"
    "os/exec"

    "strings"
)

func TokenToString(e lisp.Token, p *lisp.Lisp) string {
    switch e.Kind {
    case lisp.Int:
        fallthrough
    case lisp.Float:
        fallthrough
    case lisp.String:
        return e.String()

    case lisp.Label:
        ret, err := p.Exec(e)
        if err != nil {
            return string(e.Text.(lisp.Name))
        }
        return TokenToString(ret, p)

    case lisp.List:
        ret, err := p.Exec(e)
        if err != nil {
            return "E:" + err.Error()
        }
        return TokenToString(ret, p)

    default:
        return "UNIMPL"
    }
}

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

