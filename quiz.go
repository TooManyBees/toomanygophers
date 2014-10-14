package main

import (
  _ "github.com/lib/pq"
  "database/sql"
  "fmt"
  "math/rand"
  "time"
  "strings"
)

type Quiz struct {
  Id string
  Comics string
}

type Comic struct {
  Id int
  Title string
  Image string
}

type ComicStore struct {
  db *sql.DB
  count int
  dirty bool
  rand *rand.Rand
}

func (cs *ComicStore) init() {
  db, err := sql.Open("postgres", "dbname=spam_quiz sslmode=disable")
  if err != nil {
    fmt.Printf("Couldn't connect to database: %s\n", err)
    return
  }
  cs.dirty = true
  cs.db = db
  cs.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (cs *ComicStore) Count() int {
  if cs.dirty {
    row := cs.db.QueryRow("SELECT count(id) FROM comics")
    row.Scan(&cs.count)
    cs.dirty = false
  }
  return cs.count
}

// Return a comma-separated string of numbers from a slice of ints
// For use in a SQL 'IN' statement, because the sql library is crap
func sqlIn(arr []int) string {
  raw := fmt.Sprint(arr)
  replaced := strings.Replace(raw, " ", ",", -1)
  return replaced[1:len(replaced)-1]
}

// Return a random selection of n comics from the Comics table
// If n >= the number of saved comics, just return the table in original order
func (cs *ComicStore) random(n int) []Comic {
  ids := cs.rand.Perm(int(cs.Count()))
  returnAll := n >= cs.Count()

  var query string
  if returnAll {
    query = fmt.Sprintf("SELECT * FROM comics")
  } else {
    query = fmt.Sprintf("SELECT * FROM comics WHERE id IN (%s)", sqlIn(ids[:n]))
  }

  rows, err := cs.db.Query(query)
  if err != nil {
    fmt.Println(err)
  }

  var comics = make([]Comic, 0, n)

  for rows.Next() {
    var id int
    var title string
    var image string
    var alt_image string
    rows.Scan(&id, &title, &image, &alt_image)
    comics = append(comics, Comic{id, title, image})
  }

  return comics
}
