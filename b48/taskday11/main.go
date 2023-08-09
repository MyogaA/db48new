package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/assets", "assets")

	e.GET("/hello", helloWorld)
	e.GET("/about", aboutMe)
	e.GET("/home", home)
	e.GET("/contact", contact)
	e.GET("/myproject", myproject)
	e.GET("/testimonial", testimonial)
	e.GET("/detail/:id", detail)
	e.POST("/add-project", addProject)

	e.Logger.Fatal(e.Start("localhost:5009"))
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

	fmt.Println("title: ", title)
	fmt.Println("content: ", content)

	return c.Redirect(http.StatusMovedPermanently, "/myproject")
}
func detail(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	detail := map[string]interface{}{
		"Id":      id,
		"Title":   "Dumbways Mobile App",
		"Content": "Lorem ipsum dolor sit amet consectetur adipisicing elit. Alias temporibus ea ipsa earum quaerat, beatae aliquid. Adipisci accusamus labore aliquid enim, officiis natus repellat non eos eum perspiciatis neque itaque doloribus deleniti vero velit. Eum obcaecati, consequatur eos fugiat ipsum similique maiores quidem labore vel quasi! Consectetur eius libero repellat.",
	}

	return tmpl.Execute(c.Response(), detail)
}
