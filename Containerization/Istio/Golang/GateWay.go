package istio

import (
	"context"
	"devops2k8s/common"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GateWay() {
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
		gateway  *v1alpha3.Gateway
		services []*networkingv1alpha3.Server
	)

	service := &networkingv1alpha3.Server{
		Port: &networkingv1alpha3.Port{
			Number: 80,
			// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS.
			Protocol: "HTTP",
			Name:     "http",
		},
		// $hide_from_docs
		// The ip or the Unix domain socket to which the listener should be bound
		// to. Format: `x.x.x.x` or `unix:///path/to/uds` or `unix://@foobar`
		// (Linux abstract namespace). When using Unix domain sockets, the port
		// number should be 0.
		//Bind:                 "",
		Hosts: []string{"*.hzeng.com.cn"},
		Tls: &networkingv1alpha3.ServerTLSSettings{
			HttpsRedirect: true,
		},
	}

	services = append(services, service)
	gateway = &v1alpha3.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name:      "gw-bookinfo",
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.Gateway{
			Servers: services,
			//Selector:             nil,

		},
	}

	istioClient.NetworkingV1alpha3().Gateways(namespace).Create(context.TODO(), gateway, v1.CreateOptions{})

}
