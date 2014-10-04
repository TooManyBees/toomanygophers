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

// .get is an abstraction, assuming that the page uses a
// template of the same name, as well as _base for the layout
func (p PageCache) get(pages ...string) *template.Template {
  pages = append(pages, "_base")
  return p.getTemplate(pages[0], pages...)
}

// Checks if a pageRecord is out of date.
// Returns true if any of the files in "filenames" are newer
// than the record's modified date.
func (p PageCache) isStale(page string, filenames ...string) bool {
  for _, filename := range filenames {
    fs, _ := os.Stat(filename)
    oldRecord, cached := p[page]
    if !cached || fs.ModTime().After(oldRecord.Modified)  {
      return true
    }
  }
  return false
}

// A get/set function that gets the template keyed under "page",
// loading/compiling/saving the template if it's missing or out
// of date.
func (p PageCache) getTemplate(page string, partials ...string) *template.Template {
  filenames := partialsToFilenames(partials)

  if p.isStale(page, filenames...) {
    templ := template.Must(template.ParseFiles(filenames...))
    p[page] = PageRecord{templ, time.Now()}
    fmt.Printf("  miss (%s)", page)
  } else {
    fmt.Printf("  hit (%s)", page)
  }
  return p[page].Body
}
