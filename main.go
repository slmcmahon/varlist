package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/slmcmahon/go-azdo"
	slmcommon "github.com/slmcmahon/go-common"
)

// ByName implements sort.Interface based on the Name field.
type ByName []azdo.VarLib

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

	response, err := azdo.GetVariableLibraries(pat, org, project)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(ByName(response.Value))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID"})

	for _, varlib := range response.Value {
		table.Append([]string{varlib.Name, fmt.Sprintf("%d", varlib.ID)})
	}

	table.Render()
}
