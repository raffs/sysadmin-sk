package k8s

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/dynamic"
)

type K8sOptions struct {
	manifestPath string `type:"string" required:"true"`
}

func applyManifest(options *K8sOptions) error {
	// Need to move this into client.go
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	filestream, err := ioutil.ReadFile(options.manifestPath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(filestream))
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filestream), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			fmt.Println("Decode ERROR: " + err.Error())
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(c.Discovery())
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	if err != io.EOF {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func validateArgs(options *K8sOptions, args []string) error {
	if len(args) != 1 {
		return errors.New("Invalid number of arguments for k8s apply manifest command. Use --help for details")
	}

	return nil
}

func ApplyManifest() *cobra.Command {
	var options K8sOptions

	cmd := &cobra.Command{
		Use:   "apply-manifest",
		Short: "Apply an instantiated kubernetes manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := validateArgs(&options, args)
			if err != nil {
				return err
			}
			return applyManifest(&options)
		},
	}
	cmd.PersistentFlags().StringVarP(&options.manifestPath, "manifest-path", "p", "", "Path to the k8s manifest to apply")
	return cmd
}