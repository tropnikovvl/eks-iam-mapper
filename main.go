package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
)

func main() {

	iamK8sGroups := strings.Split(os.Getenv("IAMGROUPS"), ",")
	userRoles := generateUserRoles(iamK8sGroups)

	var newConfig []MapUserConfig
	for _, userRole := range userRoles {
		newConfig = append(newConfig, MapUserConfig{
			UserArn:  userRole.IAMArn,
			Username: userRole.IAMUsername,
			Groups:   userRole.K8sRoles,
		})
	}

	roleStr, _ := yaml.Marshal(newConfig)
	fmt.Println(string(roleStr))

}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func generateUserRoles(iamK8sGroups []string) map[string]UserRoles {
	userRoles := make(map[string]UserRoles)
	// For each iam, extract users and map them to their k8s roles
	for _, iamK8sGroup := range iamK8sGroups {
		iam, role := extractIAMK8sFromString(iamK8sGroup)
		//fmt.Printf("--\niam: %s, k8s: %s\n", iam, role)
		users := getAwsGroups(iam)
		//fmt.Printf("users:\n %s", users.GoString())
		for _, user := range users.Users {
			if _, exists := userRoles[*user.UserName]; !exists {
				userRoles[*user.UserName] = UserRoles{IAMArn: *user.Arn, IAMUsername: *user.UserName, K8sRoles: []string{}}
			}
			userRoles[*user.UserName] = userRoles[*user.UserName].SetK8sRoles(strings.Split(role, "|"))
		}
	}
	for iamUsername := range userRoles {
		userRoles[iamUsername] = userRoles[iamUsername].UniqueK8sRoles()
	}

	//fmt.Println(userRoles["ahmet.soykan"])
	return userRoles
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
