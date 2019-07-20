// Package dawg provides a Directed Acyclic Word Graph.
// A DAWG is a data structure optimized for fast string lookups.
package dawg

// A DAWG is the main type provided by this package.
type DAWG struct {
	root *node
}

// A node is a recursive structure that represents a node of a DAWG.
type node struct {
	children map[rune]*node
	eow      bool
}

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
	for _, k := range word {
		if current.children == nil {
			current.children = make(map[rune]*node)
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
	for _, k := range word {
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
func (d *DAWG) Prefixes(word []rune) (prefixes [][]rune) {
	current := d.root
	var prefix []rune
	for _, k := range word {
		if current.children == nil {
			break
		}
		if next, ok := current.children[k]; ok {
			prefix = append(prefix, k)
			if next.eow {
				prefixes = append(prefixes, prefix)
			}
			current = next
		} else {
			break
		}
	}
	return
}

// PrefixesString is like Prefixes, except that it works with strings.
func (d *DAWG) PrefixesString(word string) (prefixes []string) {
	for _, prefix := range d.Prefixes([]rune(word)) {
		prefixes = append(prefixes, string(prefix))
	}
	return
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

//====================================================================================

// Returns raw DAWG content as bytes.
//func (d *DAWG) ToBytes() []byte {
// Not implemented
/*
   cdef stringstream stream
   self.dct.Write(<ostream *> &stream)
   cdef bytes res = stream.str()
   return res
*/
//	return []byte{}
//}

// Loads DAWG from bytes ``data``.
//func (d *DAWG) FromBytes(data []byte) {
// Not implemented
/*
   cdef string s_data = data
   cdef stringstream* stream = new stringstream(s_data)

   try:
       res = self.dct.Read(<istream *> stream)

       if not res:
           self.dct.Clear()
           raise IOError("Invalid data format")

       return self
   finally:
       del stream
*/
//}

/*
// Loads DAWG from a file-like object.
func (d *DAWG) Read(f io.Reader) {
	d.FromBytes(f.Read())
}

// Writes DAWG to a file-like object.
func (d *DAWG) Write(f io.Writer) {
	f.Write(d.ToBytes())
}
*/

// Loads DAWG from a file.
//func (d *DAWG) Load(path string) {
// Not implemented
/*
   if isinstance(path, unicode):
       path = path.encode(sys.getfilesystemencoding())

   cdef ifstream stream
   stream.open(path, iostream.binary)

   res = self.dct.Read(<istream*> &stream)

   stream.close()

   if not res:
       self.dct.Clear()
       raise IOError("Invalid data format")

   return self
*/
//}

// Saves DAWG to a file.
///func (d *DAWG) Save(path string) {
// Not implemented
/*
   with open(path, 'wb') as f:
       self.write(f)
*/
//}
