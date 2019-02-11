{
  parts(_env, _params):: {
    local params = _env + _params,

    local storage = import "kubeflow/pipeline/storage.libsonnet",
    local nfs = import "kubeflow/pipeline/nfs.libsonnet",
    local minio = import "kubeflow/pipeline/minio.libsonnet",
    local mysql = import "kubeflow/pipeline/mysql.libsonnet",
    local pipeline_apiserver = import "kubeflow/pipeline/pipeline-apiserver.libsonnet",
    local pipeline_scheduledworkflow = import "kubeflow/pipeline/pipeline-scheduledworkflow.libsonnet",
    local pipeline_persistenceagent = import "kubeflow/pipeline/pipeline-persistenceagent.libsonnet",
    local pipeline_viewercrd = import "kubeflow/pipeline/pipeline-viewercrd.libsonnet",
    local pipeline_ui = import "kubeflow/pipeline/pipeline-ui.libsonnet",

    local name = params.name,
    local namespace = params.namespace,
    local apiImage = params.apiImage,
    local scheduledWorkflowImage = params.scheduledWorkflowImage,
    local persistenceAgentImage = params.persistenceAgentImage,
    local viewerCrdControllerImage = params.viewerCrdControllerImage,
    local uiImage = params.uiImage,
    local nfsImage = params.nfsImage,
    local mysqlImage = params.mysqlImage,
    local minioImage = params.minioImage,
    local mysqlPvName = params.mysqlPvName,
    local minioPvName = params.minioPvName,
    local nfsPvName = params.nfsPvName,
    local mysqlPd = params.mysqlPd,
    local minioPd = params.minioPd,
    local nfsPd = params.nfsPd,
    nfs:: if (nfsPvName != "null") || (nfsPd != "null") then
             nfs.all(namespace, nfsImage)
           else [],
    local minioPvcName = if (nfsPvName != "null") || (nfsPd != "null") then "nfs-pvc" else "minio-pvc",
    all:: minio.all(namespace, minioImage, minioPvcName) +
          mysql.all(namespace, mysqlImage) +
          pipeline_apiserver.all(namespace, apiImage) +
          pipeline_scheduledworkflow.all(namespace, scheduledWorkflowImage) +
          pipeline_persistenceagent.all(namespace, persistenceAgentImage) +
          pipeline_viewercrd.all(namespace, viewerCrdControllerImage) +
          pipeline_ui.all(namespace, uiImage) +
          storage.all(namespace, mysqlPvName, minioPvName, nfsPvName, mysqlPd, minioPd, nfsPd) +
          $.parts(_env, _params).nfs,
  },
}
