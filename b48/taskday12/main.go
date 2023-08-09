package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Title    string
	Content  string
	Author   string
	PostDate string
}

var dataProject = []Project{
	{
		Title:    "Title 1",
		Content:  "Content 1",
		Author:   "Surya Elidanto",
		PostDate: "20/07/2023",
	},
	{
		Title:    "Title 2",
		Content:  "Content 2",
		Author:   "Angga Nur",
		PostDate: "21/07/2023",
	},
}

func main() {
	e := echo.New()

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

	e.Logger.Fatal(e.Start("localhost:5008"))
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

	return tmpl.Execute(c.Response(), nil)
}

func myproject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/myproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	data := map[string]interface{}{
		"Project": dataProject,
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
	title := c.FormValue("title")
	content := c.FormValue("content")

	newProject := Project{
		Title:    title,
		Content:  content,
		Author:   "Yoga",
		PostDate: "21/07/2023",
	}
	dataProject = append(dataProject, newProject)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
func deleteProject(c echo.Context) error {
	id := c.Param("id")

	idToInt, _ := strconv.Atoi(id)
	fmt.Println("persiapan delete index:", id)
	dataProject = append(dataProject[:idToInt], dataProject[idToInt+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
func detail(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	idToint, _ := strconv.Atoi(id)

	detail := Project{}

	for index, data := range dataProject {
		if index == idToint {
			detail = Project{
				Title:    data.Title,
				Author:   data.Author,
				Content:  data.Content,
				PostDate: data.PostDate,
			}
		}
	}
	data := map[string]interface{}{
		"Id":      id,
		"Project": detail,
	}
	return tmpl.Execute(c.Response(), data)
}
func edit(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/edit.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	idToInt, _ := strconv.Atoi(id)
	project := Project{}

	for index, data := range dataProject {
		if index == idToInt {
			project = Project{
				Title:    data.Title,
				Author:   data.Author,
				Content:  data.Content,
				PostDate: data.PostDate,
			}
		}
	}

	data := map[string]interface{}{
		"Id":      id,
		"Project": project,
	}
	return tmpl.Execute(c.Response(), data)
}
func updateProject(c echo.Context) error {
	id := c.Param("id")

	title := c.FormValue("title")
	content := c.FormValue("content")

	idToInt, _ := strconv.Atoi(id)
	dataProject[idToInt].Title = title
	dataProject[idToInt].Content = content

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
