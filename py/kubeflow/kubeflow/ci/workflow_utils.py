import logging
import os

from kubeflow.testing import argo_build_util

# The name of the NFS volume claim to use for test files.
NFS_VOLUME_CLAIM = "nfs-external"
# The name to use for the volume to use to contain test data
DATA_VOLUME = "kubeflow-test-volume"

E2E_DAG_NAME = "e2e"
EXIT_DAG_NAME = "exit-handler"


LOCAL_TESTING = os.getenv("LOCAL_TESTING", "False")
CREDENTIALS_VOLUME = {"name": "aws-secret",
                      "secret": {"secretName": "aws-secret"}}
CREDENTIALS_MOUNT = {"mountPath": "/root/.aws/",
                     "name": "aws-secret"}

AWS_WORKER_IMAGE = "public.ecr.aws/j1r0q0g6/kubeflow-testing:latest"


class ArgoTestBuilder:
    def __init__(self, name=None, namespace=None, bucket=None,
                 test_target_name=None,
                 **kwargs):
        self.name = name
        self.namespace = namespace
        self.bucket = bucket
        self.template_label = "argo_test"
        self.test_target_name = test_target_name
        self.mkdir_task_name = "make-artifacts-dir"

        # *********************************************************************
        #
        # Define directory locations
        #
        # *********************************************************************

        # mount_path is the directory where the volume to store the test data
        # should be mounted.
        self.mount_path = "/mnt/" + "test-data-volume"
        # test_dir is the root directory for all data for a particular test
        # run.
        self.test_dir = self.mount_path + "/" + self.name
        # output_dir is the directory to sync to GCS to contain the output for
        # this job.
        self.output_dir = self.test_dir + "/output"

        self.artifacts_dir = "%s/artifacts/junit_%s" % (self.output_dir, name)

        # source directory where all repos should be checked out
        self.src_root_dir = "%s/src" % self.test_dir
        # The directory containing the kubeflow/kubeflow repo
        self.src_dir = "%s/kubeflow/kubeflow" % self.src_root_dir

        # Root of testing repo.
        self.testing_src_dir = os.path.join(self.src_root_dir,
                                            "kubeflow/testing")

        # Top level directories for python code
        self.kubeflow_py = self.src_dir

        # The directory within the kubeflow_testing submodule containing
        # py scripts to use.
        self.kubeflow_testing_py = "%s/kubeflow/testing/py" % self.src_root_dir

        self.go_path = self.test_dir

    def _build_workflow(self):
        """Create a scaffolding CR for the Argo workflow"""
        volumes = [{
            "name": DATA_VOLUME,
            "persistentVolumeClaim": {
                "claimName": NFS_VOLUME_CLAIM
            },
        }]
        if LOCAL_TESTING == "False":
            volumes.append(CREDENTIALS_VOLUME)

        workflow = {
            "apiVersion": "argoproj.io/v1alpha1",
            "kind": "Workflow",
            "metadata": {
                "name": self.name,
                "namespace": self.namespace,
                "labels": argo_build_util.add_dicts([
                    {
                        "workflow": self.name,
                        "workflow_template": self.template_label,
                    },
                    argo_build_util.get_prow_labels()
                ]),
            },
            "spec": {
                "entrypoint": E2E_DAG_NAME,
                # Have argo garbage collect old workflows otherwise we overload
                # the API server.
                "volumes": volumes,
                "onExit": EXIT_DAG_NAME,
                "templates": [
                    {
                        "dag": {
                            "tasks": []
                        },
                        "name": E2E_DAG_NAME
                    },
                    {
                        "dag": {
                            "tasks": []
                        },
                        "name":EXIT_DAG_NAME
                    },
                ],
            },  # spec
        }  # workflow

        return workflow

    def build_task_template(self):
        """Return a template for all the tasks"""
        volume_mounts = [{
            "mountPath": "/mnt/test-data-volume",
            "name": DATA_VOLUME
        }]
        if LOCAL_TESTING == "False":
            volume_mounts.append(CREDENTIALS_MOUNT)

        image = AWS_WORKER_IMAGE

        task_template = {
            "activeDeadlineSeconds": 3000,
            "container": {
                "command": [],
                "env": [],
                "image": image,
                "imagePullPolicy": "Always",
                "name": "",
                "resources": {
                    "limits": {
                        "cpu": "4",
                        "memory": "4Gi"
                    },
                    "requests": {
                        "cpu": "1",
                        "memory": "1536Mi"
                    },
                },
                "volumeMounts": volume_mounts,
            },
            "metadata": {
                "labels": {
                    "workflow_template": self.template_label,
                }
            },
            "outputs": {},
        }

        # Define common environment variables to be added to all steps
        common_env = [
            {
                "name": "PYTHONPATH",
                "value": ":".join([self.kubeflow_py, self.kubeflow_testing_py])
            },
            {
                "name": "GOPATH",
                "value": self.go_path
            },
        ]

        task_template["container"]["env"].extend(common_env)

        if self.test_target_name:
            task_template["container"]["env"].append(
                {
                    "name": "TEST_TARGET_NAME",
                    "value": self.test_target_name
                }
            )

        task_template = argo_build_util.add_prow_env(task_template)

        return task_template

    def _create_checkout_task(self, task_template):
        """Checkout the kubeflow/testing and kubeflow/kubeflow code"""
        main_repo = argo_build_util.get_repo_from_prow_env()
        if not main_repo:
            logging.info("Prow environment variables for repo not set")
            main_repo = "kubeflow/testing@HEAD"
        logging.info("Main repository: %s", main_repo)
        repos = [main_repo]

        checkout = argo_build_util.deep_copy(task_template)

        checkout["name"] = "checkout"
        checkout["container"]["command"] = [
            "/usr/local/bin/checkout_repos.sh",
            "--repos=" + ",".join(repos),
            "--src_dir=" + self.src_root_dir,
        ]

        return checkout

    def _create_make_dir_task(self, task_template):
        """Create the directory to store the artifacts of each task"""
        # (jlewi)
        # pytest was failing trying to call makedirs. My suspicion is its
        # because the two steps ended up trying to create the directory at the
        # same time and classing. So we create a separate step to do it.
        mkdir_step = argo_build_util.deep_copy(task_template)

        mkdir_step["name"] = self.mkdir_task_name
        mkdir_step["container"]["command"] = ["mkdir", "-p",
                                              self.artifacts_dir]

        return mkdir_step

    def build_init_workflow(self):
        """Build the Argo workflow graph"""
        workflow = self._build_workflow()
        task_template = self.build_task_template()

        # checkout the code
        checkout_task = self._create_checkout_task(task_template)
        argo_build_util.add_task_to_dag(workflow, E2E_DAG_NAME, checkout_task,
                                        [])

        # create the artifacts directory
        mkdir_task = self._create_make_dir_task(task_template)
        argo_build_util.add_task_to_dag(workflow, E2E_DAG_NAME, mkdir_task,
                                        [checkout_task["name"]])

        return workflow

    # the following methods should be implemented from the test cases
    def build(self):
        """Build the Argo Worfklow for this test"""
        raise NotImplementedError("Subclasses should implement this!")

    def create_workflow(name=None, namespace=None, bucket=None, **kwargs):
        """Return the final dict with the Argo Workflow to be submitted"""
        raise NotImplementedError("Subclasses should implement this!")
