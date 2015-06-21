package client

import (
  "fmt"
  dc "github.com/fsouza/go-dockerclient"
  "os"
  "strings"
)

type DockerClient struct {
  dockerHost      string
  dockerCertPath  string
  Client          *dc.Client
}

func NewDockerClient() *DockerClient {

  // Check for required variables
  docker_host := os.Getenv("DOCKER_HOST")
  docker_cert_path := os.Getenv("DOCKER_CERT_PATH")

  if len(strings.TrimSpace(docker_host)) == 0 {
    panic("Please set DOCKER_HOST!")
  }

  if len(strings.TrimSpace(docker_cert_path)) == 0 {
    panic("Please set DOCKER_CERT_PATH")
  }

  // Init the client
  path := os.Getenv("DOCKER_CERT_PATH")
  ca := fmt.Sprintf("%s/ca.pem", path)
  cert := fmt.Sprintf("%s/cert.pem", path)
  key := fmt.Sprintf("%s/key.pem", path)

  client, err := dc.NewTLSClient(os.Getenv("DOCKER_HOST"), cert, key, ca)
  if err != nil {
    panic(err)
  }

  return &DockerClient{
    dockerHost: docker_host,
    dockerCertPath: docker_cert_path,
    Client: client,
  }

}
