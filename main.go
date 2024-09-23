package main

import (
    "log"
    "net/http"
)


func main() {
    hub := newHub()
    go hub.run()

    http.Handle("/", http.FileServer(http.Dir("./web")))
    http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./public"))))
    http.HandleFunc("POST /upload", uploadSingleFile)
    http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request){
        serveWs(hub, w, r)
    })

    port := "12345"
    log.Printf("http://127.0.0.1:%v \n", port)
    if err := http.ListenAndServe(":" + port, nil); err != nil {
        log.Println(err)
    }
}

func init(){
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}