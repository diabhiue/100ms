package logs

import (
	// "fmt"
	// "github.com/diabhiue/100ms/trie"
	"strings"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
)

type LogStore struct {
	Size                     int
	Values                   []string
	StartCounter, EndCounter int64
	KeyCounterMap            map[int64]int64
	CounterKeyMap            map[int64]int64
	// WordTrie                 *trie.Trie
	WordMap map[string]*treeset.Set
}

func NewLogStore(size int) *LogStore {
	// Log := &LogStore{Size: size, WordTrie: trie.NewTrie()}
	Log := &LogStore{
		Size:          size,
		Values:        make([]string, size),
		StartCounter:  0,
		EndCounter:    0,
		KeyCounterMap: make(map[int64]int64),
		CounterKeyMap: make(map[int64]int64),
		WordMap:       make(map[string]*treeset.Set),
	}
	return Log
}

func GetWords(str string) []string {
	return strings.Fields(str)
}

func (Logs *LogStore) deleteWords(str string, counter int64) error {
	Words := GetWords(str)

	for _, word := range Words {
		// Logs.WordTrie.Delete(word, counter)
		Logs.WordMap[word].Remove(counter)
		if Logs.WordMap[word].Empty() {
			delete(Logs.WordMap, word)
		}
	}
	return nil
}

func (Logs *LogStore) insertWords(str string, counter int64) error {
	Words := GetWords(str)

	for _, word := range Words {
		// Logs.WordTrie.Insert(word, counter)
		_, isPresent := Logs.WordMap[word]
		if !isPresent {
			Logs.WordMap[word] = treeset.NewWith(utils.Int64Comparator)
		}
		Logs.WordMap[word].Add(counter)
	}
	return nil
}

func (Logs *LogStore) Add(key int64, value string) error {

	PrevCounter, IsKeyPresent := Logs.KeyCounterMap[key]
	if IsKeyPresent {
		Counter := PrevCounter
		index := Counter % int64(Logs.Size)
		Logs.deleteWords(Logs.Values[index], Counter)
		Logs.Values[index] = value
		Logs.insertWords(value, Counter)

		return nil
	}

	frontIndex := Logs.StartCounter % int64(Logs.Size)
	rearIndex := Logs.EndCounter % int64(Logs.Size)

	if frontIndex == rearIndex && Logs.Values[frontIndex] != "" {
		index := frontIndex
		PrevValue := Logs.Values[index]
		Logs.deleteWords(PrevValue, Logs.StartCounter)

		PrevKey := Logs.CounterKeyMap[Logs.StartCounter]
		delete(Logs.KeyCounterMap, PrevKey)
		delete(Logs.CounterKeyMap, Logs.StartCounter)
		Logs.StartCounter++
	}

	index := Logs.EndCounter % int64(Logs.Size)
	Logs.Values[index] = value
	Logs.insertWords(value, Logs.EndCounter)
	Logs.KeyCounterMap[key] = Logs.EndCounter
	Logs.CounterKeyMap[Logs.EndCounter] = key
	Logs.EndCounter++

	return nil
}

func (Logs LogStore) Search(word string, limit int) []int64 {
	// Counters := Logs.WordTrie.GetCounters(word, limit)
	Counters, isPresent := Logs.WordMap[word]
	if !isPresent {
		return make([]int64, 0)
	}
	itr := Counters.Iterator()
	itr.End()

	var Keys []int64
	// for _, counter := range Counters {
	// 	Keys = append(Keys, Logs.CounterKeyMap[counter])
	// }

	for itr.Prev() && limit > 0 {
		Keys = append(Keys, Logs.CounterKeyMap[itr.Value().(int64)])
		limit--
	}

	return Keys
}
