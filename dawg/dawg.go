// Package dawg provides a Directed Acyclic Word Graph.
// A DAWG is a data structure optimized for fast string lookups.
package dawg

// A DAWG is the main type provided by this package.
type DAWG struct {
	root *node
}

// A node is a recursive structure that represents a node of a DAWG.
type node struct {
	children map[nodeValue]*node
	eow      bool
}

// A nodeValue is the type stored in the value field of each node.
type nodeValue rune

// New creates a new DAWG from a vocabulary.
func New(vocabulary []string) *DAWG {
	d := &DAWG{&node{}}
	for _, word := range vocabulary {
		d.Insert(word)
	}
	return d
}

// Insert a word into the DAWG.
func (d *DAWG) Insert(word string) {
	current := d.root
	for _, k := range []nodeValue(word) {
		if current.children == nil {
			current.children = make(map[nodeValue]*node)
		}
		if next, ok := current.children[k]; ok {
			current = next
		} else {
			next = &node{}
			current.children[k] = next
			current = next
		}
	}
	current.eow = true
}

// Contains return true when the word is in the DAWG.
func (d *DAWG) Contains(word string) bool {
	current := d.root
	for _, k := range []nodeValue(word) {
		if current.children == nil {
			return false
		}
		if next, ok := current.children[k]; ok {
			current = next
		} else {
			return false
		}
	}
	return current.eow
}

// Returns a list of words of this DAWG that are prefixes of the given word.
func (d *DAWG) Prefixes(word string) []string {
	current := d.root
	res := []string{}
	w := []rune(word)
	for i, k := range w {
		if current.children == nil {
			break
		}
		if next, ok := current.children[nodeValue(k)]; ok {
			current = next
			if current.eow {
				res = append(res, string(w[:i+1]))
			}
		} else {
			break
		}
	}
	return res
}

// Returns a channel filled with words of this DAWG that are prefixes of the given word.
// Not implemented
func (d *DAWG) IterPrefixes(word string) chan string {
	return make(chan string)
}

// BUG(rhcarvalho): IterPrefixes is not implemented yet.

// Compact the DAWG by sharing common suffixes.
// Returns the number of trimmed branches.
func (d *DAWG) Compact() int {
	return 0
}

// BUG(rhcarvalho): Compact is not implemented yet.
