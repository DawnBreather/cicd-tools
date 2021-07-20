package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	types2 "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/thoas/go-funk"
)

type ASGs []*ASG
func (a *ASGs) CollectAllFromRegion(region string) *ASGs{
	client := autoscaling.NewFromConfig(newConfig(region))

	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	resp, err := client.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		_logger.Errorf("Unable to list ASGs in region { %s }: %v", region, err)
	}
	for _, g := range resp.AutoScalingGroups{
		var asg = ASG{}
		asg.
			SetName(*g.AutoScalingGroupName).
			SetRegion(region).
			SetAwsObject(g).
			CollectInstances()
		a.append(&asg)
	}

	return a
}

func (a *ASGs) FilterAsgsByRegion(region string) *ASGs{
	gs := funk.Filter(*a, func(g *ASG) bool {
		return g.region == region
	}).(ASGs)

	return &gs
}

func (a *ASGs) FilterAsgByNameAndRegion(name, region string) *ASG{
	gs := funk.Filter(*a, func(g *ASG) bool {
		return g.region == region && g.name == name
	}).(ASGs)

	if len(gs) > 0 {
		return gs[0]
	}

	return nil
}

func (a *ASGs) append(asg *ASG) *ASGs{
	*a = append(*a, asg)
	return a
}

type ASG struct{
	name string
	region string
	instances []*EC2
	awsObject types2.AutoScalingGroup
}

func (a *ASG) SetAwsObject(awsObject types2.AutoScalingGroup) *ASG{
	a.awsObject = awsObject
	return a
}

func (a *ASG) SetName(name string) *ASG {
	a.name = name
	return a
}

func (a *ASG) SetRegion(region string) *ASG {
	a.region = region
	return a
}

func (a *ASG) GetAwsObject() types2.AutoScalingGroup {
	return a.awsObject
}
func (a *ASG) GetName() string {
	return a.name
}
func (a *ASG) GetRegion() string {
	return a.region
}
func (a *ASG) GetInstances() []*EC2 {
	return a.instances
}

func (a *ASG) ReadAwsObject() *ASG {
	client := autoscaling.NewFromConfig(newConfig(a.region))

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			a.name,
		},
	}

	resp, err := client.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		_logger.Errorf("Unable to describe ASG { %s } in region { %s }: %v", a.name, a.region, err)
		return a
	}

	asgCount := len(resp.AutoScalingGroups)
	if asgCount == 0 {
		_logger.Warnf("Not found ASG { %s } in region { %s }", a.name, a.region)
		return a
	} else if asgCount > 1 {
		_logger.Warnf("Identified multiple ASGs with single name { %s } in region { %s }", a.name, a.region)
	}
	for _, g := range resp.AutoScalingGroups{
		a.awsObject = g
	}

	return a
}

func (a *ASG) CollectInstances() *ASG{

	var instances []*EC2

	for _, i := range a.awsObject.Instances{
		e := EC2{}
		e.SetId(*i.InstanceId).
			SetRegion(a.region).
			ReadAwsObject()
		instances = append(instances, &e)
	}

	a.instances = instances

	return a
}

func (a *ASG) GetTagValue(name string) string {
	for _, t := range a.awsObject.Tags{
		tName := *t.Key
		tValue := *t.Value
		if tName == name {
			return tValue
		}
	}

	return ""
}