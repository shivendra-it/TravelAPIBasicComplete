package main

import(
"log"
"net/http"
)

func DoCronJob(t int){
http.HandleFunc("/", handler)
err := http.ListenAndServe(":8080", nil)
if err != nil {
  log.Fatal("ListenAndServe: ", err)
}
}
