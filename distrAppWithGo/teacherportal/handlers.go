package teacherportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"goes/distrAppWithGo/grades"
	"goes/distrAppWithGo/registry"
)

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO handle via regExp please
	pathSegment := strings.Split(r.URL.Path, "/")
	switch len(pathSegment) {
	case 2:
		sh.renderStudents(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegment[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		sh.renderStudent(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegment[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		sh.renderGrades(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (sh studentsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", fmt.Sprintf("/students/%v", id))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
	title := r.FormValue("Title")
	gradeType := r.FormValue("Type")
	score, err := strconv.ParseFloat(r.FormValue("Score"), 32)
	if err != nil {
		log.Println("Failed to parse score: ", err)
		return
	}
	g := grades.Grade{
		Title: title,
		Type:  grades.GradeType(gradeType),
		Score: float32(score),
	}
	data, err := json.Marshal(g)
	if err != nil {
		log.Println("Failed to convert grade to JSON: ", g, err)
		return
	}
	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		log.Printf("Failed to retrieve instance of %v %v", registry.GradingService, err)
		return
	}
	res, err := http.Post(fmt.Sprintf("%v/students/%v/grades", serviceURL, id), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to save grade to GradingService")
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Failed to save grade to GS. Status:", res.StatusCode)
		return
	}
}

func (sh studentsHandler) renderStudent(w http.ResponseWriter, r *http.Request, id int) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving students: ", err)
		}
	}()
	serviceURL, err := registry.GetProvider(registry.GradingService) // using svc registry peer-to-peer connection
	if err != nil {
		return
	}
	res, err := http.Get(fmt.Sprintf("%v/students/%v", serviceURL, id))
	if err != nil {
		return
	}
	var s grades.Student
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("student.gohtml").Execute(w, s)
	if err != nil {
		return
	}
}

func (sh studentsHandler) renderStudents(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving students: ", err)
		}
	}()
	serviceURL, err := registry.GetProvider(registry.GradingService) // using svc registry peer-to-peer connection
	if err != nil {
		return
	}
	res, err := http.Get(serviceURL + "/students")
	if err != nil {
		return
	}
	var s grades.Students
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}
	err = rootTemplate.Lookup("students.gohtml").Execute(w, s)
	if err != nil {
		return
	}
}

func RegisterHandlers() {
	http.Handle("/", http.RedirectHandler("/students", http.StatusPermanentRedirect))
	h := new(studentsHandler)

	http.Handle("/students", h)
	http.Handle("/students/", h)
}
