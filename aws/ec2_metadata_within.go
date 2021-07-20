package aws

import (
	"github.com/go-resty/resty/v2"
)

type Ec2MetadataWithin struct {
	magicUrl             string
	instanceIdPath       string
	availabilityZonePath string

	instanceId       string
	availabilityZone string
	region           string
}

func (e Ec2MetadataWithin) GetInstanceId() string{
	return e.instanceId
}
func (e Ec2MetadataWithin) GetAvailabilityZone() string{
	return e.availabilityZone
}
func (e Ec2MetadataWithin) GetRegion() string{
	return e.region
}

func (e *Ec2MetadataWithin) Init() *Ec2MetadataWithin{
	e.magicUrl = "http://169.254.169.254"
	e.instanceIdPath = "/latest/meta-data/instance-id"
	e.availabilityZonePath = "/latest/meta-data/placement/availability-zone"

	return e
}

func (e *Ec2MetadataWithin) read(url, errorMessage string) []byte{
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Get(url)

	if err != nil {
		_logger.Errorf("%s: %v", errorMessage, err)
		return []byte{}
	}

	return resp.Body()
}

func (e *Ec2MetadataWithin) ReadInstanceId() *Ec2MetadataWithin{
	body := e.read(e.magicUrl + e.instanceIdPath, "Unable to read Instance ID")
	if len(body) != 0 {
		e.instanceId = string(body)
		_logger.Info("Instance ID identified: %s", e.instanceId)
	}

	return e
}

func (e *Ec2MetadataWithin) ReadRegionAndAvailabilityZone() *Ec2MetadataWithin{
	body := e.read(e.magicUrl + e.instanceIdPath, "Unable to read Availability Zone")
	if len(body) != 0 {
		e.availabilityZone = string(body)
		e.region = e.availabilityZone[:len(e.availabilityZone) - 1]
		_logger.Info("Availability Zone identified: %s", e.availabilityZone)
		_logger.Info("Region identified: %s", e.region)
	}

	return e
}

