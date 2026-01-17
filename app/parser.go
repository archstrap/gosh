package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

/*

   TOKENIZER will give us the raw tokens
	 LEXER will give us the word in string

*/

type TokenType int
type runeTokenClass int
type tokenizerState int

const (
	spaceRunes            = " \t\r\n"
	nonEscapingQuoteRunes = `'`
	escapinngQuoteRunes   = `"`
)

const (
	unknownRuneClass runeTokenClass = iota
	spaceRuneClass
	nonEscapingQuoteRuneClass
	escapinngQuoteRuneClass
	eofRuneClass
)

const (
	startState tokenizerState = iota
	inWordState
	nonEscapingQuoteState
	escapingQuoteState
)

const (
	wordToken TokenType = iota
)

type Token struct {
	value     string
	tokenType TokenType
}

func NewToken(value string, tokenType TokenType) *Token {
	return &Token{value: value, tokenType: tokenType}
}

/*-------------------- [ TokenClassifier ] ----------------------*/
type TokenClassifier map[rune]runeTokenClass

func NewDefaultClassifier() TokenClassifier {
	tc := TokenClassifier{}
	tc.AddClassifier(spaceRunes, spaceRuneClass)
	tc.AddClassifier(nonEscapingQuoteRunes, nonEscapingQuoteRuneClass)
	tc.AddClassifier(escapinngQuoteRunes, escapinngQuoteRuneClass)
	return tc
}

func (tc TokenClassifier) AddClassifier(k string, v runeTokenClass) {

	for _, item := range k {
		tc[item] = v
	}
}

func (tc TokenClassifier) ClassifyRune(runeInput rune) runeTokenClass {
	return tc[runeInput]
}

/*-------------------- [ Tokenizer ] ----------------------*/

type Tokenizer struct {
	input      bufio.Reader
	classifier TokenClassifier
}

func (tr *Tokenizer) scan() (*Token, error) {

	state := startState
	var value []rune
	var tokenType TokenType

	//	Debug(Red, tr.classifier)

	for {

		nextRune, _, err := tr.input.ReadRune()
		nextRuneType := tr.classifier.ClassifyRune(nextRune)

		//		Debug(Green, nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil
		} else if err != nil {
			return nil, err
		}
		//		Debug(Blue, fmt.Sprintf("state:= %v", state))

		switch state {
		case startState:
			switch nextRuneType {
			case eofRuneClass:
				return nil, io.EOF
			case spaceRuneClass:
				{
				}
			case nonEscapingQuoteRuneClass:
				state = nonEscapingQuoteState
				tokenType = wordToken
			case escapinngQuoteRuneClass:
				state = escapingQuoteState
				tokenType = wordToken
			default:
				state = inWordState
				tokenType = wordToken
				value = append(value, nextRune)
			}
		case inWordState:
			switch nextRuneType {
			case eofRuneClass:
				return NewToken(string(value), tokenType), nil
			case spaceRuneClass:
				return NewToken(string(value), tokenType), nil
			case nonEscapingQuoteRuneClass:
				state = nonEscapingQuoteState
			case escapinngQuoteRuneClass:
				state = escapingQuoteState
			default:
				tokenType = wordToken
				value = append(value, nextRune)
			}

		case nonEscapingQuoteState:
			switch nextRuneType {
			case eofRuneClass:
				return NewToken(string(value), tokenType), nil
			case nonEscapingQuoteRuneClass:
				state = inWordState
			default:
				value = append(value, nextRune)
			}
		case escapingQuoteState:
			switch nextRuneType {
			case eofRuneClass:
				return NewToken(string(value), tokenType), nil
			case escapinngQuoteRuneClass:
				state = inWordState
			default:
				value = append(value, nextRune)
			}
		default:
			return nil, fmt.Errorf("unexpected state: %v", state)
		}

		//		Debug(Cyan, value)
	}
}

func (tr *Tokenizer) Next() (*Token, error) {
	return tr.scan()
}

/*-------------------- [ Lexer ] ----------------------*/
type Lexer Tokenizer

func NewLexer(s string) *Lexer {

	//	Debug(Yellow, s)
	reader := strings.NewReader(s)
	classifier := NewDefaultClassifier()

	tr := &Tokenizer{
		input:      *bufio.NewReader(reader),
		classifier: classifier,
	}
	return (*Lexer)(tr)
}

func (lx *Lexer) Next() (string, error) {

	for {
		token, err := (*Tokenizer)(lx).Next()
		if err != nil {
			return "", err
		}
		switch token.tokenType {
		case wordToken:
			return token.value, nil
		default:
			return "", fmt.Errorf("Unknown token type: %v", token.tokenType)
		}
	}
}

func Split(s string) ([]string, error) {
	lexer := NewLexer(s)
	var values []string
	for {
		value, err := lexer.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}
