package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type githubRepository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type githubCommit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
}

type githubPush struct {
	Ref        string           `json:"ref"`
	Repository githubRepository `json:"repository"`
	HeadCommit githubCommit     `json:"head_commit"`
}

func main() {
	serverCI := http.NewServeMux()
	serverCI.HandleFunc("/", handleCI)

	http.ListenAndServe(":3000", serverCI)
}

func handleCI(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n-->Request: %v\n\n", r)

	if strings.HasSuffix(r.URL.Path, "/badge") {
		generateBadge(w, r)
	} else if strings.HasSuffix(r.URL.Path, "/output") {
		readOutput(w, r)
	} else {
		if r.URL.Path == "/" {
			runTestsAndBuild(w, r)
		}
	}
}

func runTestsAndBuild(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r: %v\n", r.RequestURI)
	push, err := proccessGithubPayload(r)

	if err == nil {
		if isMiniCi(push) {
			fmt.Printf("File .mini-ci.yml found. Lets build this app!\n")
			//1 - create the commit folder
			createCommitFolder(push)

			//2 - then run the docker image
			execDocker(push)
		}
	} else {
		fmt.Printf("Err: %v\n", err)
	}
}

//proccessGithubPayload unmarshal the github payload and returns the relevant
//data in a proper structure
func proccessGithubPayload(r *http.Request) (githubPush, error) {
	x, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var push githubPush
	err := json.Unmarshal(x, &push)

	return push, err
}

//isMiniCi verify if we need to run the CI in this push
//if the project has a file named .mini-ci.yml then we need run the CI
func isMiniCi(push githubPush) bool {
	urlYml := "https://api.github.com/repos/" + push.Repository.FullName + "/contents/.mini-ci.yml?ref=" + push.HeadCommit.ID
	resp, err := http.Get(urlYml)
	if err == nil {
		if resp.StatusCode == 200 {
			return true
		}
	}

	return false
}

//createCommitFolder creates the folder that represents the commit
//this folder will retain the output returned by the tests
func createCommitFolder(push githubPush) {
	//we gonna use this directory to save the output returned by the tests
	dir := "./repositories/" + push.Repository.FullName + "/" + push.HeadCommit.ID + "/"
	os.MkdirAll(dir, os.ModeDir|0775)
}

func execDocker(push githubPush) {
	envApp := "APP=" + push.Repository.FullName
	envCommit := "COMMIT=" + push.HeadCommit.ID

	fmt.Printf("Run the Docker image\n")
	fmt.Printf("docker run --rm -e '%v' -e '%v' -t minici:dev \n", envApp, envCommit)

	//run the docker image
	out, err := exec.Command("docker", "run", "--rm", "-e", envApp, "-e", envCommit, "-t", "minici:dev").Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Write the output to file\n")

	//write the output to file
	fileOut := "./repositories/" + push.Repository.FullName + "/" + push.HeadCommit.ID + "/output"
	ioutil.WriteFile(fileOut, out, 0644)
	fmt.Printf("%s", string(out))

	writeBuildStatus(push, string(out))
}

func writeBuildStatus(push githubPush, cmdOutput string) {
	//split the output in lines
	lines := strings.Split(cmdOutput, "\n")

	//the last line contain the exit code, we need to get len(-2) because the output comes with an \r
	exitCode := lines[len(lines)-2]
	exitCode = strings.Replace(exitCode, "\r", "", -1)

	//create the folder that will contain the build result file
	parts := strings.Split(push.Ref, "/")
	dir1 := parts[0]
	dir2 := parts[1]
	fileName := parts[2]
	os.MkdirAll("./repositories/"+push.Repository.FullName+"/"+dir1+"/"+dir2+"/", os.ModeDir|0775)

	//file that will store the build result
	buildResultFile := "./repositories/" + push.Repository.FullName + "/" + dir1 + "/" + dir2 + "/" + fileName

	//exitCode = 0 [sucess]
	//exitCode = 1 [failed]
	if exitCode == "0" {
		ioutil.WriteFile(buildResultFile, []byte("success"), 0644)
	} else {
		ioutil.WriteFile(buildResultFile, []byte("failed"), 0644)
	}
}

func readOutput(w http.ResponseWriter, r *http.Request) {
	result, err := ioutil.ReadFile("." + r.URL.Path)
	if err == nil {
		w.Write([]byte(result))
	} else {
		w.Write([]byte(err.Error()))
	}
}

func generateBadge(w http.ResponseWriter, r *http.Request) {
	result, err := ioutil.ReadFile("." + r.URL.Path[:len(r.URL.Path)-6])
	strResult := strings.Replace(string(result), "\n", "", -1)

	if err == nil {

		if strResult == "success" {
			writeBadge(w, "badges/pass.png")
		} else if strResult == "failed" {
			writeBadge(w, "badges/fail.png")
		} else {
			w.Write([]byte("Error obtaining build status"))
		}

	} else {
		w.Write([]byte(err.Error()))
	}
}

func writeBadge(w http.ResponseWriter, filename string) {
	const layout = "Mon, 02 Jan 2006 15:04:05 GMT"

	w.Header().Add("content-type", "image/png")
	w.Header().Set("cache-control", "no-cache, private")
	w.Header().Set("date", time.Now().UTC().Format(layout))
	w.Header().Set("expires", "Fri, 01 Jan 1984 00:00:00 GMT")

	icon, _ := ioutil.ReadFile(filename)
	w.Write(icon)
}
