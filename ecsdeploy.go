package ecsdeploy

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func Run(cluster string, service string, taskFamily string, container string, image string, keepService bool) error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := ecs.New(sess)

	taskDefinition, err := DescribeLatestTaskDefinition(svc, taskFamily)

	if err != nil {
		return err
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinition,
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

	if keepService == false {
		resp, err := UpdateService(svc, *newTask.TaskDefinition, cluster, service)

		if err != nil {
			return err
		} else {
			log.Println("Updated service: ", *resp.Service.TaskDefinition)
		}
	}

	return nil
}

func DescribeLatestTaskDefinition(svc *ecs.ECS, taskFamily string) (*string, error) {
	ListTaskDefinitionsInput := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(taskFamily),
		MaxResults:   aws.Int64(1),
		Sort:         aws.String("DESC"),
		Status:       aws.String("ACTIVE"),
	}

	listTaskDefinitionsOutput, err := svc.ListTaskDefinitions(ListTaskDefinitionsInput)

	if err != nil {
		return nil, err
	} else if len(listTaskDefinitionsOutput.TaskDefinitionArns) == 0 {
		return nil, errors.New("active task definition does not found")
	}

	splitArn := strings.Split(*listTaskDefinitionsOutput.TaskDefinitionArns[0], "/")
	return aws.String(splitArn[len(splitArn)-1]), nil
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
		Volumes:              task.Volumes,
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
