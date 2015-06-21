package main

import (
	"fmt"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/chrishenry/geonosis/geonosis/client"
	"github.com/chrishenry/geonosis/geonosis/image"
	"html/template"
	"io"
	"log"
	"net/http"
)

type (
	// Template provides HTML template rendering
	Template struct {
		templates *template.Template
	}

	// user struct {
	// 	ID   string `json:"id"`
	// 	Name string `json:"name"`
	// }
)

// Template provides HTML template rendering
func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Handlers
func createDeployment(c *echo.Context) error {
	return c.String(http.StatusOK, "Deployment POST\n")
}

func getDeployment(c *echo.Context) error {
	docker := client.NewDockerClient()
	containers, err := docker.Client.ListContainers(dc.ListContainersOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, containers)
}

func updateDeployment(c *echo.Context) error {
	return c.String(http.StatusOK, "Deployment PATCH\n")
}

func deleteDeployment(c *echo.Context) error {
	return c.String(http.StatusOK, "Deployment DELETE\n")
}

func getImage(c *echo.Context) error {

  r := c.Request()
  docker := client.NewDockerClient()

  var source string = ""
  var returnimages []image.MyAPIImages

  if len(r.URL.Query()["source"]) == 0 {
    source = "local"
  } else {
    source = r.URL.Query()["source"][0]
  }

  fmt.Println("source: ", source)

  if source == "local" {
    returnimages = image.GetLocalImage(docker)
  }

  return c.JSON(http.StatusOK, returnimages)

}

func main() {

	// Echo instance
	e := echo.New()

	// Debug mode
	e.SetDebug(true)

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Index("public/index.html")

	// Image Routes
	e.Get("/v1/images", getImage)

	// Deployment Routes
	e.Post("/v1/deployments", createDeployment)
	e.Get("/v1/deployments", getDeployment)
	e.Get("/v1/deployments/:id", getDeployment)
	e.Patch("/v1/deployments/:id", updateDeployment)
	e.Delete("/v1/deployments/:id", deleteDeployment)

	// Start server
	e.Run(":1323")
}
