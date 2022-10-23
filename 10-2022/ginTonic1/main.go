package main

import (
	"embed"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//go:embed public/*
var f embed.FS

type User struct {
	FirstName string `json:"firstname" binding:"-"`
	LastName  string `json:"lastname" binding:"-"`
}

// TimeoffRequest binding:"-" means optional
type TimeoffRequest struct {
	Date   time.Time `json:"date" form:"date" binding:"required,future" time_format:"2006-01-02"`
	Amount float64   `json:"amount" form:"amount" binding:"required,gt=0"`
}

var ValidatorFuture validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		return date.After(time.Now())
	}
	return true
}

func main() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("future", ValidatorFuture)
	}
	//Static Routes
	router.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello world!")
	})
	//just one file http://localhost:3000/
	router.StaticFile("/", "./public/index.html")
	// whole directory http://localhost:3000/public
	router.Static("/public", "./public")
	// embedded into binary http://localhost:3000/fs/public/app.css
	router.StaticFS("/fs", http.FileSystem(http.FS(f)))

	router.GET("/employee", func(context *gin.Context) {
		context.File("./public/employee.html")
	})

	// http://localhost:3000/employees/ajane
	router.GET("/employees/:username", func(context *gin.Context) {
		username := context.Param("username")
		context.String(http.StatusOK, username)
	})

	// http://localhost:3000/employees/ajane/roles/42
	router.GET("/employees/:username/*rest", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"username": context.Param("username"),
			"rest":     context.Param("rest"),
		})
	})

	// Route Groups
	adminGroup := router.Group("/admin")
	adminGroup.GET("/users", func(context *gin.Context) {
		context.String(http.StatusOK, "Page to administer users")
	})
	adminGroup.GET("/roles", func(context *gin.Context) {
		context.String(http.StatusOK, "Page to administer roles")
	})
	adminGroup.GET("/policies", func(context *gin.Context) {
		context.String(http.StatusOK, "Page to administer polices")
	})

	//http://localhost:3000/api/whatever
	router.GET("/api/*rest", func(context *gin.Context) {
		url := context.Request.URL.String()
		headers := context.Request.Header
		cookies := context.Request.Cookies()
		context.IndentedJSON(http.StatusOK, gin.H{
			"url":     url,
			"headers": headers,
			"cookies": cookies,
		})
	})

	// http://localhost:3000/query/?username=arthur&year=2015&month=1&month=2
	router.GET("/query/*rest", func(context *gin.Context) {
		username := context.Query("username")
		year := context.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
		month := context.QueryArray("month")
		context.JSON(http.StatusOK, gin.H{
			"username": username,
			"year":     year,
			"month":    month,
		})
	})

	// http://localhost:3000/employeeOld
	router.POST("/employeeOld", func(context *gin.Context) {
		data := context.PostForm("date")
		amount := context.PostForm("amount")
		username := context.DefaultPostForm("username", "adent")

		context.IndentedJSON(http.StatusOK, gin.H{
			"data":     data,
			"amount":   amount,
			"username": username,
		})
	})

	// http://localhost:3000/employee
	router.POST("/employee", func(context *gin.Context) {
		var timeoffRequest TimeoffRequest
		if err := context.ShouldBind(&timeoffRequest); err == nil {
			context.JSON(http.StatusOK, timeoffRequest)
		} else {
			context.String(http.StatusInternalServerError, err.Error())
		}
	})

	apiGroup := router.Group("/api")
	apiGroup.POST("/timeoff", func(context *gin.Context) {
		var timeoffRequest TimeoffRequest
		if err := context.ShouldBindJSON(&timeoffRequest); err == nil {
			context.JSON(http.StatusOK, timeoffRequest)
		} else {
			context.String(http.StatusInternalServerError, err.Error())
		}
	})

	// starting server & fatal if fail
	log.Fatal(router.Run(":3000"))
}
