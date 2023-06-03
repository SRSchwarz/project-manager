package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Project struct {
	Location string `json:"location"`
	Language string `json:"language"`
}

var projects []Project
var configFolderName = "project-manager"
var projectsFileName = "projects.json"
var homeDir string

func main() {
	initConfigFolder()
	loadProjects()
	method, _ := parseArguments()
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to determine current working directory")
	}

	if method == "add" {
		addProject(workingDirectory)
		return
	}
	if method == "open" {
		selectAndOpenProject()
		return
	}

	log.Fatalf("Command not recognized")
}

func initConfigFolder() {
	usr, _ := user.Current()
	homeDir = usr.HomeDir
	os.MkdirAll(filepath.Join(homeDir, configFolderName), os.ModePerm)
	projectsFile, _ := os.OpenFile(filepath.Join(homeDir, configFolderName, projectsFileName), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if projectsFile != nil {
		projectsFile.Write([]byte("[]"))
	}

}

func parseArguments() (string, string) {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("No arguments given")
	}
	method := args[0]
	if len(args) == 1 {
		return method, ""
	}
	path := args[1]
	if path != "." {
		log.Fatal("Adding specific paths is currently not supported")
	}
	return method, path
}

func loadProjects() {
	projectsFileContent, err := os.ReadFile(filepath.Join(homeDir, configFolderName, projectsFileName))
	if err != nil {
		log.Fatal("Failed to open projects file")
	}

	err = json.Unmarshal(projectsFileContent, &projects)
	if err != nil {
		log.Fatal("Failed to read projects file")
	}
}

func addProject(location string) {
	for _, p := range projects {
		if p.Location == location {
			log.Printf("Project '%s' already exists", location)
			return
		}
	}

	newProject := Project{Location: location, Language: "go"}
	projects = append(projects, newProject)
	newProjectFileContent, err := json.Marshal(projects)
	if err != nil {
		log.Fatal("Failed to save new project")
	}

	file, err := os.Create(filepath.Join(homeDir, configFolderName, projectsFileName))
	if err != nil {
		log.Fatal("Failed to save new project")
	}

	_, err = file.Write(newProjectFileContent)
	if err != nil {
		log.Fatal("Failed to save new project")
	}

	log.Printf("Successfully saved project '%s' for language '%s'", location, "go")
}

func selectAndOpenProject() {
	projectLocations := ""
	for _, p := range projects {
		projectLocations += p.Location + "\n"
	}
	cmd := exec.Command("bash", "-c", "printf '"+projectLocations+"' | fzf")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to select project")
	}

	selectedProject := string(output)
	cmd = exec.Command("bash", "-c", "code "+selectedProject)
	cmd.Run()
}
