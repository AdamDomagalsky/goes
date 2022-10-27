package main

import (
	"embed"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"goes/10-2022/ginTonic1/employee"
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
	router.LoadHTMLGlob("./templates/*")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("future", ValidatorFuture)
	}
	//Static Routes
	router.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello world!")
	})
	//just one file http://localhost:3000/sf
	router.StaticFile("/sf", "./public/index.html")
	// whole directory http://localhost:3000/public
	router.Static("/public", "./public")
	// embedded into binary http://localhost:3000/sfs/public/app.css
	router.StaticFS("/sfs", http.FileSystem(http.FS(f)))

	router.GET("/employeeHtml", func(context *gin.Context) {
		context.File("./public/employee.html")
	})

	// http://localhost:3000/employeesOLD/ajane
	router.GET("/employeesOLD/:username", func(context *gin.Context) {
		username := context.Param("username")
		context.String(http.StatusOK, username)
	})

	// http://localhost:3000/employeesOLD/ajane/roles/42
	router.GET("/employeesOLD/:username/*rest", func(context *gin.Context) {
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

	//http://localhost:3000/apiOLD/whatever
	router.GET("/apiOLD/*rest", func(context *gin.Context) {
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
	router.POST("/employeeOLD", func(context *gin.Context) {
		var timeoffRequest TimeoffRequest
		if err := context.ShouldBind(&timeoffRequest); err == nil {
			context.JSON(http.StatusOK, timeoffRequest)
		} else {
			context.String(http.StatusInternalServerError, err.Error())
		}
	})

	// curl --location --request POST 'localhost:3000/api/timeoff' \ --header 'Content-Type: application/json' \ --data-raw '{ "date": "2023-12-31T00:00:00Z", "amount":-1 } '
	apiGroup := router.Group("/api")
	apiGroup.POST("/timeoff", func(context *gin.Context) {
		var timeoffRequest TimeoffRequest
		if err := context.ShouldBindJSON(&timeoffRequest); err == nil {
			context.JSON(http.StatusOK, timeoffRequest)
		} else {
			context.String(http.StatusInternalServerError, err.Error())
		}
	})

	// download
	// http://localhost:3000/tales
	router.StaticFile("/tales", "./tales/index.html")
	router.GET("/tale_of_two_cities", func(context *gin.Context) {
		//context.File("./tales/a_tale_of_two_cities.txt")
		f, err := os.Open("./tales/a_tale_of_two_cities.txt")
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		data, err := io.ReadAll(f)
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		context.Data(http.StatusOK, "text/plain", data)
	})

	router.GET("/great_expectations", func(context *gin.Context) {
		talesDir := "tales"
		fileName := "great_expectations.txt"
		f, err := os.Open(filepath.Join(talesDir, fileName))
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		fi, err := os.Stat(f.Name())
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		buf := make([]byte, 512)
		_, err = f.Read(buf)
		f.Seek(0, 0)
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		contentType := http.DetectContentType(buf)
		context.DataFromReader(http.StatusOK,
			fi.Size(),
			contentType,
			f,
			map[string]string{
				"Content-Disposition": "attachment;filename=" + fileName,
				"Content-Type":        contentType,
				"Content-Length":      string(fi.Size()),
			})
	})

	// http://localhost:3000/great_expectations2 streams 4kb
	router.GET("/great_expectations2", func(context *gin.Context) {
		talesDir := "tales"
		fileName := "great_expectations.txt"
		f, err := os.Open(filepath.Join(talesDir, fileName))
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
		}
		context.Stream(streamer(f))
	})

	router.LoadHTMLGlob("./templates/*")
	router.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusTemporaryRedirect, "/empTmpl")
	})

	router.GET("empTmpl", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.tmpl", employee.GetAll())
	})

	router.POST("/employee", func(context *gin.Context) {
		context.Request.Method = http.MethodGet
		context.Redirect(http.StatusFound, "employee")
	})
	router.GET("/employees/:employeeID", func(context *gin.Context) {
		employeeIDRaw := context.Param("employeeID")
		if emp, ok := tryToGetEmployee(context, employeeIDRaw); ok {
			context.HTML(http.StatusOK, "employee.tmpl", emp)
		}
	})
	router.POST("/employees/:employeeID", func(context *gin.Context) {
		var timeoff employee.TimeOff
		err := context.ShouldBind(&timeoff)
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		timeoff.Type = employee.TimeoffTypePTO
		timeoff.Status = employee.TimeoffStatusRequested
		employeeIDRaw := context.Param("employeeID")
		if employee, ok := tryToGetEmployee(context, employeeIDRaw); ok {
			employee.TimeOff = append(employee.TimeOff, timeoff)
			context.Redirect(http.StatusFound, "/employees/"+employeeIDRaw)
		}
	})

	empApiG := apiGroup.Group("employees")
	empApiG.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, employee.GetAll())
	})
	empApiG.GET("/:employeeID", func(context *gin.Context) {
		employeeIDRaw := context.Param("employeeID")
		if emp, ok := tryToGetEmployee(context, employeeIDRaw); ok {
			context.JSON(http.StatusOK, *emp)
		}
	})
	empApiG.POST("/:employeeID", func(context *gin.Context) {
		var timeoff employee.TimeOff
		err := context.ShouldBind(&timeoff)
		if err != nil {
			context.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		timeoff.Type = employee.TimeoffTypePTO
		timeoff.Status = employee.TimeoffStatusRequested
		employeeIDRaw := context.Param("employeeID")
		if emp, ok := tryToGetEmployee(context, employeeIDRaw); ok {
			emp.TimeOff = append(emp.TimeOff, timeoff)
			context.JSON(http.StatusOK, *emp)
		}
	})
	// to generate JSON msg on the fly use gin.H
	h := gin.H{"field1": "value2", "field2": 123} // type of a map of string to the empty interface
	println(h)
	// starting server & fatal if fail
	log.Fatal(router.Run(":3000"))
}

func streamer(r io.Reader) func(writer io.Writer) bool {
	return func(step io.Writer) bool {
		for {
			data := make([]byte, 4*2^10)
			if _, err := r.Read(data); err == nil {
				_, err := step.Write(data)
				return err == nil
			} else {
				return false
			}
		}
	}
}
func tryToGetEmployee(ctx *gin.Context, employeeIDRaw string) (*employee.Employee, bool) {
	id, err := strconv.Atoi(employeeIDRaw)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("lol"))
		return nil, false
	}
	emp, err := employee.Get(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return nil, false
	}
	return emp, true
}
