package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/rivo/tview"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Import metav1 package
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	//Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	//Or uncomment to load specific auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// type Namespace struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty"`
// 	Items           []Namespace `json:"items"`
// }

func homeScreen(clientset *kubernetes.Clientset) {
	Namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// defines new app
	app := tview.NewApplication()

	// Create the list
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Namespaces") // have to define the box parameters outside the function return values otherwise dont work

	podsList := tview.NewList()
	podsList.SetBorder(true).SetTitle("Pods")

	for _, ns := range Namespaces.Items {
		list.AddItem(ns.Name, "", 0, func() {
			podsList.Clear()

			pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				list.AddItem("", "No Items to Show", 0, nil)
			} else {
				for _, pod := range pods.Items {
					podsList.AddItem(pod.Name, "", 0, nil)
				}
			}
		})
	}

	// flex columns puts the boxes in columns
	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(list, 0, 1, true).     // numbers are wight and size
		AddItem(podsList, 0, 1, false) // boolean decided what box is selected first. So list box is selected first for

	// Set the list as the root (centers the list for me) and run the application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
// type tests - https://github.com/rivo/tview/blob/b0a7293b81308ab0e44c2757f8922683a3d659df/demos/button/main.go
// https://github.com/rivo/tview/wiki/List
func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	homeScreen(clientset)

}
