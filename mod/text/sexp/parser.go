package sexp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

//go:generate stringer -output parser_string.go -type=NodeType

var ErrInvalidFormat = errors.New("invalid format")

type NodeType int

const (
	NodeNil NodeType = iota
	NodeSymbol
	NodeInt
	NodeString
	NodeCell
)

type Node struct {
	Type  NodeType
	Value any
	Car   *Node
	Cdr   *Node
}

func (node *Node) SymbolValue() (string, bool) {
	if node == nil || node.Type != NodeSymbol {
		return "", false
	}

	s, _ := node.Value.(string)
	return s, true
}

func (node *Node) IntValue() (int64, bool) {
	if node == nil || node.Type != NodeInt {
		return 0, false
	}

	n, _ := node.Value.(int64)
	return n, true
}

func (node *Node) StringValue() (string, bool) {
	if node == nil || node.Type != NodeString {
		return "", false
	}

	s, _ := node.Value.(string)
	return s, true
}

func (node *Node) String() string {
	if node == nil {
		return "nil"
	}

	var buf bytes.Buffer

	switch node.Type {
	case NodeCell:
		fmt.Fprint(&buf, "(")

		for curr := node; curr != nil; curr = curr.Cdr {
			if curr.Car != nil {
				fmt.Fprint(&buf, curr.Car)
			} else {
				fmt.Fprint(&buf, "nil")
			}

			if curr.Cdr == nil || curr.Cdr.Type == NodeNil {
				break
			}

			fmt.Fprint(&buf, " ")
		}

		fmt.Fprint(&buf, ")")
	case NodeSymbol:
		fmt.Fprint(&buf, node.Value)
	case NodeInt:
		fmt.Fprint(&buf, node.Value)
	case NodeString:
		fmt.Fprintf(&buf, "%q", node.Value)
	}

	return ""
}

type Parser struct {
	r *bufio.Reader
}

func New(rd io.Reader) *Parser {
	var br *bufio.Reader
	if r, ok := rd.(*bufio.Reader); ok {
		br = r
	} else {
		br = bufio.NewReader(rd)
	}

	return &Parser{r: br}
}

func (p *Parser) Parse() (*Node, error) {
	node, err := p.parseNode()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (p *Parser) parseNode() (*Node, error) {
	err := p.skipSpace()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, handleError(err)
	}

	r, _, err := p.r.ReadRune()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, handleError(err)
	}

	if r == '(' {
		// list
		node, err := p.parseList()
		if err != nil {
			return nil, err
		}

		r, _, err := p.r.ReadRune()
		if err != nil {
			return nil, handleError(err)
		}
		if r != ')' {
			return nil, ErrInvalidFormat
		}

		return node, nil
	}
	if isPrimitive(r) {
		p.r.UnreadRune()
		return p.parsePrimitive()
	}
	if r == '"' {
		return p.parseString()
	}

	return nil, ErrInvalidFormat
}

func (p *Parser) parseList() (*Node, error) {
	node := &Node{
		Type: NodeCell,
	}

	curr := node
	for {
		if err := p.skipSpace(); err != nil {
			return nil, handleError(err)
		}

		r, _, err := p.r.ReadRune()
		if err != nil {
			return nil, handleError(err)
		}

		_ = p.r.UnreadByte()

		if r == ')' {
			break
		}

		child, err := p.parseNode()
		if err == io.EOF {
			return nil, handleError(err)
		}
		if err != nil {
			return nil, err
		}

		if node.Car != nil {
			cdr := &Node{
				Type: NodeCell,
			}
			curr.Cdr = cdr
			curr = cdr
		}

		curr.Car = child
	}

	return node, nil
}

func (p *Parser) parseString() (*Node, error) {
	var s strings.Builder

	for {
		r, _, err := p.r.ReadRune()
		if err != nil {
			return nil, handleError(err)
		}

		if r == '"' {
			break
		}

		if r == '\\' {
			r, _, err = p.r.ReadRune()
			if err != nil {
				return nil, handleError(err)
			}

			switch r {
			case '\\':
				r = '\\'
			case 'n':
				r = '\n'
			case 'r':
				r = '\n'
			case 't':
				r = '\t'
			case 'v':
				r = '\v'
			case 'b':
				r = '\b'
			case 'f':
				r = '\f'
			case '"':
			}
		}
		s.WriteRune(r)
	}

	return &Node{
		Type:  NodeString,
		Value: s.String(),
	}, nil
}

func (p *Parser) parsePrimitive() (*Node, error) {
	var b strings.Builder

	for {
		r, _, err := p.r.ReadRune()
		if err != nil {
			return nil, handleError(err)
		}

		if !isPrimitive(r) {
			_ = p.r.UnreadRune()
			break
		}

		b.WriteRune(r)
	}

	s := b.String()

	if s == "nil" {
		return &Node{
			Type: NodeNil,
		}, nil
	}

	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &Node{
			Type:  NodeInt,
			Value: n,
		}, nil
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return &Node{
			Type:  NodeInt,
			Value: f,
		}, nil
	}

	return &Node{
		Type:  NodeSymbol,
		Value: s,
	}, nil
}

func isPrimitive(r rune) bool {
	if r >= unicode.MaxASCII {
		return false
	}

	if unicode.IsLetter(r) {
		return true
	}
	if unicode.IsDigit(r) {
		return true
	}

	if strings.ContainsRune("!#$%&'*+,-.:<=>@[]^`{/}", r) {
		return true
	}

	return false
}

func isAtomEnd(r rune) bool {
	if isSpace(r) {
		return true
	}
	if r == '(' || r == ')' {
		return true
	}

	return false
}

func (p *Parser) skipSpace() error {
	for {
		r, _, err := p.r.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if isSpace(r) {
			continue
		}

		_ = p.r.UnreadRune()
		return nil
	}
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func handleError(err error) error {
	if errors.Is(err, io.EOF) {
		return io.ErrUnexpectedEOF
	}

	return fmt.Errorf("failed to read input: %v", err)
}
