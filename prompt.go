package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func promptChoice(title string, choices ...string) (int, error) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Println(title)
        for i, choice := range choices {
            fmt.Printf("  %d) %s\n", i+1, choice)
        }
        fmt.Printf("Enter a number (1-%d): ", len(choices))

        input, err := reader.ReadString('\n')
        if err != nil {
            return 0, err
        }
        input = strings.TrimSpace(input)

        for i := range choices {
            if input == fmt.Sprintf("%d", i+1) {
                return i + 1, nil
            }
        }
        fmt.Println("Invalid choice. Try again.")
    }
}
