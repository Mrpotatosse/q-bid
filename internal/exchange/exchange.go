package exchange

import "fmt"

type Pair struct {
	Ask string
	Bid string
}

type Triple struct {
	Start Pair
	End   Pair
}

type Exchange interface {
	GetAllPairs(currencies ...string) []*Pair
	GetAllTriples(currencies ...string) []*Triple
}

func (pair Pair) String() string {
	return fmt.Sprintf("%s%s", pair.Ask, pair.Bid)
}
