package istio

import (
	"context"
	"devops2k8s/common"
	"devops2k8s/jenkins"
	"fmt"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func VirtualService(codeGroups string, projectName string, nameSpace string) error {
	restConfig, err := common.GetRestConfIstio()
	if err != nil {
		return err
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	var (
		httpRouteList            []*networkingv1alpha3.HTTPRoute
		HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
	)
	// 定义http路由
	HTTPRouteDestination := &networkingv1alpha3.HTTPRouteDestination{
		Destination: &networkingv1alpha3.Destination{
			Host:   "mlj-" + projectName,
			Subset: "v2",
		},
		// 定义权重
		Weight: 100,
	}
	HTTPRouteDestinationList = append(HTTPRouteDestinationList, HTTPRouteDestination)
	httpRouteSign := networkingv1alpha3.HTTPRoute{

		Route: HTTPRouteDestinationList,
	}
	httpRouteList = append(httpRouteList, &httpRouteSign)
	virtualService := &v1alpha3.VirtualService{
		ObjectMeta: v1.ObjectMeta{
			Name:      projectName, // 定义vs的名称
			Namespace: nameSpace,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{projectName + ".io.mlj162.com"}, // 定义可访问的hosts
			Gateways: []string{"limited-time-offers-front"},
			Http:     httpRouteList, // 为hosts 绑定路由
		},
	}
	// 创建VS
	_, err = istioClient.NetworkingV1alpha3().VirtualServices(nameSpace).Create(context.TODO(), virtualService, v1.CreateOptions{})
	//vs, err := istioClient.NetworkingV1alpha3().VirtualServices(nameSpace).Create(context.TODO(), virtualService, v1.CreateOptions{})
	if err != nil {
		return err
	}
	// 打印VS
	//log.Print(vs)
	fmt.Println("VirtualService Created successfully!")
	webhookURL, err := jenkins.CreateWebHookURL()
	if err != nil {
		return err
	}
	return jenkins.Jobs(projectName, codeGroups, webhookURL)
	// kubectl get vs | grep mljtest | awk '{print $1}' | xargs -I {} kubectl delete vs/{}
}
