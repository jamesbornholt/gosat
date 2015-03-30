package main

import (
    // "fmt"
    "log"
)

func BCP(f Formula, m Model) Formula {
    propagated := true  // did we do a unit propagation on this loop over the clauses?
    for propagated {
        propagated = false
        clauseLoop: 
        for i, c_i := range f.Clauses {
            // try to find a unit clause
            var lit Literal = 0
            for _, l := range c_i.Literals {
                if l != 0 && l != 1 {
                    if lit == 0 {  // first literal in this clause
                        lit = l
                    } else {  // more than one literal; not a unit clause
                        continue clauseLoop
                    }
                }
            }

            if lit == 0 {  // no literals in this clause
                continue clauseLoop
            }

            // c_i is a unit clause; now let's find a second clause to resolve
            found := false
            for j, c_j := range f.Clauses {
                if i == j {  // short circuit
                    continue
                }
                for k, l := range c_j.Literals {
                    if l == -lit {  // is l the negation of the unit clause's lit?
                        found = true
                        c_j.Literals[k] = 0
                    }
                    // TODO: do we need to also rewrite l == lit?
                }
                if found {  // found a second clause
                    break
                }
            }
            if found {  // made a unit resolution on lit, so keep going
                propagated = true
                m[VarFromLiteral(lit)] = IsPositiveLiteral(lit)
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

func Choose(f Formula) Literal {
    v := Literal(0)
    chooseLoop:
    for _, c := range f.Clauses {
        for _, l := range c.Literals {
            if l != 0 && l != 1 {
                v = VarFromLiteral(l)
                break chooseLoop
            }
        }
    }
    if v == 0 {
        log.Fatal("no literals left to choose!", f)
    }
    return v
}

func DPLL(f Formula, m Model) (bool, Model) {
    f = CopyFormula(f)

    // do BCP on f
    f = BCP(f, m)

    if IsTrue(f) {
        return true, m
    } else if IsFalse(f) {
        return false, nil
    }

    // we're going to have to branch. choose a variable to branch on.
    v := Choose(f)

    // create formula and model for the false branch
    f2 := CopyFormula(f)
    m2 := make(map[Literal]bool)
    for k, v := range m {
        m2[k] = v
    }

    // branched models
    m[v] = true
    m2[v] = false

    // rewrite the clauses with the branched values
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

    // perform DPLL on the true branch, then the false branch
    if r, mr := DPLL(f, m); r {
        return true, mr
    } else {
        return DPLL(f2, m2)
    }
}