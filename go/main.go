package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/user"

	"golang.design/x/clipboard"
)

// get current user home directory
func getUserHomeDir() string {
  usr, err := user.Current()
  if err != nil {
    return ""
  }
  return usr.HomeDir
}

func getPicturePath() string {
  homeDir := getUserHomeDir()
  return (homeDir + "/Pictures")
}

// generate 12 length random string 
func generateRandomString() string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, 12)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

func saveFile(path string, req *http.Request ) {
  f, err := os.Create(path) 
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  _, err = io.Copy(f, req.Body)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("image downloaded")
  
}

func uploads(w http.ResponseWriter, req *http.Request) {
  picDir := getPicturePath()
  fileName := generateRandomString()
  picPath :=  picDir + "/" + fileName
  fmt.Fprintf(w, "uploads\n" + picPath)
  saveFile(picPath, req)
}

type ShareClipboardReq struct {
  Data string 
}

func shareclipboard(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "shareclipboard\n")
  var s ShareClipboardReq 
  err := json.NewDecoder(req.Body).Decode(&s)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  } 
  clipboard.Write(clipboard.FmtText, []byte(s.Data))
}


func main() {
  port := "45531"
	http.HandleFunc("/uploads", uploads)
	http.HandleFunc("/shareclipboard", shareclipboard)

	fmt.Println("run server at ", port)
	http.ListenAndServe(":"+port, nil)
}
