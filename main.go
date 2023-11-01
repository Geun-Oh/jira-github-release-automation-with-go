package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/actions-go/toolkit/core"
)

var domain, project, version, token string
var isCreateNextVersion bool

type ValueType struct {
	Self            string `json:"self"`
	Id              string `json:"id"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	Archived        bool   `json:"archived"`
	Released        bool   `json:"released"`
	StartDate       string `json:"startDate"`
	ReleaseDate     string `json:"releaseDate"`
	Overdue         bool   `json:"overdue"`
	UserStartDate   string `json:"userStartDate"`
	UserReleaseDate string `json:"userReleaseDate"`
	ProjectId       int    `json:"projectId"`
}

type VersionType struct {
	Self       string      `json:"self"`
	MaxResults int         `json:"maxResults"`
	StartAt    int         `json:"startAt"`
	Total      int         `json:"total"`
	IsLast     bool        `json:"isLast"`
	Values     []ValueType `json:"values"`
}

func getVersionName(version string) string {
	r, _ := regexp.Compile(".*(\\d{1,5}\\.\\d{1,5}\\.\\d{1,5})")

	var match = r.FindStringSubmatch(version)

	if match == nil || len(match) < 1 {
		return ""
	}
	return match[0]
}

// func getVersions(appName, domain, project, token string) bool {

// }

type PostBody struct {
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate"`
}

func releaseVersion(versionId string) {
	fmt.Println(versionId)
	fmt.Println("Trying release version")
	reqBody := PostBody{
		Released:    true,
		ReleaseDate: time.Now().String()[:10],
	}

	rBytes, _ := json.Marshal(reqBody)
	buff := bytes.NewBuffer(rBytes)

	releaseUri := `https://drivingteacher.atlassian.net/rest/api/3/version/` + versionId
	fmt.Println(releaseUri, reqBody, rBytes, buff)
	req, err := http.NewRequest("PUT", releaseUri, buff)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	releaseClient := &http.Client{}
	resp, clientErr := releaseClient.Do(req)

	if clientErr != nil {
		fmt.Println("Error")
		panic(clientErr)
	}

	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	var dat VersionType

	if err := json.Unmarshal([]byte(str), &dat); err != nil {
		panic(err)
	}

	fmt.Println(dat)

	fmt.Println("Version Release Succeded")
}

func main() {
	domain, _ := core.GetInput("domain")
	project, _ := core.GetInput("project")
	releaseName, _ := core.GetInput("releaseName")
	// version, _ := core.GetInput("version")
	token, _ := core.GetInput("token")
	// isCreateNextVersion := core.GetBoolInput("create-next-version")

	// versionName := getVersionName(version)

	// if len(version) < 1 {
	// 	core.SetFailed("Version is Not correct.")
	// 	return
	// }

	// appName := strings.Replace(version, versionName, "", -1)

	// println(`appName: `, appName, `versionName: `, versionName)

	// getVerson

	uri := `https://` + domain + `atlassian.net/rest/api/3/project/` + project + `/version?query=` + url.QueryEscape(releaseName) + `&orderBy=name&status=unreleased`
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", `Basic `+token)

	client := &http.Client{}
	resp, clientErr := client.Do(req)

	if clientErr != nil {
		panic(clientErr)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	var dat VersionType

	if err := json.Unmarshal([]byte(str), &dat); err != nil {
		panic(err)
	}

	var currentVersion string

	for _, v := range dat.Values {
		if v.Name == "DT test" {
			currentVersion = v.Name
		}
	}

	fmt.Println(dat.IsLast, dat.Values[0].Id, currentVersion)

	// releaseVersion(dat.Values[0].Id)
}
