package storage

import "github.com/ankit8125/crud-golang-practice/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}