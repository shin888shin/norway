package mydb

import (
  // "fmt"
  // "log"
  "io/ioutil"
  "strings"
)

type Database struct {
  User string
  Password string
  Host string
}

func DatabaseSetup() (Database, error) {
  var db Database
  fileName := "config/joel.txt"
  contentBytes, err := ioutil.ReadFile(fileName)

  if err != nil {
    return db, err
  }

  result := string(contentBytes)
  lines := strings.Split(result, "\n")

  for i, line := range lines {
    val := strings.TrimSpace(line)
    if i == 0 {
      db.User = val
    } else if i == 1 {
      db.Password = val
    } else {
      db.Host = val
    }

  }
  return db, nil
}