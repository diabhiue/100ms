package logs

import (
	"github.com/diabhiue/100ms/trie"
	"strings"
)

type LogStore struct {
	Size                     int
	Values                   []string
	StartCounter, EndCounter int64
	KeyCounterMap            map[int64]int64
	CounterKeyMap            map[int64]int64
	WordTrie                 *trie.Trie
}

func NewLogStore(size int) *LogStore {
	Log := &LogStore{
		Size:          size,
		Values:        make([]string, size),
		StartCounter:  0,
		EndCounter:    0,
		KeyCounterMap: make(map[int64]int64),
		CounterKeyMap: make(map[int64]int64),
        WordTrie: trie.NewTrie(),
	}
	return Log
}

func GetWords(str string) []string {
	return strings.Fields(str)
}

func (Logs *LogStore) deleteWords(str string, counter int64) error {
	Words := GetWords(str)

	for _, word := range Words {
		Logs.WordTrie.Delete(word, counter)
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
	Counters := Logs.WordTrie.GetCounters(word, limit)

	var Keys []int64
	for _, counter := range Counters {
		Keys = append(Keys, Logs.CounterKeyMap[counter])
	}

	return Keys
}
