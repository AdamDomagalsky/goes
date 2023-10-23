  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    annotations:
      alb.ingress.kubernetes.io/backend-protocol-version: GRPC
      alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/target-type: ip
    name: argocd
    namespace: argocd
  spec:
    ingressClassName: alb
    rules:
    - http:
        paths:
        - path: /
          backend:
            service:
              name: argogrpc
              port:
                number: 80
          pathType: Prefix
        - path: /
          backend:
            service:
              name: argocd-server
              port:
                number: 80
          pathType: Prefix

# apiVersion: networking.k8s.io/v1
# kind: Ingress
# metadata:
#   annotations:
#     alb.ingress.kubernetes.io/target-type: ip
#     alb.ingress.kubernetes.io/scheme: internet-facing
#     alb.ingress.kubernetes.io/backend-protocol-version: HTTP
#     alb.ingress.kubernetes.io/backend-protocol: HTTP
#     alb.ingress.kubernetes.io/load-balancer-attributes: routing.http2.enabled=true
#     alb.ingress.kubernetes.io/success-codes: '200-307'
#   name: argocd
#   namespace: argocd
# spec:
#   ingressClassName: alb
#   rules:
#   - http:
#       paths:
#       - path: /
#         backend:
#           service:
#             name: argocd-server
#             port:
#               number: 80
#         pathType: Prefix