# CloudFauj
Deploy containers to your AWS without managing infrastructure

CloudFauj is a self-serve platform for developers to deploy containers without having to provision and manage the infrastructure for them.

It is self-hosted, creates resources in your AWS Account and aims to provide your DevOps team full flexibility and control over these resources.  

## For Devs
Developers can focus on building their applications. Once done, they add a `.cloudfauj.yml` to their project dir to declare the resources an app needs to run.

They then use Cloudfauj, either as a CLI tool or as part of their CI/CD pipeline, to deploy the app to an environment. Cloudfauj takes care of provisioning the infrastructure to run the app in the cloud. 

## For Ops
Ops teams use Cloudfauj to create and manage environments in their own AWS account.

Cloudfauj automates creating all resources to run apps in different environments. This removes toil for Ops, while still giving them an extremely high degree of control over the infrastructure & costs.

TODO

- write documentation - 25 aug
  - cli help text
  - GH readme with index
  - guide: deploy & teardown a node app using cloudfauj
  - cli: version
  - GH release binary
- polish the whole user experience
  - TEST: before modifying any state or infra, check if aws creds are supplied
- launch v1 - 25 aug

Issues:
1. During app destroy, cli times out before our ecs timeout. Also, ecs drain timeout may not be enough, so I had to re-run app destroy.
```
./cloudfauj app destroy nginx-api --env raghavdemo
   Destroying nginx-api from raghavdemo
   Error: Delete "http://127.0.0.1:6200/v1/app/nginx-api?env=raghavdemo": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

post v1:
- test whether server can handle multiple requests simultaneously
- infra management via terraform
- app destroy: stream destroy logs to client
- automatic tls + load balancer
- route53 dns integration so we can provide url for app
- destroy all apps as prerequisite when destroying an env
- add `--dry` to server, which doesn't check for aws creds, doesn't provision any real infra and manages state in-mem



- read up on websocket & gorilla lib
- do we want to rename deployment status to deployment info (or just deployment)
- replace all print statements & linebreaks with proper logging (stdout/err)
- Handle when required config files are not supplied/present (or required fields in them not present)
