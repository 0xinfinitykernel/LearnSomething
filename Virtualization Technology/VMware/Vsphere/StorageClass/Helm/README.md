## Enabling disk UUID on virtual machines
To resolve this issue, perform one of the below options:
* Create a virtual machine using a vSphere Web Client.
* After creating a virtual machine using the Host Client, add a virtual machine’s disk.EnableUUID attribute manually with following steps:
1. Open the Host Client, and log in to the ESXi.
2. Locate the Windows Server 2016 virtual machine for which you are enabling the disk UUID attribute, and power off the virtual machine.
3. After power-off, right-click the virtual machine, and choose Edit Settings.
4. Click VM Options tab, and select Advanced.
5. Click Edit Configuration in Configuration Parameters.
6. Click Add parameter.
7. In the Key column, type disk.EnableUUID.
8. In the Value column, type TRUE.
9. Click OK and click Save.
10. Power on the virtual machine

## Configure the Aggregation Layer
More information [here](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/)
```
  --allow-privileged=true \
  --requestheader-client-ca-file=/etc/kubernetes/ssl/ca.pem \
  --proxy-client-cert-file=/etc/kubernetes/ssl/kube-apiserver.pem \
  --proxy-client-key-file=/etc/kubernetes/ssl/kube-apiserver-key.pem \
  --requestheader-allowed-names=kubernetes \
  --requestheader-extra-headers-prefix=X-Remote-Extra- \
  --http2-max-streams-per-connection=3000 \
  --requestheader-group-headers=X-Remote-Group \
  --requestheader-username-headers=X-Remote-User \
  --enable-aggregator-routing=true \
```
## Configure the Kubelet
```
$ vi /usr/lib/systemd/system/kubelet.service
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/kubernetes/kubernetes
After=docker.service
Requires=docker.service

[Service]
WorkingDirectory=/var/lib/kubelet
ExecStart=/usr/local/bin/kubelet \
  --bootstrap-kubeconfig=/etc/kubernetes/kubelet-bootstrap.kubeconfig \
  --cert-dir=/etc/kubernetes/ssl \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --config=/etc/kubernetes/kubelet.json \
  --network-plugin=cni \
  --rotate-certificates \
  --pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.2 \
  --alsologtostderr=true \
  --logtostderr=false \
  --log-dir=/var/log/kubernetes \
  --minimum-image-ttl-duration=1525600m0s \
  --cloud-provider=external \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## The created PVC and node use the same datastore.
Install CPI
``` 
helm install vsphere-cpi \
     --namespace kube-system \
     ./vsphere-cpi/1.0.100 \
     --set vCenter.host=10.243.208.15 \
     --set vCenter.username=xxx@xxx.local \
     --set vCenter.password=1qaz@WSX3edc \
     --set vCenter.datacenters=Datacenter
``` 
Install CSI
``` 
helm install vsphere-csi \
     --namespace kube-system \
     ./vsphere-csi/2.3.1 \
     --set vCenter.host=10.243.208.15 \
     --set vCenter.username=xxx@xxx.local \
     --set vCenter.password=1qaz@WSX3edc \
     --set vCenter.clusterId=local \
     --set vCenter.datacenters=Datacenter \
     --set csiController.csiResizer.enabled=true \
     --set onlineVolumeExtend.enabled=true \
     --set storageClass.allowVolumeExpansion=true
```
## Migration

If using this chart to migrate volumes provisioned by the in-tree provider to the out-of-tree CPI + CSI, you need to taint all nodes with the following:
```
node.cloudprovider.kubernetes.io/uninitialized=true:NoSchedule
```

To perform this operation on all nodes in your cluster, the following script has been provided for your convenience:
```bash
# Node: If it returns no content, execute the last script
kubectl describe nodes | grep "ProviderID"
# Note: Since this script uses kubectl, ensure that you run `export KUBECONFIG=<path-to-kubeconfig-for-cluster>` before running this script
for node in $(kubectl get nodes | awk '{print $1}' | tail -n +2); do
	kubectl taint node $node node.cloudprovider.kubernetes.io/uninitialized=true:NoSchedule
done
```