package lisp

func Scan(s []byte) (list []Token, err error) {
	scanner := pattern.NewScanner(s, true)
	list = make([]Token, 0, 100)
	for {
		a, b, c := scanner.Scan()
		if c != nil {
			break
		}
		switch b {
		case 1:
			list = append(list, Token{Operator, a})
		case 2:
			list = append(list, Token{Int, a})
		case 3:
			list = append(list, Token{Float, a})
		case 4:
			list = append(list, Token{Int, a})
		case 5:
			list = append(list, Token{String, a})
		case 6:
			list = append(list, Token{Label, a})
		}
	}
	if !scanner.Over() {
		err = ErrNotOver
	}
	return
}

func Tree(tkn []Token) ([]Token, error) {
	var f Token
	var s int
	if len(tkn) == 0 {
		return nil, nil
	}
	if tkn[0].Kind == Operator {
		var t bool
		switch tkn[0].Text.(byte) {
		case '(':
			t = true
		case '[':
			t = false
		default:
			return nil, ErrUnquote
		}
		i, j, l := 1, 1, len(tkn)
		for i < l && j > 0 {
			if tkn[i].Kind == Operator {
				switch tkn[i].Text.(byte) {
				case '(', '[':
					j++
				case ')':
					j--
				}
			}
			i++
		}
		if j <= 0 {
			fold, err := Tree(tkn[1 : i-1])
			if err != nil {
				return nil, err
			}
			if t {
				f = Token{Text: fold, Kind: List}
			} else {
				f = Token{Text: fold, Kind: Fold}
			}
			s = i
		} else {
			return nil, ErrUnquote
		}
	} else {
		f = tkn[0]
		s = 1
	}
	rest, err := Tree(tkn[s:])
	if err != nil {
		return nil, err
	}
	ans := make([]Token, 1+len(rest))
	ans[0] = f
	copy(ans[1:], rest)
	return ans, nil
}

func Collect(c map[Name]bool, t *Token) {
	switch t.Kind {
	case List:
		for _, u := range t.Text.([]Token) {
			Collect(c, &u)
		}
	case Label:
		c[t.Text.(Name)] = true
	}
}

func Hard(tkn Token) Token {
	switch tkn.Kind {
	case List:
		l := tkn.Text.([]Token)
		x := make([]Token, len(l))
		for i, t := range l {
			x[i] = Hard(t)
		}
		return Token{List, x}
	case Label:
		t, ok := Global.env[tkn.Text.(Name)]
		if ok {
			return t
		}
	case Macro:
		m := tkn.Text.(*Hong)
		l := m.Text
		x := make([]Token, len(l))
		for i, t := range l {
			x[i] = Hard(t)
		}
		return Token{Macro, &Hong{m.Para, x, m.Real}}
	case Front:
		f := tkn.Text.(*Lfac)
		l := f.Text
		x := make([]Token, len(l))
		for i, t := range l {
			x[i] = Hard(t)
		}
		return Token{Front, &Lfac{f.Para, x, f.Make}}
	}
	return tkn
}

