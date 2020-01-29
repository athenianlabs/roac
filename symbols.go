package main

const MaxSymbols = 1024

var symbolTable = make(map[int]string, MaxSymbols)
var inverseSymbolTable = make(map[string]int, MaxSymbols)

func AddSymbol(s string) int {
	id := len(symbolTable)
	symbolTable[id] = s
	inverseSymbolTable[s] = id
	return id
}

func GetSymbolByID(id int) (string, bool) {
	s, ok := symbolTable[id]
	return s, ok
}

func GetSymbolIDByString(s string) (int, bool) {
	id, ok := inverseSymbolTable[s]
	return id, ok
}
