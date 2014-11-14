package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "code.google.com/p/go.net/websocket"
    "os/exec"
    "strings"
    "regexp"
)

type Memory struct {
    Total string
    Used string
}

func Free(ws *websocket.Conn) {
    for true {
        out, err := exec.Command("free", "-m").Output()
        if err != nil {
            fmt.Println("Error")
            return
        }
        memory := ParseOutput(out)
        err = websocket.JSON.Send(ws, memory)
        if err != nil {
            fmt.Println("Can't send")
        }
        time.Sleep(1000 * time.Millisecond)
    }
}

func ParseOutput(out []byte) Memory {
    tokenizedOutput := strings.Split(string(out), "\n")
    re := regexp.MustCompile(`Mem:\s+(?P<total>\d+)\s+(?P<used>\d+)`)
    matches := re.FindStringSubmatch(tokenizedOutput[1])
    memory := Memory{matches[1], matches[2]}
    return memory
}

func main() {
    http.Handle("/", websocket.Handler(Free))
    err := http.ListenAndServe(":12345", nil)
    CheckError(err)
}

func CheckError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)        
    }
}
