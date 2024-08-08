// Code generated from Calc.g4 by ANTLR 4.8. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 8, 35, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 6, 6, 25, 10,
	6, 13, 6, 14, 6, 26, 3, 7, 6, 7, 30, 10, 7, 13, 7, 14, 7, 31, 3, 7, 3,
	7, 2, 2, 8, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 3, 2, 4, 3, 2, 50, 59,
	5, 2, 11, 12, 15, 15, 34, 34, 2, 36, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2,
	2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2,
	2, 3, 15, 3, 2, 2, 2, 5, 17, 3, 2, 2, 2, 7, 19, 3, 2, 2, 2, 9, 21, 3, 2,
	2, 2, 11, 24, 3, 2, 2, 2, 13, 29, 3, 2, 2, 2, 15, 16, 7, 44, 2, 2, 16,
	4, 3, 2, 2, 2, 17, 18, 7, 49, 2, 2, 18, 6, 3, 2, 2, 2, 19, 20, 7, 45, 2,
	2, 20, 8, 3, 2, 2, 2, 21, 22, 7, 47, 2, 2, 22, 10, 3, 2, 2, 2, 23, 25,
	9, 2, 2, 2, 24, 23, 3, 2, 2, 2, 25, 26, 3, 2, 2, 2, 26, 24, 3, 2, 2, 2,
	26, 27, 3, 2, 2, 2, 27, 12, 3, 2, 2, 2, 28, 30, 9, 3, 2, 2, 29, 28, 3,
	2, 2, 2, 30, 31, 3, 2, 2, 2, 31, 29, 3, 2, 2, 2, 31, 32, 3, 2, 2, 2, 32,
	33, 3, 2, 2, 2, 33, 34, 8, 7, 2, 2, 34, 14, 3, 2, 2, 2, 5, 2, 26, 31, 3,
	8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'*'", "'/'", "'+'", "'-'",
}

var lexerSymbolicNames = []string{
	"", "MUL", "DIV", "ADD", "SUB", "NUMBER", "WHITESPACE",
}

var lexerRuleNames = []string{
	"MUL", "DIV", "ADD", "SUB", "NUMBER", "WHITESPACE",
}

type CalcLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewCalcLexer(input antlr.CharStream) *CalcLexer {

	l := new(CalcLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Calc.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// CalcLexer tokens.
const (
	CalcLexerMUL        = 1
	CalcLexerDIV        = 2
	CalcLexerADD        = 3
	CalcLexerSUB        = 4
	CalcLexerNUMBER     = 5
	CalcLexerWHITESPACE = 6
)
