package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/thoas/go-funk"
)

type EC2s []*EC2
func (e *EC2s) CollectAllFromRegion(region string) *EC2s{
	client := ec2.NewFromConfig(newConfig(region))

	input := &ec2.DescribeInstancesInput{}
	resp, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		_logger.Errorf("Unable to list EC2 instances in region { %s }: %v", region, err)
	}
	for _, r := range resp.Reservations{
		for _, i := range r.Instances {
			var ec2 = EC2{}
			ec2.
				SetId(*i.InstanceId).
				SetRegion(region).
				SetAwsObject(i)
		}
	}

	return e
}

func (e *EC2s) FilterEc2sByRegion(region string) *EC2s{
	ec2s := funk.Filter(*e, func(i *EC2) bool {
		return i.region == region
	}).(EC2s)

	return &ec2s
}

func (e *EC2s) GetEc2ByInternalIp(internalIp string) *EC2{
	ec2s := funk.Filter(*e, func(i *EC2) bool {
		for _, ni := range i.GetAwsObject().NetworkInterfaces{
			if *ni.PrivateIpAddress == internalIp{
				return true
			}
		}
		return false
	}).(EC2s)

	if len(ec2s) > 0 {
		return ec2s[0]
	} else {
		return nil
	}
}


func (e *EC2s) append(ec2 *EC2) *EC2s{
	*e = append(*e, ec2)
	return e
}



type EC2 struct {
	id string
	region string
	awsObject types.Instance
	relatedAsg *ASG
}

func (e *EC2) SetId(id string) *EC2{
	e.id = id
	return e
}

func (e *EC2) SetRegion(region string) *EC2{
	e.region = region
	return e
}

func (e *EC2) SetAwsObject(awsEc2Object types.Instance) *EC2{
	e.awsObject = awsEc2Object
	return e
}

func (e *EC2) GetId() string{
	return e.id
}
func (e *EC2) GetRegion() string{
	return e.region
}
func (e *EC2) GetAwsObject() types.Instance{
	return e.awsObject
}

func (e *EC2) ReadAwsObject() *EC2{
	client := ec2.NewFromConfig(newConfig(e.region))
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			e.id,
		},
	}

	resp, err := client.DescribeInstances(context.TODO(), input)

	if err != nil {
		_logger.Errorf("Unable to describe instance { %s }: %v", e.id, err)
		return nil
	}

	if len(resp.Reservations) > 0{
		if len(resp.Reservations[0].Instances) > 0 {
			e.awsObject = resp.Reservations[0].Instances[0]
		}
	}

	return nil
}

func (e *EC2) GetTagValue(name string) string {
	for _, t := range e.awsObject.Tags{
		tName := *t.Key
		tValue := *t.Value
		if tName == name {
			return tValue
		}
	}

	return ""
}