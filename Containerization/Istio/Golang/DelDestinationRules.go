package istio

import (
	"context"
	"devops2k8s/common"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DelDestinationRules() {
	namespace := "bookinfo"
	restConfig, err := common.GetRestConf()
	if err != nil {
		return
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return
	}
	istioClient.NetworkingV1alpha3().DestinationRules(namespace).Delete(context.TODO(), "dr-reviews", v1.DeleteOptions{})

}
