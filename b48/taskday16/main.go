package main

import (
	"backend/connection"
	"backend/middleware"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id          int
	Name        string
	Description string
	Start_date  time.Time
	End_date    time.Time
	Image       string
	Author      string
}

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	Pengalaman     []string
	Tahun          []string
	IsAdmin        bool
}

var dataProject = []Project{
	{
		Name:        "Name 1",
		Description: "Content 1",
		Start_date:  time.Now(),
		End_date:    time.Now(),
		Image:       "defaul.jpg",
	},
	{
		Name:        "Title 1",
		Description: "Content 1",
		Start_date:  time.Now(),
		End_date:    time.Now(),
		Image:       "defaul.jpg",
	},
}

func main() {
	e := echo.New()

	connection.DatabaseConnect()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("suryaganteng"))))

	// Public routes
	e.GET("/hello", helloWorld)
	e.GET("/about", aboutMe)
	e.GET("/contact", contact)
	e.GET("/testimonial", testimonial)
	e.GET("/detail/:id", detail)

	// Routes related to projects
	e.GET("/home", home)
	e.GET("/myproject", myproject, isAuthenticated())
	e.GET("/form-project", formProject, isAuthenticated())
	e.POST("/add-project", middleware.UploadFile(addProject), isAuthenticated())
	e.POST("/delete-project/:id", deleteProject, isAuthenticated())
	e.GET("/edit/:id", edit, isAuthenticated())
	e.POST("/update-project/:id", updateProject, isAuthenticated())

	// Auth routes
	e.GET("/register", RegisterForm)
	e.POST("/register", register)
	e.GET("/login", LoginForm)
	e.POST("/login", login)
	e.POST("/logout", logout)

	e.Static("/assets", "assets")
	e.Static("/uploads", "uploads")

	e.Logger.Fatal(e.Start("localhost:5002"))
}

func helloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"name":    "yoga",
		"address": "ciantra",
	})
}

func aboutMe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Halo nama saya Yoga",
	})
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userId := 1
	var dataUser User
	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, pengalaman, tahun FROM tb_user WHERE id=$1", userId).Scan(&dataUser.Id, &dataUser.Name, &dataUser.Email, &dataUser.Pengalaman, &dataUser.Tahun)

	if errQuery != nil {
		fmt.Println("masuk sini")
		return c.JSON(http.StatusInternalServerError, errQuery.Error())
	}

	dataReponse := map[string]interface{}{
		"User":       dataUser,
		"IsLoggedIn": getUserIsLoggedIn(c),
	}
	return tmpl.Execute(c.Response(), dataReponse)
}

func myproject(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "myproject"
	tmpl, err := template.ParseFiles("views/myproject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mendapatkan data proyek-proyek dari database
	dataProject, errBlogs := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_user.name, tb_project.name, tb_project.description, tb_project.start_date, tb_project.image FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id")

	if errBlogs != nil {
		return c.JSON(http.StatusInternalServerError, errBlogs.Error())
	}

	var resultProject []Project
	for dataProject.Next() {
		var each = Project{}
		var tempAuthor sql.NullString
		err := dataProject.Scan(&each.Id, &tempAuthor, &each.Name, &each.Description, &each.Start_date, &each.Image)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		fmt.Println("ini datamu: ", tempAuthor.String)

		each.Author = tempAuthor.String
		resultProject = append(resultProject, each)
	}

	data := map[string]interface{}{
		"Project":    resultProject,
		"IsLoggedIn": getUserIsLoggedIn(c),
	}
	return tmpl.Execute(c.Response(), data)
}

func formProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/form-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

func testimonial(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/testimonial.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

func contact(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

func addProject(c echo.Context) error {
	name := c.FormValue("title")
	description := c.FormValue("content")
	startDate := c.FormValue("startDate")
	image := c.Get("dataFile").(string)

	sess, _ := session.Get("session", c)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (name, description, image, start_date, author_id) VALUES ($1, $2, $3, $4, $5)", name, description, image, startDate, sess.Values["id"].(int))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}

func deleteProject(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)
	fmt.Println("persiapan delete index:", id)
	connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", idToInt)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}

func detail(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)

	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	detail := Project{}
	var tempAuthor sql.NullString
	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, tb_user.name, tb_project.name, tb_project.description, tb_project.image, tb_project.start_date FROM tb_project JOIN tb_user ON tb_project.author_id = tb_user.id WHERE tb_project.id=$1", idToInt).Scan(&detail.Id, &tempAuthor, &detail.Name, &detail.Description, &detail.Image, &detail.Start_date)
	detail.Author = tempAuthor.String
	fmt.Println("ini data blog detail: ", errQuery)
	data := map[string]interface{}{
		"Id":         id,
		"Project":    detail,
		"IsLoggedIn": getUserIsLoggedIn(c),
	}
	return tmpl.Execute(c.Response(), data)
}

func edit(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)

	tmpl, err := template.ParseFiles("views/edit.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var project Project
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description FROM tb_project WHERE id=$1", idToInt).Scan(&project.Id, &project.Name, &project.Description)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	data := map[string]interface{}{
		"Id":         id,
		"Project":    project,
		"IsLoggedIn": getUserIsLoggedIn(c),
	}
	return tmpl.Execute(c.Response(), data)
}

func updateProject(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("title")
	description := c.FormValue("content")

	idToInt, _ := strconv.Atoi(id)
	fmt.Println(idToInt)
	dataUpdate, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name=$1, description=$2 WHERE id=$3", name, description, idToInt)
	if err != nil {
		fmt.Println("error guys", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(dataUpdate)
	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}

func RegisterForm(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	sess, errSess := session.Get("session", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],
		"IsLoggedIn":   getUserIsLoggedIn(c),
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}

func register(c echo.Context) error {
	inputName := c.FormValue("inputName")
	inputEmail := c.FormValue("inputEmail") // harus valid email
	inputPassword := c.FormValue("inputPassword")

	var existingEmail string
	err := connection.Conn.QueryRow(context.Background(), "SELECT email FROM tb_user WHERE email=$1", inputEmail).Scan(&existingEmail)
	if err == nil {
		return redirectWithMessage(c, "Email is already registered!", false, "/register")
	}

	var existingName string
	err = connection.Conn.QueryRow(context.Background(), "SELECT name FROM tb_user WHERE name=$1 AND email=$2", inputName, inputEmail).Scan(&existingName)
	if err == nil {
		return redirectWithMessage(c, "Combination of name and email is already registered!", false, "/register")
	}

	// validasi (trim, validasi valid email)
	fmt.Println("Register - Name:", inputName, "Email:", inputEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	query, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES($1, $2, $3)", inputName, inputEmail, hashedPassword)
	if err != nil {
		return redirectWithMessage(c, "Register failed!", false, "/register")
	}

	if query.RowsAffected() == 0 {
		return redirectWithMessage(c, "Register failed!", false, "/register")
	}

	return redirectWithMessage(c, "Register successful!", true, "/login")
}

func LoginForm(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	sess, errSess := session.Get("session", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],
		"IsLoggedIn":   getUserIsLoggedIn(c),
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	inputEmail := c.FormValue("inputEmail")
	inputPassword := c.FormValue("inputPassword")

	user := User{}

	// check apakah ada emailnya di db
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", inputEmail).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword)

	if err != nil {
		return redirectWithMessage(c, "Login gagal!", false, "/login")
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(inputPassword))

	if errPassword != nil {
		return redirectWithMessage(c, "Login gagal!", false, "/login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 hours
	sess.Values["id"] = user.Id
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["IsLoggedIn"] = true
	sess.Save(c.Request(), c.Response())

	return redirectWithMessage(c, "Login berhasil!", true, "/home")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, "/home")
}

func getUserIsLoggedIn(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	isLoggedIn := sess.Values["IsLoggedIn"]
	if isLoggedIn == nil {
		return false
	}
	return isLoggedIn.(bool)
}

func getUserIsAdmin(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	isAdmin := sess.Values["IsAdmin"]
	if isAdmin == nil {
		return false
	}
	return isAdmin.(bool)
}

func isAuthenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !getUserIsLoggedIn(c) {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			return next(c)
		}
	}
}

func redirectWithMessage(c echo.Context, message string, status bool, url string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, url)
}
