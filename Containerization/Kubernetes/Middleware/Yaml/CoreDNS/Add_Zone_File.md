# Add a Zone File to Kubernetes CoreDNS

# Overview
The zone file is added as a ConfigMap entry which will be projected in CoreDNS pods as a file, under the zone file name. The Corefile projected as part of the same ConfigMap should be also modified to refer the new zone file with the "file" directive.

The CoreDNS deployment is then scaled down, the new configuration file is added as an "item" in configMap volume mount, and the deployment is then scaled up.

# Procedure
## Add the Zone File to ConfigMap
Get the content of the coredns ConfigMap "Corefile" entry:

```
kubectl -n kube-system get configmap coredns -o jsonpath='{.data.Corefile}' > ./Corefile
```
You should get something similar to:

```
.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       upstream
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
```
Add the following configuration extension:

```
.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       upstream
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
    file /etc/coredns/blue-zone.db blue.test {
      upstream
    }
   
}
```
This will add a zone file for the "blue.test" domain.

In the same directory, add a "blue-zone.db" file with the following content:
```

; blue.test zone
blue.test.                   IN          SOA         sns.dns.icann.org.  noc.dns.icann.org. 2019101701 7200 3600 1209600 3600
blue.test.                   IN          NS          b.iana-servers.net.
blue.test.                   IN          NS          b.iana-servers.net.
blue.test.                   IN          A           127.0.0.1
something.blue.test.         IN          CNAME       myservice.svc.cluster.local.
```
Update the ConfigMap with the new content. From the directory that contains Corefile and blue-zone.db:

```
kubectl -n kube-system create configmap coredns --from-file=Corefile --from-file=blue-zone.db --save-config=true --dry-run -o yaml > coredns.yaml
kubectl -n kube-system apply -f ./coredns.yaml
```
### Wildcard Domain
To configure a wildcard domain, use this zone file:

```
; blue.test zone
blue.test.                   IN          SOA         sns.dns.icann.org.  noc.dns.icann.org. 2019101701 7200 3600 1209600 3600
blue.test.                   IN          NS          b.iana-servers.net.
blue.test.                   IN          NS          b.iana-servers.net.
blue.test.                   IN          A           127.0.0.1
*                            IN          CNAME       myservice.svc.cluster.local.
```
### Edit the coredns Deployment
```
kubectl -n kube-system edit deployment coredns
```
In the "volumes" section, add the following key/path pair:

```
volumes:
- configMap:
    defaultMode: 420
    items:
    - key: Corefile
      path: Corefile
    - key: blue-zone.db
      path: blue-zone.db
```
### Scale Down and Up the coredns Deployment
```
kubectl -n kube-system scale --replicas=0 deployment coredns
kubectl -n kube-system scale --replicas=2 deployment coredns
```
Make sure the coredns pods start fine:

```
coredns-7f8f4bd796-khdgq                 1/1     Running   0          8s
coredns-7f8f4bd796-vbkhq                 1/1     Running   0          8s
```
# Example K8s ConfigMap
```
blue-zone.db: |
; dev.com zone
dev.com.        IN  SOA dns.dev.com. dns2.dev.com. 2015082541 7200 3600 1209600 3600
dns.dev.com.    IN  A   172.18.81.36
dns2.dev.com.    IN  A   172.18.81.36
rancher.dev.com.    IN  A   172.18.81.36
*.dev.com.    IN  A   172.18.81.36
```