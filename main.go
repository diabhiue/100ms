package main

import (
	"fmt"
	"log"
	"sort"
	logs "packages/logs"
)

func main() {
	var S int
	fmt.Scan(&S)

	Log = logs.NewLogStore(S)

	for {
		var token string
		fmt.Scan(&token)

		if token == "ADD" {
			var key int
			var value string
			fmt.Scan(&key)
			fmt.Scan(&value)
			Log.add(key, value)
		} else if token == "SEARCH" {
			var word string
			var limit int
			fmt.Scan(&word, &limit)
			fmt.Println(Log.search(word))
		} else if token == "END" {
			break
		} else {
			log.Fatal(errors.New("Not a valid token"))
		}
		S--
	}
	fmt.Println(B)
	//logs = LogStore()
	arr := []int{11, 22, 23, 34}

	x := 32
	index := sort.SearchInts(arr, x)
	if index < len(arr) && arr[index] == x {
		fmt.Printf("So %d is present at %d\n", x, index)
	} else {
		fmt.Printf("Sorry, %d is not present in the array\n", x)
	}
}
