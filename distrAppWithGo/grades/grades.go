package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s *Student) Average() float32 {
	var sum float32
	for _, grade := range s.Grades {
		sum += grade.Score
	}

	return sum / float32(len(s.Grades))
}

type Students []Student

var (
	students      Students // in-memory mock
	studentsMutex sync.Mutex
)

func (s *Students) GetByID(id int) (*Student, error) {
	for _, student := range students {
		if student.ID == id {
			return &student, nil
		}
	}
	return nil, fmt.Errorf("student[%v] not found", id)
}

type GradeType string

const (
	GradeTest     = GradeType("Test")
	GradeHomework = GradeType("Homework")
	GradeQuiz     = GradeType("Quiz")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
