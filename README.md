# ingress-group-validator

Validating admission webhook for Ingress Group. https://kubernetes-sigs.github.io/aws-load-balancer-controller/guide/ingress/annotations/#ingressgroup

## Usage

### Deploy webhook server to Kind

```console
$ kind create cluster
$ make docker-build
$ kind load docker-image controller:latest
$ kustomize build config/default | kubectl apply -f -
```

Ingress Group Validator checks the value of `alb.ingress.kubernetes.io/group.name` annotation. The accepted value format is `<namespace>/<group-name>`

```console
$ cat e2etest/invalid-ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: invalid
  namespace: default
  annotations:
    alb.ingress.kubernetes.io/group.name: test
spec:
  rules:
    - http:
        paths:
          - path: /testpath
            pathType: Prefix
            backend:
              serviceName: test
              servicePort: 8080

$ kubectl apply -f e2etest/invalid-ingress.yaml
Resource: "networking.k8s.io/v1beta1, Resource=ingresses", GroupVersionKind: "networking.k8s.io/v1beta1, Kind=Ingress"
Name: "invalid", Namespace: "default"
for: "e2etest/invalid-ingress.yaml": admission webhook "vingress.kb.io" denied the request: deny 'invalid' since the alb.ingress.kubernetes.io/group.name annotation does not start with 'default/'
```

```console
$ cat e2etest/valid-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: valid
  namespace: default
  annotations:
    alb.ingress.kubernetes.io/group.name: test
spec:
  rules:
    - http:
        paths:
          - path: /testpath
            pathType: Prefix
            backend:
              service:
                name: test
                port:
                  number: 80

$ kubectl apply -f e2etest/valid-ingress.yaml
ingress.networking.k8s.io/valid created

```
