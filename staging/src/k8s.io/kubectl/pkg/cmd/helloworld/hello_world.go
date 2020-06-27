package helloworld

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
)

// helloWorldOptions have the data required to perform the hello-world operation
type HelloWorldOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	PrintObj   printers.ResourcePrinterFunc

	// Filename options
	// resource.FilenameOptions
	RecordFlags *genericclioptions.RecordFlags

	// Common user flags

	// results of arg parsing
	Recorder                     genericclioptions.Recorder
	namespace                    string
	enforceNamespace             bool
	unstructuredClientForMapping func(mapping *meta.RESTMapping) (resource.RESTClient, error)

	genericclioptions.IOStreams
}

// NewHelloWorldOptions creates the options for hello-world
func NewHelloWorldOptions(ioStreams genericclioptions.IOStreams) *HelloWorldOptions {
	return &HelloWorldOptions{
		PrintFlags: genericclioptions.NewPrintFlags("hello-world").WithTypeSetter(scheme.Scheme),

		RecordFlags: genericclioptions.NewRecordFlags(),
		Recorder:    genericclioptions.NoopRecorder{},
		IOStreams:   ioStreams,
	}
}

// NewCmdHelloWorld creates the `hello-world` command
func NewCmdHelloWorld(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewHelloWorldOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "hello-world",
		DisableFlagsInUseLine: true,
		Short:                 "prints hello-world",
		Long:                  "prints hello-world",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.RunHelloWorld())
		},
	}

	// bind flag structs
	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	return cmd
}

// Complete adapts from the command line args and factory to the data required.
func (o *HelloWorldOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	o.RecordFlags.Complete(cmd)
	o.Recorder, err = o.RecordFlags.ToRecorder()
	if err != nil {
		return err
	}

	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	o.PrintObj = func(obj runtime.Object, out io.Writer) error {
		return printer.PrintObj(obj, out)
	}

	o.namespace, o.enforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	return nil
}

// Validate checks to the HelloWorldOptions to see if there is sufficient information run the command.
func (o HelloWorldOptions) Validate() error {
	return nil
}

// RunHelloWorld does the work
func (o HelloWorldOptions) RunHelloWorld() error {
	fmt.Println("Hello world")
	return nil
}
