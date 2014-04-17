package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
  "github.com/mgutz/ansi"
	"os"
	"strconv"
	"strings"
)

func main() {
	var listFlag bool
	var publicIp bool
	var index int
	flag.BoolVar(&listFlag, "list", false, "list instances")
	flag.BoolVar(&listFlag, "l", false, "list instances")
	flag.BoolVar(&publicIp, "public", false, "show public IP")
	flag.BoolVar(&publicIp, "p", false, "show public IP")
	flag.IntVar(&index, "n", 1, "nth instance")
	flag.Parse()

	if listFlag {
		list()
		os.Exit(0)
	}

	name := flag.Arg(0)

	if strings.Contains(name, "#") {
		parts := strings.Split(name, "#")
		name = parts[0]
		index, _ = strconv.Atoi(parts[1])
	}

	instances, err := getInstances(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	if len(instances) == 0 || index > len(instances) || index < 1 {
		fmt.Fprintf(os.Stderr, "Server Not Found\n")
		os.Exit(1)
	}
	instance := instances[index-1]
	if publicIp {
		fmt.Printf("%s\n", instance.PublicIpAddress)
	} else {
		fmt.Printf("%s\n", instance.PrivateIpAddress)
	}

}

func getInstances(name string) (instances []ec2.Instance, err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	e := ec2.New(auth, aws.USEast)
	filter := ec2.NewFilter()
	filter.Add("instance-state-name", "running")
	filter.Add("instance-state-name", "stopped")
	filter.Add("tag:Name", name)
	resp, err2 := e.Instances(nil, filter)
	instances = make([]ec2.Instance, 0, 5)
	for _, reservation := range resp.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, err2
}

func colorStatus(word string) (string) {
  if word == "running" {
    return ansi.Color("running", "green+b")
  } else if word == "stopped" {
    return ansi.Color("stopped", "red+b")
}
return word
}

func list() {
	instances, err := getInstances("*")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	for _, instance := range instances {
		tags := make(map[string]string)
		for _, tag := range instance.Tags {
			tags[tag.Key] = tag.Value
		}
		fmt.Printf("%-40s %-10s %-10s %-20s %-20s\n", tags["Name"], instance.InstanceId, colorStatus(instance.State.Name), instance.PrivateIpAddress, instance.PublicIpAddress)
	}
}
