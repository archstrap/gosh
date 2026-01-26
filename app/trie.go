package main

import "strings"

type Container map[rune]*Node

type Trie struct {
	root *Node
	// longest Common Prefix
	lcp         strings.Builder
	isFirstWord bool
}

type Node struct {
	eow      bool
	children Container
}

func NewNode() *Node {
	return &Node{eow: false, children: make(Container)}
}

func (n *Node) GetNodeAt(c rune) *Node {
	return n.children[c]
}

func (n *Node) SetNodeAt(c rune) *Node {
	nextNode := NewNode()
	n.children[c] = nextNode
	return nextNode
}

func NewTrie() *Trie {
	return &Trie{root: NewNode(), isFirstWord: true}
}

func (t *Trie) Insert(s string) {
	currentNode := t.root

	var currLcp strings.Builder

	for _, ch := range s {
		nextNode := currentNode.GetNodeAt(ch)
		if nextNode == nil {
			nextNode = currentNode.SetNodeAt(ch)
		} else {
			currLcp.WriteRune(ch)
		}

		currentNode = nextNode
	}
	currentNode.eow = true
	lcp := CommponPrefix(currLcp.String(), t.lcp.String())

	if "" == lcp && t.isFirstWord {
		t.lcp.WriteString(s)
	} else if len(lcp) < t.lcp.Len() {
		t.lcp.Reset()
		t.lcp.WriteString(lcp)
	}

}

func (t *Trie) InsertAll(input ...string) {
	for _, item := range input {
		t.Insert(item)
	}
}

func (t *Trie) LongestCommonPrefix() string {
	return t.lcp.String()
}

func CommponPrefix(a string, b string) string {

	n, m := len(a), len(b)
	var common strings.Builder
	for i := 0; i < min(n, m); i++ {
		if a[i] != b[i] {
			break
		}
		common.WriteByte(a[i])
	}
	return common.String()

}

func (t *Trie) SearchAll(s string) []string {
	currentNode := t.root
	for _, ch := range s {
		nextNode := currentNode.GetNodeAt(ch)
		if nextNode == nil {
			return []string{}
		}
		currentNode = nextNode
	}

	var result []string
	CompleteAllWords(currentNode, s, &result)
	return result
}

func CompleteAllWords(node *Node, path string, results *[]string) {
	if node.eow {
		*results = append(*results, path)
	}

	for ch, childNode := range node.children {
		CompleteAllWords(childNode, path+string(ch), results)
	}
}
