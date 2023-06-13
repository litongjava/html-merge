package main

import (
  "flag"
  "fmt"
  "github.com/fsnotify/fsnotify"
  "io/ioutil"
  "log"
  "os"
  "time"
)

func main() {
  b := flag.String("b", "body.html", "html body code file")
  c := flag.String("c", "index.css", "css code file")
  j := flag.String("j", "index.js", "java script code file")
  o := flag.String("o", "output.html", "output file")
  flag.Parse()
  log.Println("b:", *b)
  log.Println("c:", *c)
  log.Println("j:", *j)
  log.Println("o:", *o)
  start := time.Now().Unix()
  createFileIfNotExists(*b, *c, *j)
  listenAndMerge(b, c, j, o)
  end := time.Now().Unix()
  fmt.Println(end-start, "s")
  fmt.Println("done")
}

func createFileIfNotExists(filenames ...string) {
  for _, filePath := range filenames {
    // 检查文件是否存在
    _, err := os.Stat(filePath)
    if os.IsNotExist(err) {
      // 文件不存在，创建文件
      file, err := os.Create(filePath)
      if err != nil {
        fmt.Println("can't create file:", err, filePath)
        return
      }
      defer file.Close()

      fmt.Println("file created:", filePath)
    } else if err == nil {
      // 文件存在
      fmt.Println("file exists:", filePath)
    } else {
      // 其他错误
      fmt.Println("can't get file info:", err)
    }
  }
}

func addWatchFiles(watcher *fsnotify.Watcher, filenames ...string) error {
  for _, filePath := range filenames {
    if err := watcher.Add(filePath); err != nil {
      return err
    }
  }
  return nil
}

func listenAndMerge(bodyFilePath *string, cssFilePath *string, javaScriptFilePath *string, outputFilePath *string) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Println(err)
    return
  }
  defer watcher.Close()

  done := make(chan bool)
  listen(bodyFilePath, cssFilePath, javaScriptFilePath, outputFilePath, watcher)

  addWatchFiles(watcher, *bodyFilePath, *cssFilePath, *javaScriptFilePath)

  <-done
}

func listen(bodyFilePath *string, cssFilePath *string, javaScriptFilePath *string, outputFilePath *string, watcher *fsnotify.Watcher) {
  go func() {
    for {
      select {
      case event, ok := <-watcher.Events:
        if !ok {
          return
        }
        log.Println("event:", event)
        mergeFile(bodyFilePath, cssFilePath, javaScriptFilePath, outputFilePath)
      case err, ok := <-watcher.Errors:
        if !ok {
          return
        }
        log.Println("error:", err)
      }
    }
  }()
}

func mergeFile(bodyFilePath *string, cssFilePath *string, javaScriptFilePath *string, outputFilePath *string) {
  htmlContent, err := ioutil.ReadFile(*bodyFilePath)
  if err != nil {
    fmt.Println("can't read file:", *bodyFilePath, err)
  }

  //cssContent, err := ioutil.ReadFile(*cssFilePath)
  //if err != nil {
  //  fmt.Println("can't read file:", *cssFilePath, err)
  //}
  //
  //jsContent, err := ioutil.ReadFile(*javaScriptFilePath)
  //if err != nil {
  //  fmt.Println("can't read file:", *javaScriptFilePath, err)
  //}

  template := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Output</title>
  <link rel="stylesheet" href="%s">
</head>
<body>
%s
</body>
<script src="%s"></script>
</html>`
  output := fmt.Sprintf(template, *cssFilePath, htmlContent, *javaScriptFilePath)

  err = ioutil.WriteFile(*outputFilePath, []byte(output), 0644)
  if err != nil {
    fmt.Println("can't write output file:", err)
  }
}
