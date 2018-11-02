// @apiVersion 0.1
// @name io.ksonnet.pkg.tf-serving-service
// @description TensorFlow serving
// @shortDescription A TensorFlow serving model
// @param name string Name to give to each of the components
// @optionalParam serviceType string ClusterIP The k8s service type for tf serving.
// @optionalParam modelName string null The model name
// @optionalParam trafficRule string v1:100 The traffic rule, in the format of version:percentage,version:percentage,..
// @optionalParam injectIstio string false Whether to inject istio sidecar; should be true or false.

local k = import "k.libsonnet";
local tfservingService = import "kubeflow/tf-serving/tf-serving-service-template.libsonnet";
local util = import "kubeflow/tf-serving/util.libsonnet";

tfservingService.new(env, params).all
