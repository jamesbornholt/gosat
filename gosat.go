package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("usage: gosat input-file")
    }

    f := ParseFormulaFromFile(os.Args[1])

    // fmt.Printf("%v\n", f)

    b, m := Solve(f)

    fmt.Printf("%v %v\n", b, m)
}