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
  MenuStyle string
  Lines []*Line
}

func (ss *Sect) inspect() {
  fmt.Printf("%s { %s, %s }\n", ss.Title, ss.MainPic, ss.MenuPic)
  for _, l := range ss.Lines {
    l.inspect()
  }
}

func readSections() ([]Sect) {
  sectDir := "sections"
  files, _ := ioutil.ReadDir(sectDir)
  var sections []Sect
  sections = make([]Sect, len(files))
  for i, f := range files {
    j, _ := ioutil.ReadFile(path.Join(sectDir, f.Name()))
    ss := Sect{}
    json.Unmarshal(j, &ss)
    sections[i] = ss
  }
  return sections
}
