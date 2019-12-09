package main

import (
    "fmt"
    "reflect"
    "runtime"
)

func A() {
    pointers := make([]uintptr, 1)
    n := runtime.Callers(2, pointers)

    for i := 0; i < n; i++ {
        f := runtime.FuncForPC(pointers[i])

        fmt.Println(i, f.Name())
        fmt.Println(reflect.ValueOf(*f).Type().NumIn())
    }
}

func F(x string) string {

    A()
    return "x"
}

func main() {

    F("qwer")
}
