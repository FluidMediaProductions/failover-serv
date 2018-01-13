package main

import (
  "html/template"
  "net/http"
  "strings"
  "log"
  "os"
  "io"
  "time"
  "strconv"
)

type Data struct {
  Year int
}

var (
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func Init(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
    Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func error(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  path = strings.TrimPrefix(path, "/")
  if path == "healthz" {
    w.Write([]byte("OK"))
    Info.Println("Health check OK")
    return
  }

  d := &Data{Year: time.Now().Year()}

  code := r.Header.Get("X-Code")
  if code == "" {
    code = "500"
  }

  codeInt := 500

  if code, err := strconv.ParseInt(code, 10, 32); err == nil {
    codeInt = int(code)
  }

  w.WriteHeader(codeInt)

  t, err := template.ParseFiles("public-html/index.html")
  if err != nil {
    Error.Fatal("Error processing template")
  }
  t.Execute(w, d)
  Warning.Printf("Serving %v code error page\n", codeInt)
}

func main() {
  Init(os.Stdout, os.Stdout, os.Stderr)

  http.HandleFunc("/", error)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
  }
}
