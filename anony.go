package lisp

import (
	"math/rand"
	"time"
)


func init() {
	rand.Seed(time.Now().UnixNano())

	Add("macro", func(t []Token, p *Lisp) (ans Token, err error) {
		var (
			a, b, c Token
			x, y    []Name
		)

		switch len(t) {
		case 2:
			a, b = t[0], t[1]

		case 3:
			a, c, b = t[0], t[1], t[2]
			if c.Kind != List {
				return None, ErrFitType
			}
			t = c.Text.([]Token)
			y = make([]Name, len(t))
			for i, u := range t {
				if u.Kind != Label {
					return None, ErrNotName
				}
				y[i] = u.Text.(Name)
			}

		default:
			return None, ErrParaNum
		}

		if a.Kind != List || b.Kind != List {
			return None, ErrFitType
		}

		t = a.Text.([]Token)
		x = make([]Name, len(t))
		for i, u := range t {
			if u.Kind != Label {
				return None, ErrNotName
			}
			x[i] = u.Text.(Name)
		}

		ans = Token{Macro, &Hong{x, b.Text.([]Token), y}}
		return ans, nil
	})

	Add("lambda", func(t []Token, p *Lisp) (ans Token, err error) {
		if len(t) < 2 {
			return None, ErrParaNum
		}

		a, b := t[0], t[1]

        if len(t) > 2 {
            eachList := make([]Token, len(t))
            eachList[0] = Token{Label, Name("each")}
            eachList = append(eachList, t[1:]...)

            b = Token{List, eachList}
        }

		if a.Kind != List {
			return None, ErrFitType
		}

		if b.Kind != List {
			return None, ErrFitType
		}

		t = a.Text.([]Token)
		x := make([]Name, 0, len(t))
		for _, i := range t {
			if i.Kind != Label {
				return None, ErrNotName
			}
			x = append(x, i.Text.(Name))
		}

		ans = Token{Front, &Lfac{x, b.Text.([]Token), p}}
		return ans, nil
	})
}

func replace(tkn Token, lst map[Name]Token) Token {
	switch tkn.Kind {
	case List:
		l := tkn.Text.([]Token)
		x := make([]Token, len(l))
		for i, t := range l {
			x[i] = replace(t, lst)
		}
		return Token{List, x}
	case Label:
		t, ok := lst[tkn.Text.(Name)]
		if ok {
			return t
		}
	}
	return tkn
}

func genName() Name {
	u := [16]byte{'_'}
	for i := 1; i < 16; i++ {
		switch x := rand.Uint32() % 63; {
		case x < 26:
			u[i] = byte(x + 'A')
		case x < 52:
			u[i] = byte(x + 'a' - 26)
		case x < 62:
			u[i] = byte(x + '0' - 52)
		default:
			u[i] = '_'
		}
	}
	return Name(string(u[:]))
}

func evalMacro(mt Token, ls []Token, p *Lisp) (Token, error) {
    body := mt.Text.(*Hong)

    if len(ls) != len(body.Para)+1 {
        return None, ErrParaNum
    }

    xp := map[Name]Token{}
    for i, t := range ls[1:] {
        xp[body.Para[i]] = t
    }

    if body.Real == nil {
        for i, t := range body.Para {
            xp[t] = ls[1+i]
        }
    } else {
        cp := map[Name]bool{}
        for i, t := range ls[1:] {
            xp[body.Para[i]] = t
            Collect(cp, &t)
        }

        for _, t := range body.Real {
            var i Name
            for {
                i = genName()
                _, ok := cp[i]
                if !ok {
                    break
                }
            }
            xp[t] = Token{Label, i}
        }
    }
    return p.Exec(replace(Token{List, body.Text}, xp))
}

