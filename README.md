# cloudfauj
Deploy Apps to your cloud without managing infrastructure

TODO
- Set Usage command for all commands that use args
- In deploy, how do we tell cli which artifact (docker image) to deploy
- Handle when required config files are not supplied/present

Find idomatic ways of:
- How to display an error on stderr and exit instantly with failure (set RunE instead of Run in Command)
- how to create api client at top level and pass on to child commands 
- How to print to stdout
- Read diff configs for diff actions (server, project, env)

Questions:
- How do we specify the container for deploy command
- Does dev explicitly have to create a project? Should we just register the project & apps in state when they're deploying for the first time?
- Is manual management of environments good approach for our use case? Compared to per-PR URLs that people get?
