package main

import (
	"github.com/joho/godotenv"
	_ "github.com/nguyenvanduocit/autorender"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/nguyenvanduocit/autorender"
)

type Size struct{
	Width int	`json:"width"`
	Height int `json:"height"`
}

type Image struct {
	FileName string `json:"file_name"`
	Size *Size `json:"size"`
}

type Template struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	Images []*Image `json:"images"`
}

type Project struct {
	Id    int `json:"id"`
	Name       string `json:"name"`
	ResultSrc    string `json:"result_src"`
	Images map[string]string `json:"images"`
	Template *Template `json:"template"`
	Status string `json:"status"`
}

func main(){

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	endpoint := os.Getenv("API_ENDPOINT")
	aeRenderPath := os.Getenv("AERENDER_PATH")
	// Get project from server
	httpClient := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/projects/pop", endpoint), nil)
	request.Header.Set("X-CSRF-TOKEN", "qzpaRGYMJM1xbB2PaI4BKLnwFkBLyFkRIU78gJvS")
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer response.Body.Close()
	//Parse project
	var project Project
	err = json.NewDecoder(response.Body).Decode(&project)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Load Local project spec.json
	projectPath := fmt.Sprintf("./templates/%s", project.Template.Slug)
	projectSpecFile := fmt.Sprintf("%s/spec.json", projectPath)
	if _, err := os.Stat(projectSpecFile); os.IsNotExist(err) {
		log.Fatal(err.Error())
	}
	specFile, err := ioutil.ReadFile(projectSpecFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	var templateObject Template
	json.Unmarshal(specFile, &templateObject)
	//Create AutoRender
	render,err := autorender.New(aeRenderPath)
	if err != nil {
		log.Fatal(err)
	}
	assets := []autorender.Asset{}
	for _,image := range templateObject.Images{
		asset := autorender.Asset{
			Type: "image",
			Src: project.Images[image.FileName],
			Name: image.FileName,
		}
		assets = append(assets, asset)
	}
	aeTemplateFile := fmt.Sprintf("%s/template.aepx", projectPath)
	aeProject := autorender.NewProject(aeTemplateFile, "main", assets, "ahihi", 1)
	log.Println("Project created: ", aeProject.ID)
	render.Render(aeProject)
}
