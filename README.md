# servicemesh-examples



## LINKERD

```
cd k8s

kubectl apply -f backend.yaml
kubectl apply -f middleware.yaml
kubectl apply -f frontend.yaml

kubectl get deploy/frontend-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
kubectl get deploy/middleware-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
kubectl get deploy/backend-deployment -o yaml | linkerd inject - \ | kubectl apply -f -
```
