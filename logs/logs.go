package logs

import (
    "github.com/diabhiue/100ms/trie"
    "strings"
)

// LogStore struct
// This stores max size of the LogStore, all values provided by user,
// StartCounter, EndCounter, Key->Counter map, Counter->Key map and a
// trie which stores all keys for every word in values in latest `size` logs
type LogStore struct {
    Size                     int
    Values                   []string // this acts as circular queue for values
    StartCounter, EndCounter int64    // start and end counter between which all records are live
    KeyCounterMap            map[int64]int64
    CounterKeyMap            map[int64]int64
    WordTrie                 *trie.Trie
}

// Initialize the LogStore and return it
func NewLogStore(size int) *LogStore {
    Log := &LogStore{
        Size:          size,
        Values:        make([]string, size),
        StartCounter:  0,
        EndCounter:    0,
        KeyCounterMap: make(map[int64]int64),
        CounterKeyMap: make(map[int64]int64),
        WordTrie:      trie.NewTrie(),
    }
    return Log
}

// Returns whitespace separated strings as an array
func GetWords(str string) []string {
    return strings.Fields(str)
}

// Delete Counter from each of the word in `str`
func (Logs *LogStore) deleteWords(str string, counter int64) error {
    Words := GetWords(str)

    for _, word := range Words {
        Logs.WordTrie.Delete(word, counter)
    }
    return nil
}

// Insert Counter corresponding to each word in `str`
// Generally all counters will come in sorted order, but for an existing key,
// it will less than max counter, hence BST will handle this case nicely
func (Logs *LogStore) insertWords(str string, counter int64) error {
    Words := GetWords(str)

    for _, word := range Words {
        Logs.WordTrie.Insert(word, counter)
    }
    return nil
}

// Register (key, value) pair in `LogStore`
func (Logs *LogStore) Add(key int64, value string) error {

    // Case 1: Check whether Key is already registered
    // If Key is already present, then just update the values at corresponding
    // index and update trie as well by deleting the Counter from previous words
    // in string and insert PrevCounter corresponding to words in new value string
    // Note that, this implementation doesn't append new value, instead it just
    // updates key with new value
    PrevCounter, IsKeyPresent := Logs.KeyCounterMap[key]
    if IsKeyPresent {
        Counter := PrevCounter
        index := Counter % int64(Logs.Size)
        Logs.deleteWords(Logs.Values[index], Counter)
        Logs.Values[index] = value
        Logs.insertWords(value, Counter)

        return nil
    }

    // Since `Logs.Values` array acts as a circular queue, fetch front and rear
    // index
    frontIndex := Logs.StartCounter % int64(Logs.Size)
    rearIndex := Logs.EndCounter % int64(Logs.Size)

    // Case 2: Circular queue is full
    // To resolve this issue, remove (key, value) pair from `frontIndex`, and other
    // data structures like Key->Counter map, Counter->Key map and finally increment
    // StartCounter
    if frontIndex == rearIndex && Logs.Values[frontIndex] != "" {
        index := frontIndex
        PrevValue := Logs.Values[index]
        Logs.deleteWords(PrevValue, Logs.StartCounter)

        PrevKey := Logs.CounterKeyMap[Logs.StartCounter]
        delete(Logs.KeyCounterMap, PrevKey)
        delete(Logs.CounterKeyMap, Logs.StartCounter)
        Logs.StartCounter++
    }

    // Case 3: Add (key, value) pair at rear end of circular queue
    // Update other data structures as well
    index := Logs.EndCounter % int64(Logs.Size)
    Logs.Values[index] = value
    Logs.insertWords(value, Logs.EndCounter)
    Logs.KeyCounterMap[key] = Logs.EndCounter
    Logs.CounterKeyMap[Logs.EndCounter] = key
    Logs.EndCounter++

    return nil
}

// Returns latest `limit` Keys whose value contains `word`
func (Logs LogStore) Search(word string, limit int) []int64 {

    // Get latest `limit` Counters which were present in `word``
    Counters := Logs.WordTrie.GetCounters(word, limit)

    // Store Key for all `Counters` using Counter->Key map
    var Keys []int64
    for _, counter := range Counters {
        Keys = append(Keys, Logs.CounterKeyMap[counter])
    }

    return Keys
}
