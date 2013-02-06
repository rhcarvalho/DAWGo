package dawg

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// the map keys are space-separated words to construct the DAWG.
var dawgs = map[string]*DAWG{
	"": &DAWG{&node{}},
	"g": &DAWG{&node{children: map[rune]*node{
		'g': &node{eow: true},
	}}},
	"go": &DAWG{&node{children: map[rune]*node{
		'g': &node{children: map[rune]*node{
			'o': &node{eow: true},
		}},
	}}},
	"g go": &DAWG{&node{children: map[rune]*node{
		'g': &node{eow: true, children: map[rune]*node{
			'o': &node{eow: true},
		}},
	}}},
	"g t": &DAWG{&node{children: map[rune]*node{
		'g': &node{eow: true},
		't': &node{eow: true},
	}}},
	"go t": &DAWG{&node{children: map[rune]*node{
		'g': &node{children: map[rune]*node{
			'o': &node{eow: true},
		}},
		't': &node{eow: true},
	}}},
	"语 语言 信 信息 处 处理": &DAWG{&node{children: map[rune]*node{
		'处': &node{eow: true, children: map[rune]*node{
			'理': &node{eow: true},
		}},
		'语': &node{eow: true, children: map[rune]*node{
			'言': &node{eow: true},
		}},
		'信': &node{eow: true, children: map[rune]*node{
			'息': &node{eow: true},
		}},
	}}},
}

func TestNew(t *testing.T) {
	for words, d := range dawgs {
		nd := New(strings.Fields(words))
		if !dawgsEqual(nd, d) {
			t.Errorf("DAWG should be %v, got %v", d, nd)
		}
	}
}

func TestContains(t *testing.T) {
	type query struct {
		key string
		res bool
	}
	tests := []struct {
		words   string
		queries []query
	}{{
		"g", []query{
			{"g", true},
			{"go", false},
			{"z", false},
		}}, {
		"go", []query{
			{"g", false},
			{"go", true},
			{"golang", false},
		}}, {
		"g go", []query{
			{"g", true},
			{"go", true},
			{"golang", false},
		}}, {
		"g t", []query{
			{"g", true},
			{"t", true},
			{"golang", false},
			{"tornado", false},
			{"z", false},
		}}, {
		"go t", []query{
			{"g", false},
			{"go", true},
			{"t", true},
			{"golang", false},
			{"tornado", false},
			{"z", false},
		}}, {
		"语 语言 信 信息 处 处理", []query{
			{"语", true},
			{"信", true},
			{"处", true},
			{"言", false},
			{"息", false},
			{"理", false},
			{"语言", true},
			{"信息", true},
			{"处理", true},
			{"语言信息处理", false},
		}},
	}
	for _, test := range tests {
		d, ok := dawgs[test.words]
		if !ok {
			t.Errorf("Missing DAWG for words %#v", test.words)
			continue
		}
		for _, q := range test.queries {
			if ok := d.Contains(q.key); ok != q.res {
				t.Errorf("DAWG(%#v).Contains(%#v) should be %v, got %v", test.words, q.key, q.res, ok)
			}
		}
	}
}

func TestPrefixes(t *testing.T) {
	type query struct {
		key string
		res []string
	}
	tests := []struct {
		words   string
		queries []query
	}{{
		"g go", []query{
			{"", []string{}},
			{"g", []string{"g"}},
			{"go", []string{"g", "go"}},
			{"golang", []string{"g", "go"}},
			{"python", []string{}},
		}}, {
		"g t", []query{
			{"g", []string{"g"}},
			{"t", []string{"t"}},
			{"golang", []string{"g"}},
			{"tornado", []string{"t"}},
			{"z", []string{}},
		}}, {
		"", []query{
			{"", []string{}},
			{"g", []string{}},
			{"golang", []string{}},
		}}, {
		"语 语言 信 信息 处 处理", []query{
			{"语言信息处理", []string{"语", "语言"}},
		}},
	}
	for _, test := range tests {
		d, ok := dawgs[test.words]
		if !ok {
			t.Errorf("Missing DAWG for words %#v", test.words)
			continue
		}
		for _, q := range test.queries {
			if prefs := d.PrefixesString(q.key); !slicesEqual(prefs, q.res) {
				t.Errorf("DAWG(%#v).Prefixes(%#v) should be %v, got %v", test.words, q.key, q.res, prefs)
			}
		}
	}
}

// Helper functions for printings DAWGs.
func (d *DAWG) String() string {
	return fmt.Sprintf("&DAWG{%v}", d.root)
}

func (n *node) String() string {
	eowS, chldS := eowToString(n.eow), childrenToString(n.children)
	var sep string
	if eowS != "" && chldS != "" {
		sep = ", "
	} else {
		sep = ""
	}
	return fmt.Sprintf("&node{%s%s%s}", eowS, sep, chldS)
}

func eowToString(eow bool) string {
	if eow {
		return "eow: true"
	}
	return ""
}

func childrenToString(children map[rune]*node) string {
	if children == nil {
		return ""
	}
	b := bytes.NewBufferString("children: map[rune]*node{\n")
	for k, nd := range children {
		fmt.Fprintf(b, "'%s': %v,\n", k, nd)
	}
	b.WriteByte('}')
	return b.String()
}

// Helper functions for comparing DAWGs.
func dawgsEqual(d1, d2 *DAWG) bool {
	return nodesEqual(d1.root, d2.root)
}

func nodesEqual(x, y *node) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return false
	}
	if x.eow != y.eow {
		return false
	}
	for key, xChildNode := range x.children {
		if yChildNode, ok := y.children[key]; !ok {
			return false
		} else {
			if !nodesEqual(xChildNode, yChildNode) {
				return false
			}
		}
	}
	return true
}

// sliceEq defines equality based on its contents.
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
