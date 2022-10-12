## Enabling disk UUID on virtual machines
To resolve this issue, perform one of the below options:
* Create a virtual machine using a vSphere Web Client.
* After creating a virtual machine using the Host Client, add a virtual machine’s disk.EnableUUID attribute manually with following steps:
1. Open the Host Client, and log in to the ESXi.
2. Locate the virtual machine for which you are enabling the disk UUID attribute, and power off the virtual machine.
3. After power-off, right-click the virtual machine, and choose Edit Settings.
4. Click VM Options tab, and select Advanced.
5. Click Edit Configuration in Configuration Parameters.
6. Click Add parameter.
7. In the Key column, type disk.EnableUUID.
8. In the Value column, type TRUE.
9. Click OK and click Save.
10. Power on the virtual machine
```azure
# kubectl describe po -n dev ess-master-0
    
  Warning  FailedAttachVolume  30s                attachdetach-controller  AttachVolume.Attach failed for volume "pvc-424f6d5b-181b-4b4d-a599-fc404bf455c3" : rpc error: code = Internal desc = failed to attach disk: "f07f0b78-b70b-455d-a031-820719c4010f" with node: "k8s-node06" err failed to attach cns volume: "f07f0b78-b70b-455d-a031-820719c4010f" to node vm: "VirtualMachine:vm-1082 [VirtualCenterHost: 172.16.50.220, UUID: 564daa78-9e44-7461-dd09-0af37cf3c25a, Datacenter: Datacenter [Datacenter: Datacenter:datacenter-1055, VirtualCenterHost: 172.16.50.220]]". fault: "(*types.LocalizedMethodFault)(0xc000dffe60)({\n DynamicData: (types.DynamicData) {\n },\n Fault: (types.CnsFault) {\n  BaseMethodFault: (types.BaseMethodFault) <nil>,\n  Reason: (string) (len=16) \"VSLM task failed\"\n },\n LocalizedMessage: (string) (len=32) \"CnsFault error: VSLM task failed\"\n})\n". opId: "07e31050"
```
if storgeclass "volumeMode": "Block", and pod return error "CnsFault error: VSLM task failed, you can try to add "ctkEnabled": "TRUE" in vmx file. [check here](https://kb.vmware.com/s/article/1020128)
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

Install CSI
``` 
helm install vsphere-csi \
     --namespace kube-system \
     ./vsphere-csi/2.3.1 \
     --set vCenter.host=172.16.50.220 \
     --set vCenter.username=wlhairou@vsphere.local \
     --set vCenter.password=RTa@EA12F3hairou \
     --set vCenter.clusterId=local \
     --set vCenter.datacenters=WULIU \
     --set csiController.csiResizer.enabled=true \
     --set onlineVolumeExtend.enabled=true \
     --set storageClass.allowVolumeExpansion=true
```
Install CPI
``` 
helm install vsphere-cpi \
     --namespace kube-system \
     ./vsphere-cpi/1.0.100 \
     --set vCenter.host=172.16.50.220 \
     --set vCenter.username=wlhairou@vsphere.local \
     --set vCenter.password=RTa@EA12F3hairou \
     --set vCenter.datacenters=WULIU
```
```
在没有master节点时，注意容忍和节点亲和的配置。
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