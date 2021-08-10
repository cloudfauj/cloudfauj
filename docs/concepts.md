## Architecture
Cloudfauj follows a client-server architecture.

The client is used by humans or deployment systems like Jenkins to make requests to the server. The Server does all the heavy-lifting to manage infrastructure.

### Server
Server is the first thing you launch after downloading Cloudfauj.

It is responsible for creating, managing and tracking resources in the cloud. It contains data & logs of all past and ongoing deployments. The server maintains all its data inside a single [data directory](./getting-started.md#configuration).

See [Getting Started](./getting-started.md#server) on how to launch a server.

### Client
A Cloudfauj client lets you send commands to the server to [create and manage environments](./create-env.md), [deploy applications](./deploy-app.md) and get information about them.

As of this writing, Cloudfauj provides a commandline client only. This can either be invoked manually by devs or ops or it can be integrated into a CI/CD pipeline such as jenkins.

## Environment
An Environment is a collection of cloud resources that are logically isolated from others. It contains the base infrastructure to run applications.

A Cloudfauj application is always deployed to an env. You can create multiple environments such as `staging`, `test1`, `hotfix-payment-integration`, etc. Each env is designed to contain your entire set of microservices.

See [Creating an Environment](./create-env.md).

## Application
An application is the business logic you write in some programming language.

By default, Cloudfauj looks for a `.cloudfauj.yml` file in the root directory of your app. It assumes that 1 version control repository only contains 1 app. Cloudfauj currently doesn't support monorepos.

An app must be packed into a Docker image, published to AWS ECR, then deployed to an environment using the [client](#client). Cloudfauj takes care of creating & managing all infrastructure to run the containers behind the scenes.

### Deployment
Every time you ask Cloudfauj to run a new artifact for an existing application, a Deployment is created.

A successful deployment in the target environment replaces the currently running docker image with the desired one. A failed deployment exits with a non-zero code. This help you integrate cloudfauj in your CI/CD systems.

Each deployment has a unique ID. See [Deploying an Application](./deploy-app.md)

**Previous**: [Getting Started](./getting-started.md)

**Next**: [Creating an Environment](./create-env.md)
