# Tekton Pipeline to manually trigger a build
This pipeline will:
- Clone in source code from Git
- Run Golang tests against the repo
- Create/update an image build and save the build image
- Create/update a Kubernetes deployment using the built image
All task containers and the deployment will run in the current namespace.</br>
The build namespace can be customised in the PipelineRun yaml file.

## Install Tekton

```
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```
### (optional) Install Tekton Dashboard
Update line 9 of `tekton/dashboard/dashboard-ingress.yml` to have a hostname that will work behind your ingress controller.
```
kubectl apply --filename https://storage.googleapis.com/tekton-releases/dashboard/latest/tekton-dashboard-release.yaml
kubectl apply -f tekton/dashboard/dashboard-ingress.yml
```

## Pipeline setup

All commands should be run from the tekton directory.

### Export env vars
```
export REGISTRY_USERNAME=admin
export REGISTRY_PASSWORD=####
export REGISTRY_ENDPOINT=harbor.lab
export PIPELINE_NAMESPACE=tbs-tekton-validation
export GIT_SSH_KEY_LOCATION=$HOME/.ssh/git_key
```

### Dependent Resources
Use `ytt` to inject variables into `service_account.yaml`.
```
ytt -f dependencies/service_account.yml \
  -v pipeline_namespace="${PIPELINE_NAMESPACE:?}" | kubectl apply -f-
kubectl apply -f dependencies/pvc.yml -n $PIPELINE_NAMESPACE
```

### Git Secret

```
# Update $HOME/.ssh/git_key to location of git private key
kubectl create secret generic ssh-key --from-file=ssh-privatekey=${GIT_SSH_KEY_LOCATION} --type kubernetes.io/ssh-auth -n $PIPELINE_NAMESPACE
kubectl annotate secret ssh-key tekton.dev/git-0='github.com' -n $PIPELINE_NAMESPACE
```

### Reg Secret
```
kp secret create my-registry-creds --registry $REGISTRY_ENDPOINT --registry-user $REGISTRY_USERNAME -n $PIPELINE_NAMESPACE
```


### Canned tasks
```
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/golang-test/0.2/golang-test.yaml -n $PIPELINE_NAMESPACE
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.5/git-clone.yaml -n $PIPELINE_NAMESPACE
```

### Local tasks
```
kubectl apply -f tasks -n $PIPELINE_NAMESPACE
```

### Apply pipeline
```
kubectl apply -f pipelines/golang-pipeline.yml -n $PIPELINE_NAMESPACE
```

### Create a pipeline run
Runs must be created to allow name to be generated.</br>
Update the PipelineRun file `tekton\pipelines\golang-pipeline-run.yml` to ensure the values reflect your environment.
```
kubectl create -f pipelines/golang-pipeline-run.yml -n $PIPELINE_NAMESPACE
```