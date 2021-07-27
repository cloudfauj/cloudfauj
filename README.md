# cloudfauj
Automated Infrastructure provisioning in your own cloud

TODO

- write infra methods for env creation/deletion - 31 july
- write infra methods for app creation - 6 aug
- implement all other side methods - 9 aug
- cli server: setup local dir for all storage - 13 aug
- polish the whole user experience - 19 aug
- write documentation - 25 aug
- launch v1 - 25 aug

- read up on websocket & gorilla lib

- should app logs be from very beginning or only beginning of latest deployment?
- how do we show app logs if multiple containers start as part of a deployment? 
- do we want to rename deployment status to deployment info (or just deployment)

- replace all print statements & linebreaks with proper logging (stdout/err)
- make a decision on json data being interface{} or want to put in struct?
- Handle when required config files are not supplied/present (or required fields in them not present)
