package nacl

import (
	"waze/terraform/aws_terraforming/aws_generator"
	"waze/terraform/terraform_utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ignoreKey = map[string]bool{
	"id": true,
}

var allowEmptyValues = map[string]bool{
	"tags.": true,
}

type NaclGenerator struct {
	aws_generator.BasicGenerator
}

func (NaclGenerator) createResources(nacls *ec2.DescribeNetworkAclsOutput) []terraform_utils.TerraformResource {
	resoures := []terraform_utils.TerraformResource{}
	for _, nacl := range nacls.NetworkAcls {
		resourceName := ""
		if len(nacl.Tags) > 0 {
			for _, tag := range nacl.Tags {
				if aws.StringValue(tag.Key) == "Name" {
					resourceName = aws.StringValue(tag.Value)
					break
				}
			}
		}
		resoures = append(resoures, terraform_utils.TerraformResource{
			ResourceType: "aws_network_acl",
			ResourceName: resourceName,
			Item:         nil,
			ID:           aws.StringValue(nacl.NetworkAclId),
			Provider:     "aws",
		})
	}
	return resoures
}

func (g NaclGenerator) Generate(region string) ([]terraform_utils.TerraformResource, map[string]terraform_utils.ResourceMetaData, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(region)})
	svc := ec2.New(sess)
	nacls, err := svc.DescribeNetworkAcls(&ec2.DescribeNetworkAclsInput{})
	if err != nil {
		return []terraform_utils.TerraformResource{}, map[string]terraform_utils.ResourceMetaData{}, err

	}
	resources := g.createResources(nacls)
	metadata := terraform_utils.NewResourcesMetaData(resources, ignoreKey, allowEmptyValues, map[string]string{})
	return resources, metadata, nil

}