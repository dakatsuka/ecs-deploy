package ecsdeploy

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"log"
)

func Run(cluster string, service string, container string, image string) error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := ecs.New(sess)

	ecsService, err := DescribeActiveService(svc, cluster, service)

	if err != nil {
		return err
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*ecsService.TaskDefinition),
	}

	task, err := svc.DescribeTaskDefinition(params)

	if err != nil {
		return err
	} else {
		log.Println("Got current TaskDefinition: ", *task.TaskDefinition.TaskDefinitionArn)
	}

	newTask, err := UpdateImage(svc, *task.TaskDefinition, container, image)

	if err != nil {
		return err
	} else {
		log.Println("Registered TaskDefinition with new image: ", *newTask.TaskDefinition.TaskDefinitionArn)
	}

	resp, err := UpdateService(svc, *newTask.TaskDefinition, cluster, service)

	if err != nil {
		return err
	} else {
		log.Println("Updated service: ", *resp.Service.TaskDefinition)
	}

	return nil
}

func DescribeActiveService(svc *ecs.ECS, cluster string, service string) (*ecs.Service, error) {
	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(service),
		},
		Cluster: aws.String(cluster),
	}

	resp, err := svc.DescribeServices(params)

	if err != nil {
		return nil, err
	}

	for _, v := range resp.Services {
		if *v.ServiceName == service && *v.Status == "ACTIVE" {
			return v, nil
		}
	}

	return nil, errors.New("active service does not found")
}

func UpdateImage(svc *ecs.ECS, task ecs.TaskDefinition, container string, image string) (*ecs.RegisterTaskDefinitionOutput, error) {
	var newContainerDefinitions []*ecs.ContainerDefinition

	for _, v := range task.ContainerDefinitions {
		if *v.Name == container {
			v.Image = aws.String(image)
		}
		newContainerDefinitions = append(newContainerDefinitions, v)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: newContainerDefinitions,
		Family:               task.Family,
		TaskRoleArn:          task.TaskRoleArn,
	}

	return svc.RegisterTaskDefinition(params)
}

func UpdateService(svc *ecs.ECS, task ecs.TaskDefinition, cluster string, service string) (*ecs.UpdateServiceOutput, error) {
	params := &ecs.UpdateServiceInput{
		Service:        aws.String(service),
		Cluster:        aws.String(cluster),
		TaskDefinition: task.TaskDefinitionArn,
	}

	return svc.UpdateService(params)
}
