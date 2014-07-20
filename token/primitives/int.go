package primitives

import (
	"fmt"
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ConstantInt struct {
	value int
}

func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{
		value: value,
	}
}

func (p *ConstantInt) Clone() token.Token {
	return &ConstantInt{
		value: p.value,
	}
}

func (p *ConstantInt) Fuzz(r rand.Rand) {
	// do nothing
}

func (p *ConstantInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *ConstantInt) Parse(pars *token.InternalParser, cur int) (int, []error) {
	v := strconv.Itoa(p.value)
	vLen := len(v)

	nextIndex := vLen + cur

	if nextIndex > pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("Expected %q but got early EOF", v),
			Type:    token.ParseErrorUnexpectedEOF,
		}}
	}

	if got := pars.Data[cur:nextIndex]; v != got {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("Expected %q but got %q", v, got),
			Type:    token.ParseErrorUnexpectedData,
		}}
	}

	return nextIndex, nil
}

func (p *ConstantInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (p *ConstantInt) Permutations() int {
	return 1
}

func (p *ConstantInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *ConstantInt) String() string {
	return strconv.Itoa(p.value)
}

type RandomInt struct {
	value int
}

func NewRandomInt() *RandomInt {
	return &RandomInt{
		value: 0,
	}
}

func (p *RandomInt) Clone() token.Token {
	return &RandomInt{
		value: p.value,
	}
}

func (p *RandomInt) Fuzz(r rand.Rand) {
	p.value = r.Int()
}

func (p *RandomInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *RandomInt) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (p *RandomInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// TODO this could be done MUCH better
	p.value = 0

	return nil
}

func (p *RandomInt) Permutations() int {
	return 1 // TODO maybe this should be like RangeInt
}

func (p *RandomInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *RandomInt) String() string {
	return strconv.Itoa(p.value)
}

type RangeInt struct {
	from int
	to   int

	value int
}

func NewRangeInt(from, to int) *RangeInt {
	if from > to {
		panic("TODO implement that From can be bigger than To")
	}

	return &RangeInt{
		from:  from,
		to:    to,
		value: from,
	}
}

func (p *RangeInt) From() int {
	return p.from
}

func (p *RangeInt) To() int {
	return p.to
}

// Token interface methods

func (p *RangeInt) Clone() token.Token {
	return &RangeInt{
		from:  p.from,
		to:    p.to,
		value: p.value,
	}
}

func (p *RangeInt) Fuzz(r rand.Rand) {
	i := r.Intn(p.Permutations())

	p.permutation(i)
}

func (p *RangeInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *RangeInt) Parse(pars *token.InternalParser, cur int) (int, []error) {
	if cur == pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("Expected integer in range %d-%d but got early EOF", p.from, p.to),
			Type:    token.ParseErrorUnexpectedEOF,
		}}
	}

	i := cur
	v := ""

	for {
		c := pars.Data[i]

		if c < '0' || c > '9' {
			break
		}

		v += string(c)

		if ci, _ := strconv.Atoi(v); ci < p.from || ci > p.to {
			v = v[:len(v)-1]

			break
		}

		i++

		if i == pars.DataLen {
			break
		}
	}

	i--

	ci, _ := strconv.Atoi(v)

	if v == "" || (ci < p.from || ci > p.to) {
		// is the first character already invalid
		if i < cur {
			i = cur
		}

		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("Expected integer in range %d-%d but got %q", p.from, p.to, pars.Data[cur:i]),
			Type:    token.ParseErrorUnexpectedData,
		}}
	}

	p.value = ci

	return i + 1, nil
}

func (p *RangeInt) permutation(i int) {
	p.value = p.from + i
}

func (p *RangeInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	p.permutation(i - 1)

	return nil
}

func (p *RangeInt) Permutations() int {
	return p.to - p.from + 1
}

func (p *RangeInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *RangeInt) String() string {
	return strconv.Itoa(p.value)
}
