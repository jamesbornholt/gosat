package main

import (
    "bufio"
    "log"
    "os"
    "strings"
)

type Literal int

type Clause struct {
    Literals []Literal
}

type Formula struct {
    Clauses []Clause
    Variables map[string]int
    NextVariable int
}

type Model map[Literal]bool

func ParseFormulaFromFile(path string) Formula {
    f, err := os.Open(path)
    if err != nil {
        log.Fatal("couldn't open file: %s", path)
    }
    defer f.Close()

    scan := bufio.NewScanner(f)

    formula := Formula{
        Clauses: make([]Clause, 0), 
        Variables: make(map[string]int),
        NextVariable: 2,
    }

    for scan.Scan() {
        l := scan.Text()
        s := strings.Fields(l)

        c := Clause{Literals: make([]Literal, len(s))}

        for i,lit := range s {
            sgn := 1
            if lit[0] == '!' {
                sgn = -1
                lit = lit[1:]
            }
            v, ok := formula.Variables[lit]
            if !ok {
                v = formula.NextVariable
                formula.NextVariable++
                formula.Variables[lit] = v
            }

            c.Literals[i] = Literal(sgn*v)
        }

        formula.Clauses = append(formula.Clauses, c)
    }

    return formula
}

func CopyFormula(f Formula) Formula {
    f2 := Formula{
        Clauses: make([]Clause, len(f.Clauses)),
        Variables: make(map[string]int),
        NextVariable: f.NextVariable,
    }

    for i, c := range(f.Clauses) {
        f2.Clauses[i] = Clause{Literals: make([]Literal, len(c.Literals))}
        copy(f2.Clauses[i].Literals, c.Literals)
    }

    for k, v := range f.Variables {
        f2.Variables[k] = v
    }

    return f2
}

func IsTrue(f Formula) bool {
    for _, c := range f.Clauses {
        sat := false
        for _, l := range c.Literals {
            if l == 1 {
                sat = true
                break
            }
        }
        if !sat {
            return false
        }
    }
    return true
}

func IsFalse(f Formula) bool {
    for _, c := range f.Clauses {
        sat := false
        for _, l := range c.Literals {
            if l != 0 {
                sat = true
                break
            }
        }
        if !sat {
            return true
        }
    }
    return false
}

func VarFromLiteral(l Literal) Literal {
    if l < 0 {
        return -l
    } else {
        return l
    }
}

func IsPositiveLiteral(l Literal) bool {
    return l >= 0
}