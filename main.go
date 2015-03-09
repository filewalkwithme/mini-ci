package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/* curl 127.0.0.1:3000 -d '{"ref": "refs/heads/teste", "head_commit": {"id": "abf230d81caf727884a56fae9acf35788d3ae9e7", "message": "teste", "timestamp": "2015-02-25T11:56:09-03:00", "url": "https://github.com/maiconio/portugo/commit/abf230d81caf727884a56fae9acf35788d3ae9e7"}, "repository": {"id": 18707655, "name": "portugo", "full_name": "maiconio/portugo"}}'
{
	"ref": 					"refs/heads/teste",
	"head_commit": {
		"id": 				"abf230d81caf727884a56fae9acf35788d3ae9e7",
		"message": 		"teste",
		"timestamp": 	"2015-02-25T11:56:09-03:00",
		"url": 				"https://github.com/maiconio/portugo/commit/abf230d81caf727884a56fae9acf35788d3ae9e7"
	},
	"repository": {
		"id": 				18707655,
		"name": 			"portugo",
		"full_name": 	"maiconio/portugo"
	}
}
*/

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
	x, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(x)
	//docker()
	serverCI := http.NewServeMux()
	serverCI.HandleFunc("/", postReceive)

	http.ListenAndServe(":3000", serverCI)
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

//createCommitFolder creates the folder that represents the commit
//this folder will retain the output returned by the tests
func createCommitFolder(push githubPush) {
	//we gonna use this directory to save the output returned by the tests
	dir := "./repositories/" + push.Repository.FullName + "/" + push.HeadCommit.ID + "/"
	os.MkdirAll(dir, os.ModeDir|0775)
}

//isMiniCi verify if we need to run the CI in this push
//if the project has a file named .mini-ci.yml then we need run the CI
func isMiniCi(push githubPush) bool {
	urlYml := "https://api.github.com/repos/" + push.Repository.FullName + "/contents/.mini-ci.yml?ref=" + push.Ref
	resp, err := http.Get(urlYml)
	if err == nil {
		if resp.StatusCode == 200 {
			return true
		}
	}

	return false
}

//postReceive handles the github webhook.
func postReceive(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path[1:], "/")
	if urlParts[0] == "repositories" {
		fmt.Printf("." + r.URL.Path)
		tmp, err := ioutil.ReadFile("." + r.URL.Path)
		if err == nil {
			w.Write(tmp)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		push, err := proccessGithubPayload(r)

		if err == nil {
			if isMiniCi(push) {
				fmt.Printf("File .mini-ci.yml found. Lets build this app!\n")
				//1 - create the commit folder
				createCommitFolder(push)

				//2 - then run the docker image, mounting this directory as home
				execDocker(push)
			}
		} else {
			fmt.Printf("Err: %v\n", err)
		}
	}
}

func execDocker(push githubPush) {
	envApp := "APP=" + push.Repository.FullName
	envCommit := "COMMIT=" + push.HeadCommit.ID

	//run the docker image
	out, err := exec.Command("docker", "run", "--rm", "-e", envApp, "-e", envCommit, "-t", "maiconio/minici:dev").Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	//write the output to file
	fileOut := "./repositories/" + push.Repository.FullName + "/" + push.HeadCommit.ID + "/output"
	ioutil.WriteFile(fileOut, out, 0644)
	fmt.Printf("%s", string(out))

	//split the output in lines
	lines := strings.Split(string(out), "\n")

	//the last line contain the exit code, we need to get len(-2) because the output comes with an \r
	exitCode := lines[len(lines)-2]
	exitCode = strings.Replace(exitCode, "\r", "", -1)

	//exitCode = 0 [sucess]
	//exitCode = 1 [failed]

	parts := strings.Split(push.Ref, "/")
	dir1 := parts[0]
	dir2 := parts[1]
	fileName := parts[2]
	fileTree := "./repositories/" + push.Repository.FullName + "/" + dir1 + "/" + dir2 + "/" + fileName
	os.MkdirAll("./repositories/"+push.Repository.FullName+"/"+dir1+"/"+dir2+"/", os.ModeDir|0775)
	if exitCode == "0" {
		err := ioutil.WriteFile(fileTree, []byte("success"), 0644)
		fmt.Printf("%s\n", err)

		fmt.Printf("Build Success! =D [%v]\n", exitCode)
	} else {
		ioutil.WriteFile(fileTree, []byte("failed"), 0644)

		fmt.Printf("Build Failed![%v]\n", exitCode)
	}

}
