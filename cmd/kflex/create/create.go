package create

import (
	"context"
	"fmt"
	"os"
	"sync"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tenancyv1alpha1 "mcc.ibm.org/kubeflex/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"sigs.k8s.io/controller-runtime/pkg/client"

	"mcc.ibm.org/kubeflex/cmd/kflex/common"
	cont "mcc.ibm.org/kubeflex/cmd/kflex/ctx"
	"mcc.ibm.org/kubeflex/pkg/certs"
	kfclient "mcc.ibm.org/kubeflex/pkg/client"
	"mcc.ibm.org/kubeflex/pkg/kubeconfig"
	"mcc.ibm.org/kubeflex/pkg/util"
)

type CPCreate struct {
	common.CP
}

func (c *CPCreate) Create() {
	done := make(chan bool)
	var wg sync.WaitGroup
	cx := cont.CPCtx{}
	cx.Context()

	cl := kfclient.GetClient(c.Kubeconfig)

	cp := c.generateControlPlane()

	util.PrintStatus(fmt.Sprintf("Creating new control plane %s...", c.Name), done, &wg)
	if err := cl.Create(context.TODO(), cp, &client.CreateOptions{}); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
		os.Exit(1)
	}
	done <- true

	clientset := kfclient.GetClientSet(c.Kubeconfig)

	util.PrintStatus("Waiting for API server to become ready...", done, &wg)
	kubeconfig.WatchForSecretCreation(clientset, c.Name, certs.AdminConfSecret)

	if err := util.WaitForDeploymentReady(clientset, "kube-apiserver", util.GenerateNamespaceFromControlPlaneName(cp.Name)); err != nil {
		fmt.Fprintf(os.Stderr, "Error waiting for deployment to become ready: %v\n", err)
		os.Exit(1)
	}
	done <- true

	if err := kubeconfig.LoadAndMerge(c.Ctx, clientset, c.Name); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading and merging kubeconfig: %v\n", err)
		os.Exit(1)
	}

	wg.Wait()
}

func (c *CPCreate) generateControlPlane() *tenancyv1alpha1.ControlPlane {
	return &tenancyv1alpha1.ControlPlane{
		ObjectMeta: v1.ObjectMeta{
			Name: c.Name,
		},
	}
}