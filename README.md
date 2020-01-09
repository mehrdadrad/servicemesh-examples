# servicemesh-examples

![topo](/imgs/topo.png?raw=true "linkerd")


## LINKERD
#### install linkerd
```
curl -sL https://run.linkerd.io/install | sh
export PATH=$PATH:$HOME/.linkerd2/bin
linkerd install | kubectl apply -f -
```
#### deploy microservices
```
git clone https://github.com/mehrdadrad/servicemesh-examples.git
cd servicemesh-examples/k8s

kubectl apply -f backend.yaml
kubectl apply -f middleware.yaml
kubectl apply -f frontend.yaml
```
#### linkerd injection
```
kubectl get deploy/frontend-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
kubectl get deploy/middleware-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
kubectl get deploy/backend-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
```
#### pods status
```
kubectl get pods
NAME                                     READY   STATUS    RESTARTS   AGE
backend-deployment-7c65c5cdbf-6gs9z      2/2     Running   0          3h19m
backend-deployment-7c65c5cdbf-m2ckt      2/2     Running   0          3h19m
frontend-deployment-697fb6566-prpd7      2/2     Running   0          4m5s
frontend-deployment-697fb6566-z78sx      2/2     Running   0          4m12s
middleware-deployment-7599d78474-d4ftv   2/2     Running   0          4m24s
middleware-deployment-7599d78474-lxrjm   2/2     Running   0          4m11s
```
#### load test
```
kubectl get svc frontend-external
NAME                TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)          AGE
frontend-external   LoadBalancer   10.107.208.163   192.168.55.102   8082:32415/TCP   16h

while true; do curl 192.168.55.102:8082/time; sleep 2; done
```

#### linkerd cmd stat deployment
```
linkerd stat deploy
NAME                    MESHED   SUCCESS      RPS   LATENCY_P50   LATENCY_P95   LATENCY_P99   TCP_CONN
backend-deployment         2/2   100.00%   0.5rps          75ms          98ms         100ms          6
frontend-deployment        2/2   100.00%   0.5rps          75ms          98ms         100ms          2
middleware-deployment      2/2   100.00%   0.5rps          75ms          98ms         100ms          6
```

#### linkerd top

```
linkerd top deploy

(press q to quit)

Source                                  Destination                             Method      Path    Count    Best   Worst    Last  Success Rate
frontend-deployment-697fb6566-5s4rw     middleware-deployment-7599d78474-kn4hn  GET         /time      24    74ms   108ms   107ms       100.00%
frontend-deployment-697fb6566-5s4rw     middleware-deployment-7599d78474-wqw8r  GET         /time      24    78ms    96ms    88ms       100.00%
10.44.0.0                               frontend-deployment-697fb6566-5s4rw     GET         /time      24    87ms   111ms    93ms       100.00%
middleware-deployment-7599d78474-kn4hn  backend-deployment-5649869497-dtscv     GET         /time      22    70ms    97ms    83ms       100.00%
backend-deployment-5649869497-dtscv     12.23.14.17                             GET         /          22    68ms    89ms    72ms       100.00%
middleware-deployment-7599d78474-wqw8r  backend-deployment-5649869497-dtscv     GET         /time      22    70ms    87ms    77ms       100.00%
middleware-deployment-7599d78474-kn4hn  backend-deployment-5649869497-rvvwc     GET         /time      22    74ms    99ms    77ms       100.00%
backend-deployment-5649869497-rvvwc     12.23.14.17                             GET         /          22    70ms    93ms    76ms       100.00%
middleware-deployment-7599d78474-wqw8r  backend-deployment-5649869497-rvvwc     GET         /time      22    74ms    88ms    80ms       100.00%
frontend-deployment-697fb6566-57h48     middleware-deployment-7599d78474-wqw8r  GET         /time      20    77ms    98ms    84ms       100.00%
frontend-deployment-697fb6566-57h48     middleware-deployment-7599d78474-kn4hn  GET         /time      20    83ms   107ms    87ms       100.00%
10.44.0.0                               frontend-deployment-697fb6566-57h48     GET         /time      20    85ms   112ms    94ms       100.00%
```

#### linkerd tap
```
linkerd tap deploy/middleware-deployment

req id=17:0 proxy=in  src=10.32.0.15:56874 dst=10.32.0.20:8081 tls=true :method=GET :authority=middleware.default.svc.cluster.local:8081 :path=/time
req id=17:1 proxy=out src=10.32.0.20:46488 dst=10.32.0.17:8080 tls=true :method=GET :authority=backend.default.svc.cluster.local:8080 :path=/time
rsp id=17:1 proxy=out src=10.32.0.20:46488 dst=10.32.0.17:8080 tls=true :status=200 latency=114343µs
end id=17:1 proxy=out src=10.32.0.20:46488 dst=10.32.0.17:8080 tls=true duration=73µs response-length=100B
rsp id=17:0 proxy=in  src=10.32.0.15:56874 dst=10.32.0.20:8081 tls=true :status=200 latency=117368µs
end id=17:0 proxy=in  src=10.32.0.15:56874 dst=10.32.0.20:8081 tls=true duration=72µs response-length=22B
```
```
linkerd tap deploy/frontend-deployment

req id=15:0 proxy=in  src=10.44.0.0:57314 dst=10.32.0.15:8082 tls=not_provided_by_remote :method=GET :authority=192.168.55.102:8082 :path=/time
req id=15:1 proxy=out src=10.32.0.15:48382 dst=10.32.0.4:8081 tls=true :method=GET :authority=middleware.default.svc.cluster.local:8081 :path=/time
rsp id=15:1 proxy=out src=10.32.0.15:48382 dst=10.32.0.4:8081 tls=true :status=200 latency=89552µs
end id=15:1 proxy=out src=10.32.0.15:48382 dst=10.32.0.4:8081 tls=true duration=78µs response-length=22B
rsp id=15:0 proxy=in  src=10.44.0.0:57314 dst=10.32.0.15:8082 tls=not_provided_by_remote :status=200 latency=95348µs
end id=15:0 proxy=in  src=10.44.0.0:57314 dst=10.32.0.15:8082 tls=not_provided_by_remote duration=37µs response-length=35B
```
```
linkerd tap deploy/backend-deployment

req id=18:0 proxy=in  src=10.32.0.20:44304 dst=10.32.0.3:8080 tls=true :method=GET :authority=backend.default.svc.cluster.local:8080 :path=/time
req id=18:1 proxy=out src=10.32.0.3:34434 dst=12.23.14.17:80  tls=not_provided_by_service_discovery :method=GET :authority=time.jsontest.com :path=/
rsp id=18:1 proxy=out src=10.32.0.3:34434 dst=12.23.14.17:80  tls=not_provided_by_service_discovery :status=200 latency=71557µs
end id=18:1 proxy=out src=10.32.0.3:34434 dst=12.23.14.17:80  tls=not_provided_by_service_discovery duration=63µs response-length=108B
rsp id=18:0 proxy=in  src=10.32.0.20:44304 dst=10.32.0.3:8080 tls=true :status=200 latency=74675µs
end id=18:0 proxy=in  src=10.32.0.20:44304 dst=10.32.0.3:8080 tls=true duration=54µs response-length=100B
```








#### linkerd and grafana dashboards screenshots

```
linkerd dashboard &
```

![linkerd](/imgs/Linkerd.png?raw=true "linkerd")
![linkerd](/imgs/Linkerd_Control_Plane.png?raw=true "linkerd")
![linkerd](/imgs/Linkerd_Backend_Grafana.png?raw=true "linkerd")
![linkerd](/imgs/Linkerd_Middleware_Grafana.png?raw=true "linkerd")
