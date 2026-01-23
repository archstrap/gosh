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
	escapeRunes           = `\`
	ioRedirectRunes       = `><`
	digitRunes            = `0123456789`
)

const (
	unknownRuneClass runeTokenClass = iota
	spaceRuneClass
	nonEscapingQuoteRuneClass
	escapingQuoteRuneClass
	escapeRuneClass
	ioRedirectRuneClass
	digitRuneClass
	eofRuneClass
)

const (
	startState tokenizerState = iota
	inWordState
	nonEscapingQuoteState
	escapingQuoteState
	quotedEscapingState
	ioRedirectState
	escapeState
)

const (
	wordToken TokenType = iota
	ioRedirectionToken
)

var (
	specialRune map[string]bool = map[string]bool{
		escapinngQuoteRunes: true,
		escapeRunes:         true,
		`$`:                 true,
	}
)

type Token struct {
	value     string
	tokenType TokenType
}

type Redirection struct {
	fileName   string
	appendOnly bool
}

func NewRedirection(fileName string) *Redirection {
	return &Redirection{fileName: fileName}
}

type Command struct {
	name         string
	args         []string
	redirections map[int]*Redirection
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
	tc.AddClassifier(escapinngQuoteRunes, escapingQuoteRuneClass)
	tc.AddClassifier(ioRedirectRunes, ioRedirectRuneClass)
	tc.AddClassifier(digitRunes, digitRuneClass)
	tc.AddClassifier(escapeRunes, escapeRuneClass)
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

func (tr *Tokenizer) getRuneDetails() (rune, runeTokenClass, error) {
	currentRune, _, err := tr.input.ReadRune()
	currentRuneType := tr.classifier.ClassifyRune(currentRune)

	if err == io.EOF {
		err = nil
		currentRuneType = eofRuneClass
	}

	return currentRune, currentRuneType, err

}

func (tr *Tokenizer) scan() (*Token, error) {

	state := startState
	var prevEscapeRune rune
	var value []rune
	var tokenType TokenType
	
	for {

		nextRune, nextRuneType, err := tr.getRuneDetails()

		if err != nil {
			return nil, err
		}

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
			case escapingQuoteRuneClass:
				state = escapingQuoteState
				tokenType = wordToken
			case escapeRuneClass:
				state = escapeState
				tokenType = wordToken
			case digitRuneClass:
				value = append(value, nextRune)
				_, nextToNextRuneType, _ := tr.getRuneDetails()
				if nextToNextRuneType == ioRedirectRuneClass {
					state = ioRedirectState
					tokenType = ioRedirectionToken
				} else {
					state = inWordState
					tokenType = wordToken
				}
				tr.input.UnreadRune()
			case ioRedirectRuneClass:
				state = ioRedirectState
				tokenType = ioRedirectionToken
				value = append(value, nextRune)
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
			case escapingQuoteRuneClass:
				state = escapingQuoteState
			case escapeRuneClass:
				state = escapeState
			case ioRedirectRuneClass:
				tr.input.UnreadRune()
				return NewToken(string(value), tokenType), nil
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
				return nil, fmt.Errorf("EOF after escape character")
			case escapingQuoteRuneClass:
				state = inWordState
			case escapeRuneClass:
				state = quotedEscapingState
				prevEscapeRune = nextRune
			default:
				value = append(value, nextRune)
			}
		case escapeState:
			switch nextRuneType {
			case eofRuneClass:
				return nil, fmt.Errorf("EOF after escape character")
			default:
				state = inWordState
				value = append(value, nextRune)
			}
		case quotedEscapingState:
			switch nextRuneType {
			case eofRuneClass:
				return nil, fmt.Errorf("EOF found while expecting a closing quote")
			default:
				state = escapingQuoteState
				if !specialRune[string(nextRune)] {
					value = append(value, prevEscapeRune)
				}
				value = append(value, nextRune)
			}
		case ioRedirectState:
			switch nextRuneType {
			case eofRuneClass:
				return nil, fmt.Errorf("EOF while expecting a io redirection")
			case ioRedirectRuneClass:
				state = ioRedirectState
				tokenType = ioRedirectionToken
				value = append(value, nextRune)
			case spaceRuneClass:
				return NewToken(string(value), tokenType), nil
			default:
				tr.input.UnreadRune()
				return NewToken(string(value), tokenType), nil
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

func (lx *Lexer) Parse() (*Token, error) {
	return (*Tokenizer)(lx).Next()
}

func Parse(s string) Command {

	lexer := NewLexer(s)
	var tokens []*Token
	for {
		token, err := lexer.Parse()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		tokens = append(tokens, token)
	}
	var commandName string
	var args []string
	redirections := make(map[int]*Redirection)

	position := 0

	for position < len(tokens) {
		token := tokens[position]
		if position == 0 {
			commandName = token.value
			position++
			continue
		}

		if token.tokenType == wordToken {
			args = append(args, token.value)
			position++
			continue
		}

		if token.tokenType == ioRedirectionToken {
			var fileName string
			if position+1 < len(tokens) {
				fileName = tokens[position+1].value
			}

			redirectionValue := strings.TrimSpace(token.value)
			fileDescriptor := -1
			appendOnly := false

			switch redirectionValue {
			case "<":
				fileDescriptor = 0
				appendOnly = false
			case ">", "1>":
				fileDescriptor = 1
				appendOnly = false
			case "2>":
				fileDescriptor = 2
				appendOnly = false
			case ">>", "1>>":
				fileDescriptor = 1
				appendOnly = true
			case "2>>":
				fileDescriptor = 2
				appendOnly = true

			}

			redirections[fileDescriptor] = &Redirection{fileName: fileName, appendOnly: appendOnly}
			position += 2

		}

	}

	return Command{name: commandName, args: args, redirections: redirections}

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
