package sqlite

import (
	"database/sql"

	"github.com/ankit8125/crud-golang-practice/internal/config"
	// "github.com/ankit8125/crud-golang-practice/internal/storage"
	_ "github.com/mattn/go-sqlite3" // This is being used indirectly, that's why we use '_'
)

type Sqlite struct {
	Db *sql.DB
}

// There is no concept of constructors in 'Go'
// So we make a function which makes the instance of the struct and return it
// Convention - use function name as 'New'
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath) // driver - sqlite3
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT, 
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error){
	
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}