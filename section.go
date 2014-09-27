package main

import (
  "fmt"
  "io/ioutil"
  "path"
  "encoding/json"
)

type Line struct {
  Title string
  Href string
  Id string
  Class string
}

func (l *Line) inspect() {
  fmt.Printf(" %s { %s }", l.Title, l.Href)
  if l.Id != "" {
    fmt.Printf("#%s", l.Id)
  }
  if l.Class != "" {
    fmt.Printf(".%s", l.Class)
  }
  fmt.Println()
}

type Sect struct {
  Title string
  MainPic string
  MenuPic string
  Lines []*Line
}

func (ss *Sect) inspect() {
  fmt.Printf("%s { %s, %s }\n", ss.Title, ss.MainPic, ss.MenuPic)
  for i := range ss.Lines {
    ss.Lines[i].inspect()
  }
}

const sectDir = "./sections"
func readSections() ([]Sect) {
  files, _ := ioutil.ReadDir(sectDir)
  var sections []Sect
  for i := range files {
    j, _ := ioutil.ReadFile(path.Join(sectDir, files[i].Name()))
    ss := &Sect{}
    json.Unmarshal(j, ss)
    ss.inspect()
  }
  return sections
}
