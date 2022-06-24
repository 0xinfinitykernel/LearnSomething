package istio

import (
	"context"
	"devops2k8s/common"
	"log"

	"istio.io/api/networking/v1beta1"
	networkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	namespace := "bookinfo"
	restConfig, err := common.GetRestConf()
	if err != nil {
		return
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return
	}

	var (
		serviceEntry *networkingv1beta1.ServiceEntry
		ports        []*v1beta1.Port
	)

	ports = append(ports, &v1beta1.Port{
		Number:   80,
		Protocol: "HTTP",
		Name:     "http",
	})
	ports = append(ports, &v1beta1.Port{
		Number:   443,
		Protocol: "HTTPS",
		Name:     "https",
	})
	serviceEntry = &networkingv1beta1.ServiceEntry{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "baidu",
		},
		Spec: v1beta1.ServiceEntry{
			//hosts字段用于在VirtualServices和DestinationRules中选择匹配的主机
			//对于HTTP traffic， Host/Authority header 与主机字段匹配
			//对于包含服务器名称指示(SNI)的HTTPs或TLS通信,sni的值与host进行匹配
			Hosts:      []string{"www.baidu.com"},
			Ports:      ports,
			Resolution: v1beta1.ServiceEntry_DNS,
		},
	}

	_, err = istioClient.NetworkingV1beta1().ServiceEntries(namespace).Create(context.TODO(), serviceEntry, v1.CreateOptions{})

	if err != nil {
		log.Print(err)
	}
}
