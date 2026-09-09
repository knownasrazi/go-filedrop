package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

type File struct {
    Data []byte
    Expires time.Time
}

var store = sync.Map{}

func upload(w http.ResponseWriter, r *http.Request) {
    // Hand-crafted: limit to 10MB, like a person would
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
    data, _ := io.ReadAll(r.Body)
    id := fmt.Sprintf("%d", time.Now().UnixNano() % 100000)
    store.Store(id, File{Data: data, Expires: time.Now().Add(24 * time.Hour)})
    fmt.Fprintf(w, "filedrop/%s", id)
}

func download(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/filedrop/"):]
    v, ok := store.Load(id)
    if !ok { http.NotFound(w, r); return }
    f := v.(File)
    if time.Now().After(f.Expires) { store.Delete(id); http.NotFound(w, r); return }
    w.Write(f.Data)
}

func main() {
    http.HandleFunc("/upload", upload)
    http.HandleFunc("/filedrop/", download)
    fmt.Println("go-filedrop listening on :8080")
    http.ListenAndServe(":8080", nil)
}
