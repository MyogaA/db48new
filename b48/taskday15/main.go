package main

import (
	"backend/connection"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Import package yang diperlukan

type Project struct {
	Id          int
	Name        string
	Description string
	Start_date  time.Time
	End_date    time.Time
	Image       string
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

// Data dummy untuk proyek-proyek
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
	// Membuat aplikasi Echo baru
	e := echo.New()

	// Koneksi ke database
	connection.DatabaseConnect()

	// Middleware untuk session menggunakan cookie
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("suryaganteng"))))

	// Route publik
	// Endpoint untuk mencetak "Hello World" dalam format JSON
	e.GET("/hello", helloWorld)
	e.GET("/about", aboutMe)
	e.GET("/contact", contact)
	e.GET("/testimonial", testimonial)
	e.GET("/detail/:id", detail)

	// Route terkait proyek
	e.GET("/home", home)
	e.GET("/myproject", myproject, isAuthenticated())
	e.GET("/form-project", formProject, isAuthenticated())
	e.POST("/add-project", addProject, isAuthenticated())
	e.POST("/delete-project/:id", deleteProject, isAuthenticated())
	e.GET("/edit/:id", edit, isAuthenticated())
	e.POST("/update-project/:id", updateProject, isAuthenticated())

	// Route terkait otentikasi (auth)
	e.GET("/register", RegisterForm)
	e.POST("/register", register)
	e.GET("/login", LoginForm)
	e.POST("/login", login)
	e.POST("/logout", logout)

	// Static file serving untuk folder "assets"
	e.Static("/assets", "assets")

	// Menjalankan server Echo pada alamat localhost:5001
	e.Logger.Fatal(e.Start("localhost:5001"))
}

// Fungsi untuk menghandle endpoint "/hello"
func helloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"name":    "yoga",
		"address": "ciantra",
	})
}

// Fungsi untuk menghandle endpoint "/about"
func aboutMe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Halo nama saya Yoga",
	})
}

// Fungsi untuk menghandle endpoint "/home"
func home(c echo.Context) error {
	// Menggunakan template HTML untuk halaman utama
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mendapatkan data user dari database berdasarkan ID tertentu
	userId := 1
	var dataUser User
	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, pengalaman, tahun FROM tb_user WHERE id=$1", userId).Scan(&dataUser.Id, &dataUser.Name, &dataUser.Email, &dataUser.Pengalaman, &dataUser.Tahun)

	// Menandai user sebagai admin berdasarkan email
	if strings.HasSuffix(dataUser.Email, "admin@gmail.com") {
		dataUser.IsAdmin = true
	} else {
		dataUser.IsAdmin = false
	}

	if errQuery != nil {
		fmt.Println("masuk sini")
		return c.JSON(http.StatusInternalServerError, errQuery.Error())
	}

	dataReponse := map[string]interface{}{
		"User":       dataUser,
		"IsLoggedIn": getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}
	return tmpl.Execute(c.Response(), dataReponse)
}

// Fungsi untuk menghandle endpoint "/myproject"
func myproject(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "myproject"
	tmpl, err := template.ParseFiles("views/myproject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mendapatkan data proyek-proyek dari database
	dataProject, errBlogs := connection.Conn.Query(context.Background(), "SELECT id, name, description, start_date, end_date, image  FROM tb_project")

	if errBlogs != nil {
		return c.JSON(http.StatusInternalServerError, errBlogs.Error())
	}

	var resultProject []Project
	for dataProject.Next() {
		var each = Project{}
		err := dataProject.Scan(&each.Id, &each.Name, &each.Description, &each.Start_date, &each.End_date, &each.Image)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		resultProject = append(resultProject, each)
	}

	data := map[string]interface{}{
		"Project":    resultProject,
		"IsAdmin":    getUserIsAdmin(c),
		"IsLoggedIn": getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}
	return tmpl.Execute(c.Response(), data)
}

// Fungsi untuk menghandle endpoint "/form-project"
func formProject(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "form-project"
	tmpl, err := template.ParseFiles("views/form-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

// Fungsi untuk menghandle endpoint "/testimonial"
func testimonial(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "testimonial"
	tmpl, err := template.ParseFiles("views/testimonial.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

// Fungsi untuk menghandle endpoint "/contact"
func contact(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "contact"
	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), nil)
}

// Fungsi untuk menambahkan proyek baru dengan metode POST
func addProject(c echo.Context) error {
	name := c.FormValue("title")
	description := c.FormValue("content")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")

	// Menyimpan data proyek baru ke database
	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (name, description, image, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", name, description, "default.jpg", startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}

// Fungsi untuk menghapus proyek dengan metode POST dan parameter ID
func deleteProject(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)
	fmt.Println("persiapan delete index:", id)
	connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", idToInt)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}

// Fungsi untuk menghandle endpoint "/detail/:id"
func detail(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)

	// Menggunakan template HTML untuk halaman "detail"
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	detail := Project{}

	// Mendapatkan data proyek berdasarkan ID tertentu
	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, image, start_date, end_date FROM tb_project WHERE id=$1", idToInt).Scan(&detail.Id, &detail.Name, &detail.Description, &detail.Image, &detail.Start_date, &detail.End_date)
	fmt.Println("ini data blog detail: ", errQuery)
	data := map[string]interface{}{
		"Id":         id,
		"Project":    detail,
		"IsLoggedIn": getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}
	return tmpl.Execute(c.Response(), data)
}

// Fungsi untuk menghandle endpoint "/edit/:id"
func edit(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)

	// Menggunakan template HTML untuk halaman "edit"
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
		"IsLoggedIn": getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}
	return tmpl.Execute(c.Response(), data)
}

// Fungsi untuk mengupdate proyek dengan metode POST dan parameter ID
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

// Fungsi untuk menghandle endpoint "/register"
func RegisterForm(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "register"
	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mendapatkan session
	sess, errSess := session.Get("session", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	// Menyimpan pesan flash dari session ke variabel
	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],
		"IsLoggedIn":   getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}

	// Menghapus pesan flash dari session setelah digunakan
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}

// Fungsi untuk melakukan registrasi dengan metode POST
func register(c echo.Context) error {
	inputName := c.FormValue("inputName")
	inputEmail := c.FormValue("inputEmail") // harus valid email
	inputPassword := c.FormValue("inputPassword")

	// Validasi apakah email sudah terdaftar di database
	var existingEmail string
	err := connection.Conn.QueryRow(context.Background(), "SELECT email FROM tb_user WHERE email=$1", inputEmail).Scan(&existingEmail)
	if err == nil {
		return redirectWithMessage(c, "Email sudah terdaftar!", false, "/register")
	}

	// Validasi apakah kombinasi nama dan email sudah terdaftar di database
	var existingName string
	err = connection.Conn.QueryRow(context.Background(), "SELECT name FROM tb_user WHERE name=$1 AND email=$2", inputName, inputEmail).Scan(&existingName)
	if err == nil {
		return redirectWithMessage(c, "Kombinasi nama dan email sudah terdaftar!", false, "/register")
	}

	// Validasi data (trim, validasi email yang valid)
	fmt.Println("Register - Name:", inputName, "Email:", inputEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Menyimpan data user baru ke database
	query, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES($1, $2, $3)", inputName, inputEmail, hashedPassword)
	if err != nil {
		return redirectWithMessage(c, "Registrasi gagal!", false, "/register")
	}

	// Memastikan ada baris yang terpengaruh untuk menandakan registrasi berhasil
	if query.RowsAffected() == 0 {
		return redirectWithMessage(c, "Registrasi gagal!", false, "/register")
	}

	return redirectWithMessage(c, "Registrasi berhasil!", true, "/login")
}

// Fungsi untuk menghandle endpoint "/login"
func LoginForm(c echo.Context) error {
	// Menggunakan template HTML untuk halaman "login"
	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mendapatkan session
	sess, errSess := session.Get("session", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	// Menyimpan pesan flash dari session ke variabel
	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],
		"IsLoggedIn":   getUserIsLoggedIn(c), // Menambahkan nilai IsLoggedIn ke data template.
	}

	// Menghapus pesan flash dari session setelah digunakan
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}

// Fungsi untuk melakukan login dengan metode POST
func login(c echo.Context) error {
	inputEmail := c.FormValue("inputEmail")
	inputPassword := c.FormValue("inputPassword")

	user := User{}

	// Memeriksa apakah email ada di database
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", inputEmail).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword)

	if err != nil {
		return redirectWithMessage(c, "Login gagal!", false, "/login")
	}

	// Memeriksa apakah password sesuai dengan yang ada di database
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(inputPassword))

	if errPassword != nil {
		return redirectWithMessage(c, "Login gagal!", false, "/login")
	}

	// Menyimpan data user ke dalam session setelah berhasil login
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 jam
	sess.Values["id"] = user.Id
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["IsLoggedIn"] = true
	sess.Save(c.Request(), c.Response())

	return redirectWithMessage(c, "Login berhasil!", true, "/home")
}

// Fungsi untuk melakukan logout dengan metode POST
func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Values["IsLoggedIn"] = false // Mengatur IsLoggedIn menjadi false saat logout
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, "/home")
}

// Fungsi untuk memeriksa apakah pengguna sudah login
func getUserIsLoggedIn(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	isLoggedIn := sess.Values["IsLoggedIn"]
	if isLoggedIn == nil {
		return false
	}
	return isLoggedIn.(bool)
}

// Fungsi untuk memeriksa apakah pengguna adalah admin
func getUserIsAdmin(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	isAdmin := sess.Values["IsAdmin"]
	if isAdmin == nil {
		return false
	}
	return isAdmin.(bool)
}

// Middleware untuk memeriksa apakah pengguna sudah terotentikasi (sudah login)
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

// Fungsi untuk menangani redirect dengan pesan flash
func redirectWithMessage(c echo.Context, message string, status bool, url string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, url)
}
