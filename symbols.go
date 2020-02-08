package main

import "fmt"

const MaxSymbols = 1024

type Symbol struct {
	name     string
	t        NodeType
	st       StructuralNodeType
	id       int
	endLabel int
}

func (s Symbol) String() string {
	return fmt.Sprintf("Symbol '%s' ID: %d", s.name, s.id)
}

var (
	symbolTable        = make(map[int]*Symbol, MaxSymbols)
	inverseSymbolTable = make(map[string]int, MaxSymbols)
)

func AddSymbol(s string, t NodeType, st StructuralNodeType, endLabel int) *Symbol {
	if _, exists := inverseSymbolTable[s]; exists {
		fatal("symbol %s already declared", s)
	}
	id := len(symbolTable)
	symbolTable[id] = &Symbol{
		name:     s,
		t:        t,
		st:       st,
		id:       id,
		endLabel: endLabel,
	}
	inverseSymbolTable[s] = id
	return symbolTable[id]
}

func GetSymbolByID(id int) *Symbol {
	s, ok := symbolTable[id]
	if !ok {
		fatal("symbol with id %d does not exists", id)
	}
	return s
}

func GetSymbolByString(s string) *Symbol {
	id, ok := inverseSymbolTable[s]
	if !ok {
		fatal("symbol %s does not exists", s)
	}
	return GetSymbolByID(id)
}
