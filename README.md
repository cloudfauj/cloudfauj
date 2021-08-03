# cloudfauj
Automated Infrastructure provisioning in your own cloud

TODO

- polish the whole user experience
  - review all commands, their outputs
  - tail ecs deployment logs - 19 aug
- write documentation - 25 aug
- launch v1 - 25 aug

post v1:
- infra management via terraform
- automatic tls + load balancer
- route53 dns integration so we can provide url for app
- destroy all apps as prerequisite when destroying an env



- read up on websocket & gorilla lib
- do we want to rename deployment status to deployment info (or just deployment)
- replace all print statements & linebreaks with proper logging (stdout/err)
- Handle when required config files are not supplied/present (or required fields in them not present)
