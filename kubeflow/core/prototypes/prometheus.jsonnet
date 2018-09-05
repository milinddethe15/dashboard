// @apiVersion 0.1
// @name io.ksonnet.pkg.prometheus
// @description Provides prometheus prototype in kubeflow core.
// @shortDescription Prometheus Service.
// @param name string Name for the component
// @param projectId string GCP project id.
// @param clusterName string GKE cluster name.
// @param zone string GKE cluster zone.

local prometheus = import "kubeflow/core/prometheus.libsonnet";
prometheus.new(env, params).list
