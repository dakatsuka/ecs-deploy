package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var version string

func main() {
	var (
		clusterName   = kingpin.Flag("cluster", "Set cluster name").Required().String()
		serviceName   = kingpin.Flag("service", "Set service name").Required().String()
		containerName = kingpin.Flag("container", "Set container name").Required().String()
		image         = kingpin.Flag("image", "Set image").Required().String()
	)
	kingpin.Version(version)
	kingpin.Parse()

	sess, err := session.NewSession()

	if err != nil {
		fmt.Println("Failed to create session", err)
		os.Exit(1)
	}

	svc := ecs.New(sess)

	service, err := describeActiveService(svc, clusterName, serviceName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*service.TaskDefinition),
	}

	task, err := svc.DescribeTaskDefinition(params)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		log.Println("Got current TaskDefinition: ", *task.TaskDefinition.TaskDefinitionArn)
	}

	newTask, err := updateImage(svc, *task.TaskDefinition, containerName, image)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		log.Println("Registered TaskDefinition with new image: ", *newTask.TaskDefinition.TaskDefinitionArn)
	}

	resp, err := updateService(svc, *newTask.TaskDefinition, clusterName, serviceName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	log.Println("Updated service: ", *resp.Service.TaskDefinition)
}

func describeActiveService(svc *ecs.ECS, cluster *string, service *string) (*ecs.Service, error) {
	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(*service),
		},
		Cluster: aws.String(*cluster),
	}

	resp, err := svc.DescribeServices(params)

	if err != nil {
		return nil, err
	}

	for _, v := range resp.Services {
		if *v.ServiceName == *service && *v.Status == "ACTIVE" {
			return v, nil
		}
	}

	return nil, errors.New("active service does not found")
}

func updateImage(svc *ecs.ECS, task ecs.TaskDefinition, container *string, image *string) (*ecs.RegisterTaskDefinitionOutput, error) {
	var newContainerDefinitions []*ecs.ContainerDefinition

	for _, v := range task.ContainerDefinitions {
		if *v.Name == *container {
			v.Image = image
		}
		newContainerDefinitions = append(newContainerDefinitions, v)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: newContainerDefinitions,
		Family:               aws.String(*task.Family),
		TaskRoleArn:          aws.String(*task.TaskRoleArn),
	}

	return svc.RegisterTaskDefinition(params)
}

func updateService(svc *ecs.ECS, task ecs.TaskDefinition, cluster *string, service *string) (*ecs.UpdateServiceOutput, error) {
	params := &ecs.UpdateServiceInput{
		Service:        aws.String(*service),
		Cluster:        aws.String(*cluster),
		TaskDefinition: aws.String(*task.TaskDefinitionArn),
	}

	return svc.UpdateService(params)
}
