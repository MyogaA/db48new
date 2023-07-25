package main

import (
	"backend/connection"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Id          int
	Name        string
	Description string
	Start_date  time.Time
	End_date    time.Time
	Image       string
}
type User struct {
	Id         int
	Name       string
	Email      string
	Password   string
	Pengalaman []string
	Tahun      []string
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

	e.Static("/assets", "assets")

	e.GET("/hello", helloWorld)
	e.GET("/about", aboutMe)
	e.GET("/home", home)
	e.GET("/contact", contact)
	e.GET("/myproject", myproject)
	e.GET("/form-project", formProject)
	e.GET("/testimonial", testimonial)
	e.GET("/detail/:id", detail)
	e.POST("/add-project", addProject)
	e.POST("/delete-project/:id", deleteProject)
	e.GET("/edit/:id", edit)
	e.POST("/update-project/:id", updateProject)

	e.Logger.Fatal(e.Start("localhost:5003"))
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

	// id nantinya dapet dari user login

	if errQuery != nil {
		fmt.Println("masuk sini")
		return c.JSON(http.StatusInternalServerError, errQuery.Error())
	}

	dataReponse := map[string]interface{}{
		"User": dataUser,
	}
	return tmpl.Execute(c.Response(), dataReponse)
}

func myproject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/myproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
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
		"Project": resultProject,
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
	start_date := c.FormValue("startDate")
	end_date := c.FormValue("endDate")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project ( name, description,image, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", name, description, "default.jpg", start_date, end_date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// newProject := Project{
	// 	Name:        Name,
	// 	Description: Description,
	// 	Start_date:  time.Now(),
	// 	End_date:    time.Now(),
	// }
	// dataProject = append(dataProject, newProject)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
func deleteProject(c echo.Context) error {
	id := c.Param("id")

	idToInt, _ := strconv.Atoi(id)
	fmt.Println("persiapan delete index:", id)
	// dataProject = append(dataProject[:idToInt], dataProject[idToInt+1:]...)
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

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, image, start_date,end_date FROM tb_project WHERE id=$1", idToInt).Scan(&detail.Id, &detail.Name, &detail.Description, &detail.Image, &detail.Start_date, &detail.End_date)
	fmt.Println("ini data blog detail: ", errQuery)
	data := map[string]interface{}{
		"Id":      id,
		"Project": detail,
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
		"Id":      id,
		"Project": project,
	}
	return tmpl.Execute(c.Response(), data)
}
func updateProject(c echo.Context) error {
	id := c.Param("id")

	name := c.FormValue("title")
	description := c.FormValue("content")

	idToInt, _ := strconv.Atoi(id)
	fmt.Println(idToInt)
	dataUpdate, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name=$1, description=$2 WHERE id=$3", name, description, id)
	if err != nil {
		fmt.Println("error guys", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(dataUpdate)
	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
