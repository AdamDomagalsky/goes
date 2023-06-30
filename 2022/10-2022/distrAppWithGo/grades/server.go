package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type studentsHandler struct {
}

// /students
// /students/{id}
// /students/{id}/grades
func (sh *studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO handle via regExp please
	pathSegment := strings.Split(r.URL.Path, "/")
	switch len(pathSegment) {
	case 2:
		sh.getAll(w)
	case 3:
		id, err := strconv.Atoi(pathSegment[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegment[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (sh *studentsHandler) getAll(w http.ResponseWriter) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh *studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Errorf("failed getByID student: %q", err))
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		log.Println(fmt.Errorf("failed to seralize studnets: %q", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh *studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(fmt.Errorf("failed getByID student: %q", err))
		return
	}

	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh *studentsHandler) toJSON(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)

	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func RegisterHandlers() {
	handler := new(studentsHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}
