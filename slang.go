// Package slang (short for Sabre lang) provides a tiny LISP dialect built
// using  factilities provided by Sabre. See New() for  initializing Slang
// and using it.
package slang

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/spy16/sabre"
)

const (
	nsSeparator = '/'
	defaultNS   = "user"
)

// New returns a new instance of Slang interpreter.
func New() *Slang {
	sl := &Slang{
		mu:       &sync.RWMutex{},
		bindings: map[nsSymbol]sabre.Value{},
	}

	if err := BindAll(sl); err != nil {
		panic(err)
	}
	sl.checkNS = true

	_ = sl.SwitchNS(sabre.Symbol{Value: defaultNS})
	_ = sl.BindGo("ns", sl.SwitchNS)
	return sl
}

// Slang represents an instance of slang interpreter.
type Slang struct {
	mu        *sync.RWMutex
	currentNS string
	checkNS   bool
	bindings  map[nsSymbol]sabre.Value
}

// Eval evaluates the given value in Slang context.
func (slang *Slang) Eval(v sabre.Value) (sabre.Value, error) {
	return sabre.Eval(slang, v)
}

// ReadEval reads from the given reader and evaluates all the forms
// obtained in Slang context.
func (slang *Slang) ReadEval(r io.Reader) (sabre.Value, error) {
	return sabre.ReadEval(slang, r)
}

// ReadEvalStr reads the source and evaluates it in Slang context.
func (slang *Slang) ReadEvalStr(src string) (sabre.Value, error) {
	return sabre.ReadEvalStr(slang, src)
}

// Bind binds the given name to the given Value into the slang interpreter
// context.
func (slang *Slang) Bind(symbol string, v sabre.Value) error {
	slang.mu.Lock()
	defer slang.mu.Unlock()

	nsSym, err := slang.splitSymbol(symbol)
	if err != nil {
		return err
	}

	if slang.checkNS && nsSym.NS != slang.currentNS {
		return fmt.Errorf("cannot bind outside current namespace")
	}

	slang.bindings[*nsSym] = v
	return nil
}

// Resolve finds the value bound to the given symbol and returns it if
// found in the Slang context and returns it.
func (slang *Slang) Resolve(symbol string) (sabre.Value, error) {
	slang.mu.RLock()
	defer slang.mu.RUnlock()

	if symbol == "ns" {
		symbol = "user/ns"
	}

	nsSym, err := slang.splitSymbol(symbol)
	if err != nil {
		return nil, err
	}

	return slang.resolveAny(symbol, *nsSym, nsSym.WithNS("core"))
}

// BindGo is similar to Bind but handles conversion of Go value 'v' to
// sabre Value type.
func (slang *Slang) BindGo(symbol string, v interface{}) error {
	return slang.Bind(symbol, sabre.ValueOf(v))
}

// SwitchNS changes the current namespace to the string value of given symbol.
func (slang *Slang) SwitchNS(sym sabre.Symbol) error {
	slang.mu.Lock()
	slang.currentNS = sym.String()
	slang.mu.Unlock()

	return slang.Bind("*ns*", sym)
}

// CurrentNS returns the current active namespace.
func (slang *Slang) CurrentNS() string {
	slang.mu.RLock()
	defer slang.mu.RUnlock()

	return slang.currentNS
}

// Parent always returns nil to represent this is the root scope.
func (slang *Slang) Parent() sabre.Scope {
	return nil
}

func (slang *Slang) resolveAny(symbol string, syms ...nsSymbol) (sabre.Value, error) {
	for _, s := range syms {
		v, found := slang.bindings[s]
		if found {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unable to resolve symbol: %v", symbol)
}

func (slang *Slang) splitSymbol(symbol string) (*nsSymbol, error) {
	sep := string(nsSeparator)
	if symbol == sep {
		return &nsSymbol{
			NS:   slang.currentNS,
			Name: symbol,
		}, nil
	}

	parts := strings.SplitN(symbol, sep, 2)
	if len(parts) < 2 {
		return &nsSymbol{
			NS:   slang.currentNS,
			Name: symbol,
		}, nil
	}

	if strings.Contains(parts[1], sep) && parts[1] != sep {
		return nil, fmt.Errorf("invalid qualified symbol: '%s'", symbol)
	}

	return &nsSymbol{
		NS:   parts[0],
		Name: parts[1],
	}, nil
}

type nsSymbol struct {
	NS   string
	Name string
}

func (s nsSymbol) WithNS(ns string) nsSymbol {
	s.NS = ns
	return s
}
