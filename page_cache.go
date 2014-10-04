package main

import (
  "time"
  "os"
  "fmt"
  "path/filepath"
  "html/template"
)

const templateDir = "templates"

type PageRecord struct {
  Body *template.Template
  Modified time.Time
}

type PageCache map[string]PageRecord
var pageCache = PageCache{}

func partialsToFilenames(partials []string) []string {
  mapped := make([]string, len(partials))
  for i, partial := range partials {
  mapped[i] = filepath.Join(templateDir, partial+".html")
  }
  return mapped
}

func (p PageCache) get(page string, partials ...string) *template.Template {
  needsUpdating := false
  filenames := partialsToFilenames(partials)

  for _, filename := range filenames {
    fs, _ := os.Stat(filename)
    oldRecord, cached := p[page]
    if !cached || fs.ModTime().After(oldRecord.Modified)  {
      needsUpdating = true
      break
    }
  }
  if needsUpdating {
    templ := template.Must(template.ParseFiles(filenames...))
    p[page] = PageRecord{templ, time.Now()}
    fmt.Printf("  miss (%s)", page)
  } else {
    fmt.Printf("  hit (%s)", page)
  }
  return p[page].Body
}
