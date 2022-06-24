package istio

import (
	"context"
	"devops2k8s/common"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DelServiceEntry() {
	namespace := "bookinfo"
	restConfig, err := common.GetRestConf()
	if err != nil {
		return
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return
	}
	istioClient.NetworkingV1beta1().ServiceEntries(namespace).Delete(context.TODO(), "baidu", v1.DeleteOptions{})
}
