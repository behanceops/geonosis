package image

import (
  "fmt"
  dc "github.com/fsouza/go-dockerclient"
  "github.com/chrishenry/geonosis/geonosis/client"
  "log"
)

var Test string = "test"

type MyAPIImages struct {
  ID          string    `json:"Id" yaml:"Id"`
  RepoTag     string    `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
  Source      string    `json:"Source,omitempty" yaml:"Source, omitempty"`
  APIImages   dc.APIImages
}

func GetLocalImage(c *client.DockerClient) []MyAPIImages {

  var returnimages []MyAPIImages

  // Get local images
  images, err := c.Client.ListImages(dc.ListImagesOptions{All: true})
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
