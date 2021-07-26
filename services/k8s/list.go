package k8s

import (
	"context"
	"fmt"
	"github.com/raffs/sysadmin-sk/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type ListOptions struct {
	kind 			string `type:"string" required:"true"`
	namespace		string `type:"string" required:"true"`
}

func listResource(options *ListOptions) error {
	client, err := utils.K8sClient()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	switch strings.Title(options.kind) {
	case "Deployment":
		deps, err := client.AppsV1().Deployments(options.namespace).List(context.TODO(),v1.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		for _, res := range (*deps).Items {
			fmt.Println("***********************************************************")
			fmt.Printf("deployments-name=%v\n", res.Name)
			fmt.Printf("deployments-status=%v\n", res.Status)
			fmt.Printf("deployments-creationtime=%v\n", res.CreationTimestamp)
			fmt.Printf("deployments-available-replicas=%v\n", res.Status.AvailableReplicas)
			fmt.Printf("deployments-unavailable-replicas=%v\n", res.Status.UnavailableReplicas)
			fmt.Printf("deployments-namespace=%v\n", res.Namespace)
		}
	case "Daemonset":
		dems, err := client.AppsV1().DaemonSets(options.namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		for _, res := range (*dems).Items {
			fmt.Println("***********************************************************")
			fmt.Printf("daemonsets-name=%v\n", res.Name)
			fmt.Printf("daemonset-status=%v\n", res.Status)
			fmt.Printf("daemonset-creationtime=%v\n", res.CreationTimestamp)
			fmt.Printf("daemonset-replicas-available=%v\n", res.Status.NumberAvailable)
			fmt.Printf("daemonset-replicas-unavailable=%v\n", res.Status.NumberUnavailable)
			fmt.Printf("deployments-namespace=%v\n", res.Namespace)
		}
	case "StatefulSet":
		ss, err := client.AppsV1().StatefulSets(options.namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		for _, res := range (*ss).Items {
			fmt.Println("***********************************************************")
			fmt.Printf("statefulsets-name=%v\n", res.Name)
			fmt.Printf("statefulsets-status=%v\n", res.Status)
			fmt.Printf("statefulset-creationtime=%v\n", res.CreationTimestamp)
			fmt.Printf("statement-desired-replicas=%v\n", res.Status.Replicas)
			fmt.Printf("daemonset-ready-replicas=%v\n", res.Status.ReadyReplicas)
			fmt.Printf("deployments-namespace=%v\n", res.Namespace)
		}
	case "Service":
		srv, err := client.CoreV1().Services(options.namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		for _, res := range (*srv).Items {
			fmt.Println("***********************************************************")
			fmt.Printf("services-name=%v\n", res.Name)
			fmt.Printf("services-status=%v\n", res.Status)
			fmt.Printf("service-creationtime=%v\n", res.CreationTimestamp)
			fmt.Printf("daemonset-labels=%v\n", res.Labels)
			fmt.Printf("deployments-namespace=%v\n", res.Namespace)
		}
	case "Pod":
		p, err := client.CoreV1().Pods(options.namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		for _, res := range (*p).Items {
			fmt.Println("***********************************************************")
			fmt.Printf("pod-name=%v\n", res.Name)
			fmt.Printf("pod-status=%v\n", res.Status)
			fmt.Printf("pod-creationtime=%v\n", res.CreationTimestamp)
			fmt.Printf("pod-labels=%v\n", res.Labels)
			fmt.Printf("pod-namespace=%v\n", res.Namespace)
		}
	}
	return nil
}

func ListResources() *cobra.Command {
	var options ListOptions

	cmd := &cobra.Command{
		Use:   "list-resource",
		Short: "List k8s namespaced resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := utils.ValidateArgs(args)
			if err != nil {
				return err
			}
			return listResource(&options)
		},
	}
	cmd.PersistentFlags().StringVarP(&options.kind, "kind", "t", "", "List ResourceType")
	cmd.PersistentFlags().StringVarP(&options.namespace, "namespace", "t", "", "List Resource Namespace")
	return cmd
}