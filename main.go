package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os/exec"
)

type GithubRepository struct {
        Id       int    `json:"id"`
        Name     string `json:"name"`
        FullName string `json:"full_name"`
}

type GithubCommit struct {
        Id        string "id"
        Message   string `json:"message"`
        Timestamp string `json:"timestamp"`
        Url       string `json:"url"`
}

type GithubPush struct {
        Ref        string           `json:"ref"`
        Repository GithubRepository `json:"repository"`
        HeadCommit GithubCommit     `json:"head_commit"`
}

func main() {
		docker()
        //serverCI := http.NewServeMux()
        //serverCI.HandleFunc("/", postReceive)

        //http.ListenAndServe(":3000", serverCI)
}

func postReceive(w http.ResponseWriter, r *http.Request) {
        x, _ := ioutil.ReadAll(r.Body)
        defer r.Body.Close()

        var push GithubPush
        err := json.Unmarshal(x, &push)

        if err == nil {
                fmt.Printf("%v\n", push)
                urlYml := "https://api.github.com/repos/" + push.Repository.FullName + "/contents/.mini-ci.yml?ref=" + push.Ref
                fmt.Printf("GET %v\n", urlYml)
                resp, err := http.Get(urlYml)
                if err == nil {
                        if resp.StatusCode == 200 {
                                fmt.Printf("File .mini-ci.yml found.Lets build this app!\n")
                        } else {
                                fmt.Printf("File .mini-ci.yml not found. Exiting.[%v]\n", resp.Status)
                        }
                } else {
                        fmt.Printf("Err: %v\n", err)
                }
        } else {
                fmt.Printf("Err: %v\n", err)
        }
}

func docker(){
	out, err := exec.Command("docker", "run", "--rm", "-t", "-v", "/home/maicon/go/src/github.com/maiconio/mini-ci/docker-stuff:/home/docker", "--name", "minici",  "maiconio/minici:dev").Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", out)
}
