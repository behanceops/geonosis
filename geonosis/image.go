package image

type Test string = "test"

type MyAPIImages struct {
  ID          string    `json:"Id" yaml:"Id"`
  RepoTag     string    `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
  Source      string    `json:"Source,omitempty" yaml:"Source, omitempty"`
  APIImages   dc.APIImages
}

func GetImage(c *echo.Context) error {

  r := c.Request()
  var source string = ""
  var returnimages []MyAPIImages

  if len(r.URL.Query()["source"]) == 0 {
    source = "local"
  } else {
    source = r.URL.Query()["source"][0]
  }

  fmt.Println("source: ", source)

  if source == "local" {

    returnimages = getLocalImage()

  }

  return c.JSON(http.StatusOK, returnimages)

}

func getLocalImage() []MyAPIImages {

  docker := getDockerClient()

  var returnimages []MyAPIImages

  // Get local images
  images, err := docker.ListImages(dc.ListImagesOptions{All: true})
  if err != nil {
    log.Fatal(err)
  }

  for _, img := range images {

    if img.RepoTags[0] != "<none>:<none>" {

      for _, tag := range img.RepoTags {
        fmt.Println("RepoTags: ", tag)
        returnimages = append(returnimages, MyAPIImages{img.ID, tag, "local", img})
      }

    }

  }

  return returnimages

}
