# Install Steps
These steps follow the [documented process](https://docs.vmware.com/en/Tanzu-Build-Service/1.3/vmware-tanzu-build-service-v13/GUID-installing.html#relocate-images-to-a-registry) to install TBS 1.3.4 based on images in a local registry. 

## Create a Cluster
Note `/var/lib/containerd` must be extended to allow for a large image cache.</br>

### Instructions for vSphere with Tanzu.
Update `tbs-cluster.yml` to have the correct storage policy and target namespace for the workload cluster. </br>
When logged into the supervisor cluster:
```
# Create the cluster
kubectl apply -f tbs-cluster.yml
# Check that cluster has become ready
kubectl get tkc -A
# Log in to the workload cluster
kubectl vsphere login --server=172.20.2.2 --vsphere-username administrator@vsphere.local --tanzu-kubernetes-cluster-name tbs-1 --tanzu-kubernetes-cluster-namespace tbs --insecure-skip-tls-verify
# Clear the default PSP
kubectl create clusterrolebinding psp:authenticated  --clusterrole=psp:vmware-system-privileged --group=system:authenticated
```

## Setup TBS
Ensure the cluster has the registry CA added to the system trust store. In the example above it is defined in the yaml file.</br>
Setup a project in the registry host TBS images and another for apps.</br>
Follow the docs... </br>

Update and export necessary env vars:
```
# REGISTRY refers to your container registry which will host all the images.
export REGISTRY_USERNAME=admin
export REGISTRY_PASSWORD=####
export REGISTRY_ENDPOINT=harbor.lab
# The path is the <registry endpoint>/<project name>
export REGISTRY_TBS_PATH=${REGISTRY_ENDPOINT}/tbs
export REGISTRY_APPS_PATH=${REGISTRY_ENDPOINT}/test-apps
export TANZU_NET_USERNAME=user@vmware.com
export TANZU_NET_PASSWORD=####
```

Import TBS bundle into a local registry.
```
imgpkg copy -b "registry.tanzu.vmware.com/build-service/bundle:1.3.4" --to-repo ${REGISTRY_TBS_PATH}/dependencies
```
Pull the bundle into the local temp directory.
```
imgpkg pull -b "${REGISTRY_TBS_PATH}/dependencies:1.3.4" -o /tmp/bundle
```

Install the TBS kapp in the cluster.
```
ytt -f /tmp/bundle/values.yaml \
    -f /tmp/bundle/config/ \
    -f /home/matt/workspace/secrets/certs/rootCA.pem \
	-v kp_default_repository="${REGISTRY_TBS_PATH}/dependencies" \
	-v kp_default_repository_username=$REGISTRY_USERNAME \
	-v kp_default_repository_password=$REGISTRY_PASSWORD \
	-v pull_from_kp_default_repo=true \
	-v tanzunet_username="$TANZU_NET_USERNAME" \
	-v tanzunet_password="$TANZU_NET_PASSWORD" \
	| kbld -f /tmp/bundle/.imgpkg/images.yml -f- \
	| kapp deploy -a tanzu-build-service -f- -y
```

## Validation
Set the registry Secret.
```
kp secret create my-registry-creds --registry $REGISTRY_ENDPOINT --registry-user $REGISTRY_USERNAME
```

Golang test app
```
kp image create test-image-go --tag ${REGISTRY_APPS_PATH}/test-app-go --git https://github.com/laidbackware/tanzu-build-service-experiments --sub-path ./example-apps/golang --wait
```

Java test app
```
kp image create my-image --tag ${REGISTRY_APPS_PATH}/test-app --git https://github.com/buildpacks/samples --sub-path ./apps/java-maven --wait
```

Python test app
```
kp image create test-image-python --tag ${REGISTRY_APPS_PATH}/test-app-python --git https://github.com/laidbackware/tanzu-build-service-experiments --sub-path ./example-apps/python --wait
```