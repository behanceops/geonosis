package main

import (
	"encoding/json"
	"fmt"
	"github.com/chrishenry/geonosis/geonosis/client"
	"github.com/chrishenry/geonosis/geonosis/image"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"io/ioutil"
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

type test_struct struct {
	Image string
	Name  string
	Cmd   []string
	Env   []string
}

// Handlers
func createDeployment(c *echo.Context) error {

	r := c.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var t test_struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Image)

	docker := client.NewDockerClient()

	cmd := t.Cmd
	if len(t.Cmd) == 0 {
		cmd = nil
	}

	env := t.Env
	if len(t.Env) == 0 {
		env = nil
	}

	createContConf := dc.Config{
		Image: t.Image,
		Cmd:   cmd,
		Env:   env,
	}

	createContHostConfig := dc.HostConfig{}

	createContOps := dc.CreateContainerOptions{
		Name:       t.Name,
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	container, err := docker.Client.CreateContainer(createContOps)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\ncontainer = %s\n", container)

	err = docker.Client.StartContainer(container.ID, nil)
	if err != nil {
		panic(err)
	}

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

	var source string = ""
	var returnimages []image.MyAPIImages

	if len(r.URL.Query()["source"]) == 0 {
		source = "local"
	} else {
		source = r.URL.Query()["source"][0]
	}

	fmt.Println("source: ", source)

	if source == "local" {
		returnimages = image.GetLocalImage(client.NewDockerClient())
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
