# cloudfauj
Deploy Apps to your cloud without managing infrastructure

TODO
- refactor http requests logic, take out common code from all methods
- refactor: replace Sprintf() with simple string concat where possible
- replace all print statements & linebreaks with proper logging (stdout/err)

- should app logs be from very beginning or only beginning of latest deployment?
- how do we show app logs if multiple containers start as part of a deployment? 
- do we want to rename deployment status to deployment info (or just deployment)
- read up on websocket & gorilla lib

- make a decision on json data being interface{} or want to put in struct?
- Handle when required config files are not supplied/present (or required fields in them not present)
