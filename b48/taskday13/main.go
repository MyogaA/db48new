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
	Name        string
	Description string
	Start_date  time.Time
	End_date    time.Time
	Image       string
}

var dataProject = []Project{
	{
		Name:        "Title 1",
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
	// e.GET("/detail/:id", detail)
	e.POST("/add-project", addProject)
	// e.POST("/delete-project/:id", deleteProject)
	// e.GET("/edit/:id", edit)
	// e.POST("/update-project/:id", updateProject)

	e.Logger.Fatal(e.Start("localhost:5004"))
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
	dataProject, errBlogs := connection.Conn.Query(context.Background(), "SELECT name, description, start_date, end_date, image  FROM tb_project")

	if errBlogs != nil {
		return c.JSON(http.StatusInternalServerError, errBlogs.Error())
	}

	var resultProject []Project
	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Name, &each.Description, &each.Start_date, &each.End_date, &each.Image)
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
	Name := c.FormValue("title")
	Description := c.FormValue("content")

	newProject := Project{
		Name:        Name,
		Description: Description,
		Start_date:  time.Now(),
		End_date:    time.Now(),
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

// func detail(c echo.Context) error {
// 	id := c.Param("id")

// 	tmpl, err := template.ParseFiles("views/detail.html")

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	idToint, _ := strconv.Atoi(id)

// 	detail := Project{}

// 	for index, data := range dataProject {
// 		if index == idToint {
// 			detail = Project{
// 				Title:    data.Title,
// 				Author:   data.Author,
// 				Content:  data.Content,
// 				PostDate: data.PostDate,
// 			}
// 		}
// 	}
// 	data := map[string]interface{}{
// 		"Id":      id,
// 		"Project": detail,
// 	}
// 	return tmpl.Execute(c.Response(), data)
// }
// func edit(c echo.Context) error {
// 	id := c.Param("id")

// 	tmpl, err := template.ParseFiles("views/edit.html")
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	idToInt, _ := strconv.Atoi(id)
// 	project := Project{}

// 	for index, data := range dataProject {
// 		if index == idToInt {
// 			project = Project{
// 				Title:    data.Title,
// 				Author:   data.Author,
// 				Content:  data.Content,
// 				PostDate: data.PostDate,
// 			}
// 		}
// 	}

// 	data := map[string]interface{}{
// 		"Id":      id,
// 		"Project": project,
// 	}
// 	return tmpl.Execute(c.Response(), data)
// }
// func updateProject(c echo.Context) error {
// 	id := c.Param("id")

// 	title := c.FormValue("title")
// 	content := c.FormValue("content")

// 	idToInt, _ := strconv.Atoi(id)
// 	dataProject[idToInt].Title = title
// 	dataProject[idToInt].Content = content

// 	return c.Redirect(http.StatusMovedPermanently, "/myproject")
// }
