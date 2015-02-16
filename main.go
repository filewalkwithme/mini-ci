package main

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "os"
)

func main(){
  serverCI := http.NewServeMux()
  serverCI.HandleFunc("/", postReceive)

  http.ListenAndServe(":3000", serverCI)
}

func postReceive(w http.ResponseWriter, r *http.Request) {
    x, _ := ioutil.ReadAll(r.Body)

    var out bytes.Buffer
    json.Indent(&out, x, "", "\t")
    out.WriteTo(os.Stdout)
}

