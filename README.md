# ecs-deploy

This is a tool to easily update the docker image of task definition on Amazon ECS/ECR.

## Usage

```
$ export AWS_REGION=
$ export AWS_ACCESS_KEY_ID=
$ export AWS_SECRET_ACCESS_KEY=

$ ecs-deploy \
  --cluster=<cluster-name> \
  --service=<service-name> \
  --task-family=<task-family>
  --container=<container-name> \
  --image=<new image> \
```

## Help

```
$ ecs-deploy --help
usage: ecs-deploy --cluster=CLUSTER --service=SERVICE --container=CONTAINER --image=IMAGE [<flags>]

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --cluster=CLUSTER          Set cluster name
  --service=SERVICE          Set service name
  --task-family=TASK-FAMILY  Set task definition family
  --container=CONTAINER      Set container name
  --image=IMAGE              Set image
  --keep-service             Not update service
  --version                  Show application version.
  --help                     Show context-sensitive help (also try --help-long and --help-man).
```
