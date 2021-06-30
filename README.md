# cloudfauj
Deploy Apps to your cloud without managing infrastructure

Find idomatic ways of:
- how to create api client at top level and pass on to child commands 
- How to print to stdout
- How to display an error on stderr and exit instantly with failure
- Read diff configs for diff actions (server, project, env)

Questions:
- How do we specify the container for deploy command
- Does dev explicitly have to create a project? Should we just register the project & apps in state when they're deploying for the first time?
- Is manual management of environments good approach for our use case? Compared to per-PR URLs that people get?
