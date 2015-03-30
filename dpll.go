package main

import (
    // "fmt"
    "log"
)

func BCP(f Formula, m Model) Formula {
    propagated := true
    for propagated {
        propagated = false
        clauseLoop: 
        for i, c_i := range f.Clauses {
            var lit Literal = 0
            for _, l := range c_i.Literals {
                if l != 0 && l != 1 {
                    if lit == 0 {
                        lit = l
                    } else {
                        continue clauseLoop
                    }
                }
            }
            if lit != 0 {
                // we found a unit clause
                found := false
                for j, c_j := range f.Clauses {
                    if i == j {
                        continue
                    }
                    for k, l := range c_j.Literals {
                        if l == -lit {
                            found = true
                            c_j.Literals[k] = 0
                        }
                    }
                    if found {
                        break
                    }
                }
                if found {
                    // found a resolvant; we can have another go
                    propagated = true
                    m[VarFromLiteral(lit)] = IsPositiveLiteral(lit)
                }
            }
        }
    }
    return f
}

func Solve(f Formula) (bool, Model) {
    m := make(map[Literal]bool)
    sat, m2 := DPLL(f, m)
    if sat {
        return true, m2
    } else {
        return false, nil
    }
}

func DPLL(f Formula, m Model) (bool, Model) {
    // log.Printf("DPLL: %v; model %v", f.Clauses, m)
    f = CopyFormula(f)
    f = BCP(f, m)
    if IsTrue(f) {
        // log.Printf("DPLL: %v is true; model %v", f.Clauses, m)
        return true, m
    } else if IsFalse(f) {
        // log.Printf("DPLL: %v is false; model %v", f.Clauses, m)
        return false, nil
    } else {
        // choose a variable
        v := Literal(0)
        varChooseLoop:
        for _, c := range f.Clauses {
            for _, l := range c.Literals {
                if l != 0 && l != 1 {
                    v = VarFromLiteral(l)
                    break varChooseLoop
                }
            }
        }

        if v == 0 {
            log.Fatal("no literals left", f)
        }

        f2 := CopyFormula(f)
        m2 := make(map[Literal]bool)
        for k, v := range m {
            m2[k] = v
        }

        m[v] = true
        m2[v] = false

        for i, c := range f.Clauses {
            for j, l := range c.Literals {
                if l == v {
                    f.Clauses[i].Literals[j] = 1
                    f2.Clauses[i].Literals[j] = 0
                } else if l == -v {
                    f.Clauses[i].Literals[j] = 0
                    f2.Clauses[i].Literals[j] = 1
                }
            }
        }

        // log.Printf("DPLL: branching on %v; f=%v, f2=%v", v, f.Clauses, f2.Clauses)

        if r, mr := DPLL(f, m); r {
            return true, mr
        } else {
            return DPLL(f2, m2)
        }
    }
}