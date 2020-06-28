package main

import (
	"os"
	"strings"
	"time"

	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	iamK8sGroups := strings.Split(os.Getenv("GROUPSANDROLES"), ",")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		userRoles := generateUserRoles(iamK8sGroups)
		cf, err := clientset.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		var newConfig []MapUserConfig

		for _, userRole := range userRoles {
			newConfig = append(newConfig, MapUserConfig{
				UserArn:  userRole.IAMArn,
				Username: userRole.IAMUsername,
				Groups:   userRole.K8sRoles,
			})
		}

		// If there are no users to add, the config map will be empty
		// Since this will never be the intended purpose of the user
		// and this case is more likely to happen due to a bug
		// we'll just skip the changes
		if len(newConfig) == 0 {
			golog.Info("No users found, config will not be changed")
			continue
		}

		roleStr, err := yaml.Marshal(newConfig)
		if err != nil {
			golog.Error(err)
		}
		cf.Data["mapUsers"] = string(roleStr)

		_, err = clientset.CoreV1().ConfigMaps("kube-system").Update(cf)
		if err != nil {
			golog.Error(err)
		} else {
			golog.Info("Successfully updated user roles")
		}
		time.Sleep(1 * time.Minute)
	}
}
