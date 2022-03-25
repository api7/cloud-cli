Deploy APISIX on Kubernetes
=======================

In this section, you'll learn how to deploy APISIX on Kubernetes through Cloud CLI.

> Note, before you go ahead, make sure you read the section
> [How to Configure Cloud CLI](./configuring-cloud-cli.md)

Cloud CLI will create [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment), 
[Service](https://kubernetes.io/docs/concepts/services-networking/service), 
[ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap) and 
[Secret](https://kubernetes.io/docs/concepts/configuration/secret)
on Kubernetes for APISIX, each resource provides different functionality.

* The Cloud Lua Module is stored in the ConfigMap (default name is cloud-module).

The Cloud Lua Module contains codes to communicate with API7 Cloud (such as
heartbeat, status reporting, etc.), it'll be downloaded every time you run the command.

> Currently, the Cloud Lua Module will be downloaded from [api7/cloud-scripts](https://github.com/api7/cloud-scripts).

* TLS Bundle is stored in the Secret (default name is cloud-ssl).

TLS Bundle (Certificate, Private Key, CA Bundle) will be downloaded from API7
Cloud, only instances with a valid client certificate can be connected to API7 Cloud.

> See the
> [DP Certificate API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getCertificates)
> to learn the details.

Cloud CLI deploys APISIX on Kubernetes by using [helm](https://helm.sh/), so please make sure helm was installed before you go ahead.
the essential parts that APISIX needs to run, the configuration items in values are referenced into the deployment.yaml.

> See [Helm Values Template API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getControlPlaneStartupConfig)
> you can get a value.yaml template of the helm.


Run Command
-----------

```shell
cloud-cli deploy kubernetes \
  --name my-apisix \
  --namespace apisix \
  --replica-count 1 \
  --apisix-image apache/apisix:2.11.0-centos \
  --helm-install-arg --output=table

Congratulations! Your APISIX cluster was deployed successfully on Kubernetes.
The Helm release name is: my-apisix
The APISIX Deployment name is: "my-apisix"
The APISIX Service name is: "my-apisix-gateway"

Workloads:
Pod Name: my-apisix-7959ffd978-bmlv8 APISIX ID: e9ecb37c-6631-49ef-9990-bc1370278834
```

In this command, we:

1. name the helm release to `my-apisix`;
2. specify the namespace is `apisix`;
3. specify the APISIX pods replica is `1`;
4. specify the APISIX image `apache/apisix:2.11.0-centos`;
5. prints the output in the table format for helm install command.

And the following operations were done in the above command:

1. create helm release that name is `my-apisix`;
2. create namespace on Kubernetes that name is `apisix`, if it not already existed;
3. create secret with name is `cloud-ssl` on namespace which name is `apisix`, if it not already existed;
4. create configMap with name is `cloud-module` on namespace which name is `apisix`, if it not already existed;
5. create Deployment, Service, Pod on namespace which name is `apisix`.

If you see the similar output about the Helm release name, APISIX Deployment name, APISIX Service name, Pod Name and APISIX ID, then your
APISIX instance was deployed successfully. You can redirect to the API7 Cloud console
to check the status of your APISIX cluster.
![img.png](./deploy-apisix-on-kubernetes-succeed.png)

> You can also run the `kubectl get` command to check the status for this deployment.

Besides, you can go into the Kubernetes and access APISIX cluster through by service or pods.

Stop Instance
-------------

If you want to stop the APISIX instance, just run the command below:

```shell
cloud-cli stop kubernetes \
  --name my-apisix \
  --namespace apisix
```

In this command, the following operations will be done:

1. delete helm release that name is `my-apisix`, it will be delete the Deployment and Service;
2. delete secret with name is `cloud-ssl` on namespace which name is `apisix`;
3. delete configMap with name is `cloud-module` on namespace which name is `apisix`.


Command Option Reference
------------------------

You can run `cloud-cli deploy kubernetes --help` or `cloud-cli stop kubernetes --help` to learn 
the command line option meanings.