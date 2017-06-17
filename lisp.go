package lisp

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
)

type Lisp struct {
    dad *Lisp
    env map[Name]Token
}

func NewLisp() *Lisp {
    x := new(Lisp)
    x.env = map[Name]Token{}
    x.dad = Global
    return x
}

func Add(s string, f func([]Token, *Lisp) (Token, error)) {
    Global.env[Name(s)] = Token{Back, Gfac(f)}
}

func (l *Lisp) Exec(f Token) (ans Token, err error) {
    var (
        ls []Token
        ct Token
        ok bool
    )

    switch f.Kind {
    case Fold:
        return Token{List, f.Text.([]Token)}, nil

    case Label:
        nm := f.Text.(Name)
        for ; l != nil; l = l.dad {
            ct, ok = l.env[nm]
            if ok {
                return ct, nil
            }
        }
        return None, ErrNotFind(string(nm))

    case List:
        ls = f.Text.([]Token)
        if len(ls) == 0 {
            return False, nil
        }
        ct = ls[0]

        switch ct.Kind {
        case Label:
            nm := ct.Text.(Name)
            for v := l; v != nil; v = v.dad {
                ct, ok = v.env[nm]
                if ok {
                    break
                }
            }

            if !ok {
                return None, ErrNotFind(string(nm))
            }

        case List:
            ct, err = l.Exec(ct)
            if err != nil {
                return None, err
            }
        }

        switch ct.Kind {

        case Back:
            return ct.Text.(Gfac)(ls[1:], l)

        case Front:
            lp := ct.Text.(*Lfac)
            if len(ls) != len(lp.Para) + 1 {
                return None, ErrParaNum
            }

            q := &Lisp{dad: lp.Make, env: map[Name]Token{}}
            q.env[Name("self")] = ct
            for i, t := range ls[1:] {
                q.env[lp.Para[i]], err = l.Exec(t)
                if err != nil {
                    return None, err
                }
            }
            return q.Exec(Token{List, lp.Text})

        case Macro:
            return evalMacro(ct, ls, l)

        case Chan:

            switch len(ls) {
            case 1:
                u, ok := <-ct.Text.(chan Token)
                if ok {
                    return u, nil
                } else {
                    return None, ErrIsClose
                }

            case 2:
                u, err := l.Exec(ls[1])
                if err != nil {
                    return None, err
                }
                t := func() (s string) {
                    defer func() {
                        e := recover()
                        if e != nil {
                            s = fmt.Sprint(e)
                        }
                    }()
                    ct.Text.(chan Token) <- u
                    return
                }()
                if t != "" {
                    return None, errors.New(t)
                }
                return u, nil

            default:
                return None, ErrParaNum
            }

        default:
            return None, ErrNotFunc(f.String())
        }

    default:
        return f, nil
    }

    return None, nil
}

func (l *Lisp) Eval(s []byte) (Token, error) {
    var (
        a, b []Token
        c, d Token
        e    error
    )

    a, e = Scan(s)
    if e != nil {
        return None, e
    }

    b, e = Tree(a)
    if e != nil {
        return None, e
    }

    for _, c = range b {
        d, e = l.Exec(c)
        if e != nil {
            return None, e
        }
    }
    return d, nil
}

func (l *Lisp) Load(s string) (Token, error) {
    var file *os.File
    var data []byte
    var err error

    file, err = os.Open(s)
    if err != nil {
        file, err = os.Open(s + ".lsp")
        if err != nil {
            return None, err
        }
    }
    defer file.Close()

    data, err = ioutil.ReadAll(file)
    if err != nil {
        return None, err
    }

    buf := bytes.NewBuffer(data)
    one := section{}
    for {
        data, err := buf.ReadBytes('\n')
        if err != nil {
            if err != io.EOF {
                return None, err
            }
            err = one.feed(data)
            break
        }

        err = one.feed(data)
        if err != nil {
            return None, err
        }
    }

    if !one.over() {
        return None, ErrUnquote
    }

    return l.Eval([]byte(one.total))
}

