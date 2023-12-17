package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	slmcommon "github.com/slmcmahon/go-common"
)

type VarLib struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type VarLibResponse struct {
	Count   int      `json:"count"`
	VarLibs []VarLib `json:"value"`
}

// ByName implements sort.Interface based on the Name field.
type ByName []VarLib

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func main() {
	var (
		patFlag     string
		orgFlag     string
		projectFlag string
	)

	flag.StringVar(&orgFlag, "org", "", "Azure Devops Organization")
	flag.StringVar(&projectFlag, "project", "", "Azure DevOps Project")
	flag.StringVar(&patFlag, "pat", "", "Personal Access Token")
	flag.Parse()

	pat, err := slmcommon.CheckEnvOrFlag(patFlag, "AZDO_PAT")
	if err != nil {
		log.Fatal(err)
	}

	org, err := slmcommon.CheckEnvOrFlag(orgFlag, "AZDO_ORG")
	if err != nil {
		log.Fatal(err)
	}

	project, err := slmcommon.CheckEnvOrFlag(projectFlag, "AZDO_PROJECT")
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/distributedtask/variablegroups?api-version=6.0-preview.2", org, project)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("", pat)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response VarLibResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		panic(err)
	}

	sort.Sort(ByName(response.VarLibs))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID"})

	for _, varlib := range response.VarLibs {
		table.Append([]string{varlib.Name, fmt.Sprintf("%d", varlib.Id)})
	}

	table.Render()
}
