package hellokubernetes

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	helloLong = templates.LongDesc(i18n.T(`
		Prints config of resource by filename or stdin.
		
		Prints config of resource anf date of creation when type/name is passed as arg

		JSON and YAML formats are accepted.`))

	helloExample = templates.Examples(i18n.T(`
		# Print the configuration in pod.json to a pod.
		kubectl hello-kubernetes -f ./pod.json

		# Print resources from a directory containing kustomization.yaml - e.g. dir/kustomization.yaml.
		kubectl hello-kubernetes -k dir/

		# Print JSON passed into stdin to a pod.
		cat pod.json | kubectl hello-kubernetes -f -
		
		# Print type/name
		kubectl hello-kubernetes type/name`))
)

// HelloKubernetesOptions have the data required to perform the hello-kubernetes operation
type HelloKubernetesOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	PrintObj   func(obj kruntime.Object) error

	// Filename options
	FilenameOptions resource.FilenameOptions
	RecordFlags     *genericclioptions.RecordFlags
	args            []string

	// Common user flags

	// results of arg parsing
	Recorder                     genericclioptions.Recorder
	namespace                    string
	enforceNamespace             bool
	unstructuredClientForMapping func(mapping *meta.RESTMapping) (resource.RESTClient, error)

	genericclioptions.IOStreams
}

// NewHelloKubernetesOptions creates the options for hello-world
func NewHelloKubernetesOptions(ioStreams genericclioptions.IOStreams) *HelloKubernetesOptions {
	return &HelloKubernetesOptions{
		PrintFlags: genericclioptions.NewPrintFlags("hello-kubernetes").WithTypeSetter(scheme.Scheme),

		RecordFlags: genericclioptions.NewRecordFlags(),
		Recorder:    genericclioptions.NoopRecorder{},
		IOStreams:   ioStreams,
	}
}

// NewCmdHelloKubernetes creates the `hello-world` command
func NewCmdHelloKubernetes(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewHelloKubernetesOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "hello-kubernetes",
		DisableFlagsInUseLine: true,
		Short:                 "prints Hello <resource-name> <kind> or Hello <kind of resource> <name of resource> <creation time>",
		Long:                  helloLong,
		Example:               helloExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.RunHelloKubernetes(f, cmd))
		},
	}

	// bind flag structs
	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	usage := "for printing \"Hello <resource-name> <kind>\""
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)

	return cmd
}

// Complete adapts from the command line args and factory to the data required.
func (o *HelloKubernetesOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	o.args = args
	o.RecordFlags.Complete(cmd)
	o.Recorder, err = o.RecordFlags.ToRecorder()
	if err != nil {
		return err
	}

	o.namespace, o.enforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.PrintObj = func(obj kruntime.Object) error {
		// return printer.PrintObj(obj, o.Out)
		// fmt.Printf("The out we expecting: %+v\n", obj) // below is for flags -f
		// The out we expecting: &{Object:map[apiVersion:v1 kind:ReplicationController metadata:map[labels:map[name:test-rc] name:test-rc namespace:test] spec:map[replicas:1 template:map[metadata:map[labels:map[name:test-rc]] spec:map[containers:[map[image:nginx name:test-rc ports:[map[containerPort:80]]]]]]]]}
		// below out is for type/name
		// The out we expecting: &{Object:map[apiVersion:v1 kind:Pod metadata:map[creationTimestamp:2020-06-30T18:48:12Z generateName:hello-node-576f8c6496- labels:map[app:hello-node pod-template-hash:576f8c6496] managedFields:[map[apiVersion:v1 fieldsType:FieldsV1 fieldsV1:map[f:metadata:map[f:generateName:map[] f:labels:map[.:map[] f:app:map[] f:pod-template-hash:map[]] f:ownerReferences:map[.:map[] k:{"uid":"5a20d3d7-c6af-4fb8-80a7-650a4a353b80"}:map[.:map[] f:apiVersion:map[] f:blockOwnerDeletion:map[] f:controller:map[] f:kind:map[] f:name:map[] f:uid:map[]]]] f:spec:map[f:containers:map[k:{"name":"echoserver"}:map[.:map[] f:image:map[] f:imagePullPolicy:map[] f:name:map[] f:resources:map[] f:terminationMessagePath:map[] f:terminationMessagePolicy:map[]]] f:dnsPolicy:map[] f:enableServiceLinks:map[] f:restartPolicy:map[] f:schedulerName:map[] f:securityContext:map[] f:terminationGracePeriodSeconds:map[]]] manager:kube-controller-manager operation:Update time:2020-06-30T18:48:12Z] map[apiVersion:v1 fieldsType:FieldsV1 fieldsV1:map[f:status:map[f:conditions:map[k:{"type":"ContainersReady"}:map[.:map[] f:lastProbeTime:map[] f:lastTransitionTime:map[] f:status:map[] f:type:map[]] k:{"type":"Initialized"}:map[.:map[] f:lastProbeTime:map[] f:lastTransitionTime:map[] f:status:map[] f:type:map[]] k:{"type":"Ready"}:map[.:map[] f:lastProbeTime:map[] f:lastTransitionTime:map[] f:status:map[] f:type:map[]]] f:containerStatuses:map[] f:hostIP:map[] f:phase:map[] f:podIP:map[] f:podIPs:map[.:map[] k:{"ip":"172.17.0.3"}:map[.:map[] f:ip:map[]]] f:startTime:map[]]] manager:kubelet operation:Update time:2020-06-30T18:48:43Z]] name:hello-node-576f8c6496-7btxc namespace:default ownerReferences:[map[apiVersion:apps/v1 blockOwnerDeletion:true controller:true kind:ReplicaSet name:hello-node-576f8c6496 uid:5a20d3d7-c6af-4fb8-80a7-650a4a353b80]] resourceVersion:449 selfLink:/api/v1/namespaces/default/pods/hello-node-576f8c6496-7btxc uid:d0c6003e-420b-4f7d-ad34-cea40c2101a4] spec:map[containers:[map[image:k8s.gcr.io/echoserver:1.4 imagePullPolicy:IfNotPresent name:echoserver resources:map[] terminationMessagePath:/dev/termination-log terminationMessagePolicy:File volumeMounts:[map[mountPath:/var/run/secrets/kubernetes.io/serviceaccount name:default-token-2f7wp readOnly:true]]]] dnsPolicy:ClusterFirst enableServiceLinks:true nodeName:127.0.0.1 priority:0 restartPolicy:Always schedulerName:default-scheduler securityContext:map[] serviceAccount:default serviceAccountName:default terminationGracePeriodSeconds:30 tolerations:[map[effect:NoExecute key:node.kubernetes.io/not-ready operator:Exists tolerationSeconds:300] map[effect:NoExecute key:node.kubernetes.io/unreachable operator:Exists tolerationSeconds:300]] volumes:[map[name:default-token-2f7wp secret:map[defaultMode:420 secretName:default-token-2f7wp]]]] status:map[conditions:[map[lastProbeTime:<nil> lastTransitionTime:2020-06-30T18:48:12Z status:True type:Initialized] map[lastProbeTime:<nil> lastTransitionTime:2020-06-30T18:48:43Z status:True type:Ready] map[lastProbeTime:<nil> lastTransitionTime:2020-06-30T18:48:43Z status:True type:ContainersReady] map[lastProbeTime:<nil> lastTransitionTime:2020-06-30T18:48:12Z status:True type:PodScheduled]] containerStatuses:[map[containerID:docker://2a3238e6e509447215da42ed6432a7c1ac85e41ab1b2169b142676614f8dfa4c image:k8s.gcr.io/echoserver:1.4 imageID:docker-pullable://k8s.gcr.io/echoserver@sha256:5d99aa1120524c801bc8c1a7077e8f5ec122ba16b6dda1a5d3826057f67b9bcb lastState:map[] name:echoserver ready:true restartCount:0 started:true state:map[running:map[startedAt:2020-06-30T18:48:43Z]]]] hostIP:127.0.0.1 phase:Running podIP:172.17.0.3 podIPs:[map[ip:172.17.0.3]] qosClass:BestEffort startTime:2020-06-30T18:48:12Z]]}
		var template []byte
		if len(args) == 0 {
			template = []byte("{{printf \"Hello %s %s\\n\" .metadata.name .kind}}")
		} else {
			template = []byte("{{printf \"Hello %s %s %s\\n\" .metadata.name .kind .metadata.creationTimestamp}}")
		}
		printer, err := printers.NewGoTemplatePrinter([]byte(template))
		if err != nil {
			return err
		}
		return printer.PrintObj(obj, o.Out)
	}

	return nil
}

// Validate checks to the HelloKubernetesOptions to see if there is sufficient information run the command.
func (o HelloKubernetesOptions) Validate(cmd *cobra.Command, args []string) error {
	// fmt.Println("Printing args passed, ", args)
	if len(args) == 0 && cmdutil.IsFilenameSliceEmpty(o.FilenameOptions.Filenames, o.FilenameOptions.Kustomize) {
		return cmdutil.UsageErrorf(cmd, "Error: must specify one of -f and -k or type/name as arg\n\n")
	}
	return nil
}

// RunHelloKubernetes does the work
func (o HelloKubernetesOptions) RunHelloKubernetes(f cmdutil.Factory, cmd *cobra.Command) error {
	cmdNamespace, enforceNamespace, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	r := f.NewBuilder().
		Unstructured().
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &o.FilenameOptions).
		ResourceTypeOrNameArgs(false, o.args...).
		Flatten().
		Do()

	err = r.Err()
	if err != nil {
		return err
	}

	count := 0
	err = r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		if err := o.Recorder.Record(info.Object); err != nil {
			klog.V(4).Infof("error recording current command: %v", err)
		}

		count++

		return o.PrintObj(info.Object)
	})
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no objects passed to print")
	}
	return nil
}
