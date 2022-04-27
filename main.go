package main

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/diabhiue/100ms/logs"
    "log"
    "os"
    "strconv"
    "strings"
)

func main() {
    // Initialize reader and read line by line
    reader := bufio.NewReader(os.Stdin)

    // Input max size of Log store
    scannedText, _ := reader.ReadString('\n')
    inputs := strings.Fields(scannedText)
    S, _ := strconv.Atoi(inputs[0])

    // Create Log object of max size S
    Log := logs.NewLogStore(S)

    for {
        // Read line from stdin and split it on whitespace
        scannedText, _ := reader.ReadString('\n')
        inputs := strings.Fields(scannedText)

        // Store the keyword to identify action to be done
        var token string
        token = inputs[0]

        // Based on input, add, search or end the process
        if token == "ADD" {
            var key int64
            var value string
            key, _ = strconv.ParseInt(inputs[1], 10, 64)
            value = strings.Join(inputs[2:], " ")
            Log.Add(key, value)
        } else if token == "SEARCH" {
            var word string
            var limit int
            word = inputs[1]
            limit, _ = strconv.Atoi(inputs[2])
            outputKeys := Log.Search(word, limit)
            if len(outputKeys) == 0 {
                fmt.Print("NONE")
            } else {
                for _, key := range outputKeys {
                    fmt.Print(key, " ")
                }
            }
            fmt.Println()
        } else if token == "END" {
            fmt.Println("END")
            break
        } else {
            log.Fatal(errors.New("Not a valid token"))
        }
    }
}
