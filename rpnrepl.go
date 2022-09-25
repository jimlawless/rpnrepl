// Copyright 2013 - by Jim Lawless
// License: MIT / X11
// See: http://www.mailsend-online.com/license2013.php
//
// Bear with me ... I'm a Go noob.
 
package main
 
import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)
 
const smax = 20
 
var stack []interface{}
var sptr int
var words map[string]func()
 
var patDquote = "[\"][^\"]*[\"]"
var patSquote = "['][^']*[']"
var patNumber = "[-]?\\d+"
var patWord = "\\S+"
 
var reDquote, reSquote *regexp.Regexp
var reNumber, reAll *regexp.Regexp
 
func push(val interface{}) {
    if sptr == 0 {
        log.Fatal("Stack overflow.")
    }
    sptr--
    stack[sptr] = val
}
 
func pop() interface{} {
    if sptr >= smax {
        log.Fatal("Stack underflow")
    }
    tmp := stack[sptr]
    sptr++
    return tmp
}
 
func add() {
    i, j := popTwo()
    push(i + j)
}
 
func sub() {
    i, j := popTwo()
    push(i - j)
}
 
func mul() {
    i, j := popTwo()
    push(i * j)
}
 
func div() {
    i, j := popTwo()
    push(i / j)
}
 
func dot() {
    fmt.Printf("%v", pop())
}
 
func cr() {
    fmt.Println()
}
 
func bye() {
    os.Exit(0)
}
 
func popTwo() (int, int) {
    j := pop().(int)
    i := pop().(int)
    return i, j
}
 
func repl() {
    var s string
    var k int
    var f func()
 
    reader := bufio.NewReader(os.Stdin)
    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        tokens := reAll.FindAllString(input, -1)
        for i := 0; i < len(tokens); i++ {
            s = tokens[i]
            switch {
            case reDquote.MatchString(s):
                s = s[1 : len(s)-1]
                push(s)
            case reSquote.MatchString(s):
                s = s[1 : len(s)-1]
                push(s)
            case reNumber.MatchString(s):
                k, err = strconv.Atoi(s)
                if err != nil {
                    log.Fatal(err)
                }
                push(k)
 
            default:
                f = words[s]
                if f != nil {
                    f()
                } else {
                    log.Fatal("Word '" + s + "' not found!")
                }
            }
        }
    }
}
 
func main() {
    sptr = smax
    stack = make([]interface{}, sptr+1)
    words = make(map[string]func())
    words["+"] = add
    words["-"] = sub
    words["*"] = mul
    words["/"] = div
    words["."] = dot
    words["cr"] = cr
    words["bye"] = bye
 
    reDquote, _ = regexp.Compile(patDquote)
    reSquote, _ = regexp.Compile(patSquote)
    reNumber, _ = regexp.Compile(patNumber)
    reAll, _ = regexp.Compile(patDquote + "|" + patSquote + "|" + patNumber + "|" + patWord)
    repl()
}
