# CloudFauj
**Deploy containers to your AWS without managing infrastructure**

**NOTE**: Cloudfauj is under active development. We do not recommend it for production environments or if you're not comfortable with AWS. There may be breaking changes in future releases.

CloudFauj is a self-serve platform for developers to deploy containers without having to provision and manage the infrastructure for them.

It is self-hosted, creates resources in your AWS Account and aims to provide your DevOps team full flexibility and control over these resources.  

## For Devs
Developers can focus on building their applications. Once done, they add a `.cloudfauj.yml` to their project dir to declare the resources an app needs to run.

They then use Cloudfauj, either as a CLI tool or as part of their CI/CD pipeline, to deploy the app to an environment. Cloudfauj takes care of provisioning the infrastructure to run the app in the cloud. 

## For Ops
Ops teams use Cloudfauj to create and manage environments in their own AWS account.

Cloudfauj automates creating all resources to run apps in different environments. This removes toil for Ops, while still giving them an extremely high degree of control over the infrastructure & costs.

## Documentation
1. [Getting Started](./docs/getting-started.md)
2. [Concepts](./docs/concepts.md)
    1. [Architecture](./docs/concepts.md#architecture)
    2. [Server](./docs/concepts.md#server)
    3. [Client](./docs/concepts.md#client)
    4. [Environment](./docs/concepts.md#environment)
    5. [Application](./docs/concepts.md#application)
3. [Creating an Environment](./docs/create-env.md)
4. [Deploying an Application](./docs/deploy-app.md)

## License
MPL-2.0

Please see the [License](./LICENSE) for details.
