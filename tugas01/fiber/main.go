package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        art := `
   _   _      _ _        __        __         _     _   _
  | | | | ___| | | ___   \ \      / /__  _ __| | __| | | |
  | |_| |/ _ \ | |/ _ \   \ \ /\ / / _ \| '__| |/ _` + "`" + ` | | |
  |  _  |  __/ | | (_) |   \ V  V / (_) | |  | | (_| | |_|
  |_| |_|\___|_|_|\___/     \_/\_/ \___/|_|  |_|\__,_| |_|
  
Hello, World!
`
        fmt.Fprint(w, art)
    })

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}