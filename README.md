# cloudfauj
Automated Infrastructure provisioning in your own cloud

TODO

- write infra methods for env creation/deletion
- write infra methods for app creation
- implement all other side methods
- cli server: setup local dir for all storage
- revisit: we may not need the whole app logs command & its backend right now, we can just ship logs to cloudwatch
- polish the whole user experience
- write documentation
 launch v1

- read up on websocket & gorilla lib

- should app logs be from very beginning or only beginning of latest deployment?
- how do we show app logs if multiple containers start as part of a deployment? 
- do we want to rename deployment status to deployment info (or just deployment)

- replace all print statements & linebreaks with proper logging (stdout/err)
- make a decision on json data being interface{} or want to put in struct?
- Handle when required config files are not supplied/present (or required fields in them not present)
