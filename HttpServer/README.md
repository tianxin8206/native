### 模块十作业

安装helm步骤略过……

安装loki、promethus
```shell
helm pull grafana/loki-stack

tar -xvf loki-stack-2.4.1.tgz

sed -i s#rbac.authorization.k8s.io/v1beta1#rbac.authorization.k8s.io/v1#g $(find ./loki-stack -name "*.yaml")

helm upgrade --install loki ./loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

Grafana Dashboard:
![image](https://user-images.githubusercontent.com/32237514/146003464-c6cae003-78c7-4d6a-9ecf-8e0b481648e3.png)

Promethus界面:
![image](https://user-images.githubusercontent.com/32237514/146003664-811da5b9-7f46-45ec-a4df-a165c2cbd0e7.png)
