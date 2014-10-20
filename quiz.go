package main

import (
  _ "github.com/lib/pq"
  "database/sql"
  "fmt"
  "math/rand"
  "time"
  "strings"
  "strconv"
)

type Quiz struct {
  Id string
  Comics string
}

type Comic struct {
  Id int
  Title string
  Image string
  AltImage string
}

type ComicStore struct {
  db *sql.DB
  count int
  dirty bool
  rand *rand.Rand
}

type QuizResults struct {
  Answered int
  Unanswered int
  Total int
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

func fromRows(rows sql.Rows, callback func(int, string, string, string)) {

  for rows.Next() {
    var id int
    var title string
    var image string
    var altImage string
    rows.Scan(&id, &title, &image, &altImage)
    callback(id, title, image, altImage)
  }
}

func (cs *ComicStore) Find(ids []int) *sql.Rows {
  query := fmt.Sprintf("SELECT * FROM comics WHERE id IN (%s)", sqlIn(ids))
  rows, _ := cs.db.Query(query)
  return rows
}

func (cs *ComicStore) All() *sql.Rows {
  rows, _ := cs.db.Query("SELECT * FROM comics ORDER BY id")
  return rows
}

// Return a random selection of n comics from the Comics table
// If n >= the number of saved comics, just return the table in original order
func (cs *ComicStore) random(n int) []Comic {
  ids := cs.rand.Perm(int(cs.Count()))
  returnAll := n >= cs.Count()

  var rows sql.Rows
  if returnAll {
    rows = *cs.All()
  } else {
    rows = *cs.Find(ids[:n])
  }

  var comics = make([]Comic, 0, len(ids))
  fromRows(rows, func(id int, title string, image string, altImage string) {
    comics = append(comics, Comic{id, title, image, altImage})
  })
  return comics
}

func (cs *ComicStore) ParseAnswers(form map[string][]string) QuizResults {
  var ids []int
  var answers = make(map[int]string)
  var comics = make(map[int]string)
  for k, v := range form {
    if true {
      id, _ := strconv.ParseInt(k, 10, 32)
      ids = append(ids, int(id))
      answers[int(id)] = v[0]
    }
  }

  { // Make map(id => image) from database
    query := fmt.Sprintf("SELECT * FROM comics WHERE id IN (%s)", sqlIn(ids))
    rows, _ := cs.db.Query(query)
    fromRows(*rows, func(id int, title string, image string, altImage string) {
      comics[id] = image
    })
  }

  var unanswered = 0
  var answered = 0
  for _, id := range ids {
    if answers[id] == "" {
      unanswered += 1
    } else if answers[id] == comics[id] {
      answered += 1
    }
  }

  return QuizResults{answered, unanswered, len(ids)}
}
