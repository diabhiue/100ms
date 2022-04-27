## Problem Statement

To implement a Log Managment tool which stores key, value with some constraints like max key-value pairs to keep. For detail, please refer to the shared document.

Operations on the Log store
1. Add a `key`, `value` pair where key will be an integer and value will be string consisting of words separated by spaces. Note that if a key is already present, then instead of appending new value, update the value itself.
2. Given a `word`, search upto latest `limit` number of keys whose value contains that word.

## Instructions to run the code

Please follow this link - https://replit.com/@AbhiueAnand/100ms and run the code as demonstrated in the GIF below.

Example 1:


<table>
<tr>
<th>Input</th>
<th>Expected output</th>
</tr>
<tr>
<td>
  
```
2
ADD 25 the first
SEARCH the 1
ADD 56 the second log
SEARCH the 2
ADD 25 again second log
SEARCH second 2
ADD 16 the third log
SEARCH the 3
SEARCH fourth 1
END
```
  
</td>
<td>

```
25 
56 25 
56 25 
16 56 
NONE
END
```

</td>
</tr>
</table>

![](https://raw.githubusercontent.com/diabhiue/100ms/main/Demonstration%20-%20100ms%20-%201.gif)

Example 2

<table>
<tr>
<th>Input</th>
<th>Expected Output</th>
</tr>
<tr>
<td>
  
```
3
ADD 16 big university
ADD 36 peace and the war
SEARCH the 3
ADD 45 the war university
SEARCH war 1
SEARCH big 1
ADD 36 big war
SEARCH big 2
ADD 48 big cholocate
SEARCH big 2
END
```
  
</td>
<td>

```
36 
45 
16 
36 16 
48 36 
END
```

</td>
</tr>
</table>

![](https://raw.githubusercontent.com/diabhiue/100ms/main/Demonstration%20-%20100ms%20-%202.gif)

## Implementation details

In this particular problem, since we have to create a Log Management tool with max size(S), we need to create some data structure so that we can evict as per FIFO order. Only caveat is that if a key is already present, we just need to update the value. There are three operations which need to be done:

1. Initialize `LogStore` with max size `S`
2. `Add(Key int64, Value string)`
3. `Search(Word string, limit int)`

To accommodate 1., we need to have some sort of data structure with a max limit of key-value pairs, and if this data structure is full, we need to evict the oldest record. Since eviction policy will follow FIFO order, we can create a queue to store all key-value pairs. But if we want to fetch value for a key, it will take O(S) to search linearly for key-value pair. To make this efficient, we can create a queue but in an array(circular queue). So, now we can store the values in an array and maintain frontIndex and rearIndex which can be helpful while pushing/popping elements; and we can have a map of key->index. So, whenever we get a key, we can fetch index from that map and from that index we can fetch value.

For 2., adding key value will work fine as well with the above approach

For 3., since we need to find keys whose value has the `word` present in it, whenever we `Add` a key-value pair, we also maintain a map of `word` -> "key store" so that key is added to the "key store" of each word. For repeat key case, we can **delete** this key from "key store" corresponding to words present in old value, and populate this key corresponding to words from new value. We've this abstract "key store" but we have not really gave a thought about which data structure to use so that we can a) add key to a word b) remove key from a word and c) return latest K keys from the word.

Since we need to have an order in "key store", instead of having keys itself (which can come in any order), we can use some global `Counter` which increases every time a new key is added. We can initialize Counter = 0 and this can be used to find index as well: index = Counter%size. Storing increasing Counter sounds like having an array which will always be sorted, so searching(O(limit)) and adding(O(1)) will be optimal; but **deleting** will be time consuming (O(log(size) + size)). To make this efficient as well, we can use balanced BST, and now adding and deleting Counters will both take O(log(size)). While searching, once we have an array of Counters, we can use Counter -> Key map (which can be created) to fetch keys.

One more optimization: instead of using word -> BST(key store) hashmap, we can use trie where each node in trie will have a balanced BST(Red-Black tree). Searching for word through trie will take constant time, as word has maximum length of 15 chars.

Above, we noticed that we used Counter to store in BST. This Counter can be also used to fetch index and we can have Key -> Counter map instead of Key -> Index map, and similarly we have Counter -> Key map.

Instead of frontIndex and rearIndex, we can now have StartCounter and EndCounter to track front and rear of circular queue of values.

Below is the LogStore struct attributes:
```
type LogStore struct {
    Size                     int
    Values                   []string // this acts as circular queue for values
    StartCounter, EndCounter int64    // start and end counter between which all records are live. This also act as front and rear of circular queue
    KeyCounterMap            map[int64]int64
    CounterKeyMap            map[int64]int64
    WordTrie                 *trie.Trie // Trie which will have balanced BST at each node.
}
```
For balanced BST, I've used TreeSet: https://github.com/emirpasic/gods/blob/master/sets/treeset/ which is implemented with the help of Red-Black tree under the hood.


<ins>Time Complexity</ins>:

1. `Add(Key int64, Value string)`: `O(1)` (for inserting value in circular queue) + `O(15*log(S))` (for adding Counter/Key in each word's BST) + `O(15*log(S))` (for deleting Counter in each word's BST in case of existing key) where S is the max log size
2. `Search(Word string, limit int)`: `O(log(S))` for reaching End iterator of treeset + `O(limit)` for traversing through the treeset(RBT).

Note that I've assumed golang's map to have O(1) time complexity for inserting and deletion

<ins>Space Complexity</ins>:

O(S) -> storing values in circular queue made out of array

O(S) -> storing KeyCounterMap and CounterKeyMap

O(15\*S) -> storing BST in trie node with some counters. 15\*S beacase there can be 15 words and each word can be present in each value strings.
