### 查看容器内mysql日志并过滤

k exec -ti mysql-0 -n tp -- bash -c '/opt/bitnami/mysql/bin/mysqlbinlog --no-defaults /bitnami/mysql/data/binlog.000170 --base64-output=decode-rows -vv --skip-gtids=true  | grep -C 1 -i "F_2022102741002"'

### 端口流量异常找不到原因时

iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -X