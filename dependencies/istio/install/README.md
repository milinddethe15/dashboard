# Istio manifest for Kubeflow

The manifest is taken from OSS Istio pre-release
[1.1-20190111-09-15](https://gcsweb.istio.io/gcs/istio-prerelease/daily-build/release-1.1-20190111-09-15/).

We need custom configuration as discussed in the
[issue](https://github.com/kubeflow/kubeflow/issues/1909#issuecomment-438409215).

Specifically, the things we changed are:

1. The service type of `istio-ingressgateway` is changed from `LoadBalancer` to `NodePort`.
1. Add annotation to service `istio-ingressgateway`: `beta.cloud.google.com/backend-config: XXX`.
   This is to [enable IAP](https://cloud.google.com/iap/docs/enabling-kubernetes-howto#kubernetes-configure).

*Notice*: 
* The policy of configmap `istio-sidecar-injector` is `enabled`. To turn on injection for certain namespace, add namespace 
label `istio-injection: enabled`. More details see [table](https://github.com/istio/istio/issues/6476#issuecomment-399219937).

*TODO*: To allow egress, we need to know cluster specific IP ranges. Also, Istio's recommended way is to use
serviceEntry. We will figure out this part later.
