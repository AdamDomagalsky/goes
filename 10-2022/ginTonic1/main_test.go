package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWithoutHTTP(t *testing.T) {
	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user1": "adent",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	res := rec.Result()

	// Assertions stuff here
	var p interface{}
	err := json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		panic(err)
	}
	m := map[string]string{"user1": "adent"}
	eq := fmt.Sprint(p) == fmt.Sprint(m) // easy way to compare 2 maps - here map[string]string vs interface | map[string]interfrace{}
	if !eq {
		t.Error("Expect this BS to be the same")
	}
}

func TestWithHTTP(t *testing.T) {
	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user1": "adent",
		})
	})

	srv := httptest.NewServer(router)
	defer srv.Close()
	client := srv.Client()

	res, err := client.Get(srv.URL + "/users")
	if err != nil {
		t.Error(err)
	}
	// Assertions stuff here
	var p interface{}
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		panic(err)
	}
	m := map[string]string{"user1": "adent"}
	eq := fmt.Sprint(p) == fmt.Sprint(m) // easy way to compare 2 maps - here map[string]string vs interface | map[string]interfrace{}
	if !eq {
		t.Error("Expect this BS to be the same")
	}

}

func TestApiEmployees(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/employees/", nil)
	r := gin.New()
	registerRoutes(r)
	r.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Fail()
	}
	t.Log(rec.Body.String())
}
