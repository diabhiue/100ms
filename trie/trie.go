package trie

import (
    "github.com/emirpasic/gods/sets/treeset"
    "github.com/emirpasic/gods/utils"
)

const (
    ALPHANUMERIC_SIZE = 62
)

var charToIndex = map[string]int{
    "0": 0,
    "1": 1,
    "2": 2,
    "3": 3,
    "4": 4,
    "5": 5,
    "6": 6,
    "7": 7,
    "8": 8,
    "9": 9,
    "a": 10,
    "b": 11,
    "c": 12,
    "d": 13,
    "e": 14,
    "f": 15,
    "g": 16,
    "h": 17,
    "i": 18,
    "j": 19,
    "k": 20,
    "l": 21,
    "m": 22,
    "n": 23,
    "o": 24,
    "p": 25,
    "q": 26,
    "r": 27,
    "s": 28,
    "t": 29,
    "u": 30,
    "v": 31,
    "w": 32,
    "x": 33,
    "y": 34,
    "z": 35,
    "A": 36,
    "B": 37,
    "C": 38,
    "D": 39,
    "E": 40,
    "F": 41,
    "G": 42,
    "H": 43,
    "I": 44,
    "J": 45,
    "K": 46,
    "L": 47,
    "M": 48,
    "N": 49,
    "O": 50,
    "P": 51,
    "Q": 52,
    "R": 53,
    "S": 54,
    "T": 55,
    "U": 56,
    "V": 57,
    "W": 58,
    "X": 59,
    "Y": 60,
    "Z": 61,
}

type Node struct {
    Char     string
    Numbers  treeset.Set
    Children [ALPHANUMERIC_SIZE]*Node
}

func NewNode(char string) *Node {
    node := &Node{
        Char:    char,
        Numbers: *treeset.NewWith(utils.Int64Comparator),
    }

    for i := 0; i < ALPHANUMERIC_SIZE; i++ {
        node.Children[i] = nil
    }
    return node
}

type Trie struct {
    Root *Node
}

func NewTrie() *Trie {
    root := NewNode("\000")
    return &Trie{Root: root}
}

func (t *Trie) Insert(word string, num int64) error {
    current := t.Root

    for i := 0; i < len(word); i++ {
        index := charToIndex[string(word[i])]
        if current.Children[index] == nil {
            current.Children[index] = NewNode(word)
        }
        current = current.Children[index]
    }
    current.Numbers.Add(num)
    return nil
}

func (t *Trie) Delete(word string, num int64) error {
    current := t.Root
    for i := 0; i < len(word); i++ {
        index := charToIndex[string(word[i])]
        current = current.Children[index]
    }
    current.Numbers.Remove(num)
    return nil
}

func (t *Trie) GetCounters(word string, limit int) []int64 {
    current := t.Root
    for i := 0; i < len(word); i++ {
        index := charToIndex[string(word[i])]
        if current.Children[index] == nil {
            return make([]int64, 0)
        }
        current = current.Children[index]
    }
    itr := current.Numbers.Iterator()
    itr.End()

    var Counters []int64
    for itr.Prev() && limit > 0 {
        value := itr.Value().(int64)
        Counters = append(Counters, value)
        limit--
    }
    return Counters
}
