# App Deploy

The aim of this application is to deploy deploy any application by receiving config through handler.
After received the config it could execute pre-deploy and post-deploy scripts defined on its own configuration.
The deployment is handling by helm version 3

It built to run into kubernetes cluster

## Point of view

### Application configuration

- It may contain post-deployment and pre-deployment script
- It have to contain the chart to use for the deployment and the credential if needed

Configuration could set in [configmap](https://kubernetes.io/docs/concepts/configuration/configmap) or [secret](https://kubernetes.io/docs/concepts/configuration/secret)

## Examples

- TODO