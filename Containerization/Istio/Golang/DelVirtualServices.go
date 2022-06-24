package istio

import (
	"context"
	"devops2k8s/common"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DelVirtualServices() {
	namespace := "bookinfo"
	restConfig, err := common.GetRestConf()
	if err != nil {
		return
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return
	}
	// 删除啊 vs
	istioClient.NetworkingV1alpha3().VirtualServices(namespace).Delete(context.TODO(), "vs-test", v1.DeleteOptions{})

}
