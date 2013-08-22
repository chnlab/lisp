package lisp

import "fmt"

func init() {
	Global.Add("raise", func(t []Token, p *Lisp) (Token, error) {
		if len(t) != 1 {
			return None, ErrParaNum
		}
		ans, err := p.Exec(t[0])
		if err != nil {
			return None, err
		}
		if ans.Kind != String {
			return None, ErrFitType
		}
		return None, fmt.Errorf(ans.Text.(string))
	})
	Global.Add("catch", func(t []Token, p *Lisp) (Token, error) {
		if len(t) != 1 {
			return None, ErrParaNum
		}
		_, err := p.Exec(t[0])
		if err != nil {
			return Token{String, fmt.Sprint(err)}, nil
		}
		return None, nil
	})
}