package logs

import (
	"strings"
	"trie"
)

type LogStore struct {
	Size                     int
	Values                   []string
	StartCounter, EndCounter int64
	KeyCounterMap            map[int64]int64
	CounterKeyMap            map[int64]int64
	WordTrie                 *packages.trie.Trie
}

func NewLogStore(size int) *LogStore {
	Log := &LogStore{Size: size}
	return Log
}

func GetWords(str string) []string {
	return strings.Fields(str)
}

func (Logs *LogStore) deleteWords(str string, counter int64) error {
	Words := GetWords(str)

	for _, word := range Words {
		Logs.WordTrie.Delete(word)
	}
	return nil
}

func (Logs *LogStore) insertWords(str string, counter int64) error {
	Words := GetWords(str)

	for _, word := range Words {
		Logs.WordTrie.Insert(word, counter)
	}
	return nil
}

func (Logs *LogStore) add(key int64, value string) error {

	PrevValue, IsKeyPresent := Logs.KeyCounterMap[key]
	if IsKeyPresent {
		Counter := KeyCounterMap[key]
		index := Counter % int64(Logs.Size)
		Logs.deleteWords(Logs.Values[index], Counter)
		Logs.Values[index] = value
		Logs.insertWords(value, Counter)
		return nil
	}

	if Logs.StartCounter%int64(Logs.Size) == Logs.EndCounter%int64(Logs.Size) {
		index := Logs.StartCounter % int64
		PrevValue := Logs.Values[index]
		Logs.deleteWords(PrevValue, Logs.StartCounter)

		PrevKey = Logs.CounterKeyMap[Logs.StartCounter]
		delete(Logs.KeyCounterMap, PrevKey)
		delete(Logs.CounterKeyMap, Logs.StartCounter)
		Logs.StartCounter++
	}

	index := Logs.EndCounter % int64(Logs.Size)
	Logs.Values[index] = value
	Logs.insertWords(word, Logs.EndCounter)
	Logs.KeyCounterMap[key] = Logs.EndCounter
	Logs.CounterKeyMap[Logs.EndCounter] = key
	Logs.EndCounter++

	return nil
}

func (Logs LogStore) search(word string, limit int) []int64 {
	Counters := Logs.WordTrie.GetNumbers(word)
	itr := Counters.Iterator()
	itr.End()

	var Keys []int64
	for itr.Prev() && limit > 0 {
		Keys = append(Keys, Logs.CounterKeyMap[itr.Value()])
		limit--
	}
	return Keys
}
