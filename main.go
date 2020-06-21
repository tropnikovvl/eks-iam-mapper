package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kataras/golog"
)

func main() {

	iamK8sGroups := strings.Split(os.Getenv("IAMGROUPS"), ",")
	generateUserRoles(iamK8sGroups)
}

func generateUserRoles(iamK8sGroups []string) {

	// For each iam, extract users and map them to their k8s roles
	for _, iamK8sGroup := range iamK8sGroups {
		iam, role := extractIAMK8sFromString(iamK8sGroup)
		fmt.Printf("iam: %s, k8s: %s\n", iam, role)
	}

}

func extractIAMK8sFromString(str string) (string, string) {
	splits := strings.Split(str, "::")
	if len(splits) != 2 {
		golog.Infof("Flag hatali, [<groupname>]::[<ns>:<role>] seklinde 2 indexli olmalidir!")
		os.Exit(1)
	}
	iam := splits[0]
	k8s := splits[1]
	return iam, k8s
}
