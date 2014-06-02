package filters

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type FuncFilter struct {
	fuzzFunc   func(r rand.Rand, tok token.Token) interface{}
	stringFunc func(state interface{}, tok token.Token) string
	state      interface{}
	token      token.Token
}

func NewFuncExpression(
	tok token.Token,
	fuzzFunc func(r rand.Rand, tok token.Token) interface{},
	stringFunc func(state interface{}, tok token.Token) string,
) *FuncFilter {
	return &FuncFilter{
		fuzzFunc:   fuzzFunc,
		stringFunc: stringFunc,
		state:      nil,
		token:      tok,
	}
}

func (f *FuncFilter) Clone() token.Token {
	return &FuncFilter{
		fuzzFunc:   f.fuzzFunc,
		stringFunc: f.stringFunc,
		state:      f.state,
		token:      f.token,
	}
}

func (f *FuncFilter) Fuzz(r rand.Rand) {
	f.state = f.fuzzFunc(r, f.token)
}

func (f *FuncFilter) String() string {
	return f.stringFunc(f.state, f.token)
}
