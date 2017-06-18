
package platform

import (
    "github.com/chnlab/lisp"
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

func Lines(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    lines := strings.Split(TokenToString(ls[0], p), "\n")

    ret := make([]lisp.Token, len(lines))
    for i, line := range lines {
        ret[i] = lisp.Token{lisp.String, line}
    }
    return lisp.Token{lisp.List, ret}, nil
}

func Words(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    words := strings.Fields(TokenToString(ls[0], p))

    ret := make([]lisp.Token, len(words))
    for i, word := range words {
        ret[i] = lisp.Token{lisp.String, word}
    }
    return lisp.Token{lisp.List, ret}, nil
}

func Split(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    if len(ls) < 2 {
        return lisp.None, lisp.ErrParaNum
    }

    str := TokenToString(ls[0], p)
    sep := TokenToString(ls[1], p)
    words := strings.Split(str, sep)

    ret := make([]lisp.Token, len(words))
    for i, word := range words {
        ret[i] = lisp.Token{lisp.String, word}
    }
    return lisp.Token{lisp.List, ret}, nil
}

func Join(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    if len(ls) < 2 {
        return lisp.None, lisp.ErrParaNum
    }

    var err error
    element := ls[0]
    if element.Kind == lisp.Label || element.Kind == lisp.List {
        element, err = p.Exec(element)
        if err != nil {
            return lisp.None, err
        }
    }

    if element.Kind != lisp.List {
        return lisp.None, lisp.ErrFitType
    }

    sep := TokenToString(ls[1], p)
    ts := element.Text.([]lisp.Token)
    words := make([]string, len(ts))

    for i, t := range ts {
        words[i] = TokenToString(t, p)
    }

    return lisp.Token{lisp.String, strings.Join(words, sep)}, nil
}

func Length(ls []lisp.Token, p *lisp.Lisp) (lisp.Token, error) {
    var err error
    element := ls[0]

    if element.Kind == lisp.Label || element.Kind == lisp.List {
        element, err = p.Exec(element)
        if err != nil {
            return lisp.None, err
        }
    }

    if element.Kind != lisp.String {
        return lisp.None, lisp.ErrFitType
    }

    size := len(element.Text.(string))
    return lisp.Token{lisp.Int, size}, nil
}

