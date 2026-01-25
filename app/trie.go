package main

type Container map[rune]*Node

type Trie struct {
	root *Node
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
	return &Trie{root: NewNode()}
}

func (t *Trie) Insert(s string) {
	currentNode := t.root
	for _, ch := range s {
		nextNode := currentNode.GetNodeAt(ch)
		if nextNode == nil {
			nextNode = currentNode.SetNodeAt(ch)
		}
		currentNode = nextNode
	}
	currentNode.eow = true
}

func (t *Trie) InsertAll(input ...string) {
	for _, item := range input {
		t.Insert(item)
	}
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
