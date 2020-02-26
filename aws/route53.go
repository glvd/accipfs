package aws

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/goextension/tool"
)

const (
	// HostedZoneID ...
	HostedZoneID = "ZULHCRKQ5YDLC"
	// RecordName ...
	RecordName = "gate.dhash.app"
)

// awsRoute ...
type awsRoute struct {
}

// Router ...
type Router interface {
	GetRecordSets() ([]*route53.ResourceRecordSet, error)
	ChangeSets([]*route53.ResourceRecordSet, string) (*route53.ChangeResourceRecordSetsOutput, error)
	BuildMultiValueRecordSets([]string) []*route53.ResourceRecordSet
	FilterFailedRecords([]*route53.ResourceRecordSet) []*route53.ResourceRecordSet
}

// NewRoute53 ...
func NewRoute53() Router {
	return &awsRoute{}
}

// GetRecordSets ...
func (r awsRoute) GetRecordSets() ([]*route53.ResourceRecordSet, error) {
	client := route53.New(session.New())
	recordSetsInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(HostedZoneID),
		StartRecordName: aws.String(RecordName),
	}
	sets, err := client.ListResourceRecordSets(recordSetsInput)
	if err != nil {
		fmt.Println("[get record sets failed]", err.Error())
		return []*route53.ResourceRecordSet{}, err
	}
	return sets.ResourceRecordSets, nil
}

// ChangeSets options: 'CREATE', 'DELETE', 'UPSERT'
func (r awsRoute) ChangeSets(sets []*route53.ResourceRecordSet, option string) (*route53.ChangeResourceRecordSetsOutput, error) {
	var changes []*route53.Change
	client := route53.New(session.New())

	for _, set := range sets {
		change := &route53.Change{
			Action:            aws.String(option),
			ResourceRecordSet: set,
		}
		changes = append(changes, change)
	}

	changeInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
			Comment: aws.String("gateway"),
		},
		HostedZoneId: aws.String(HostedZoneID),
	}
	res, err := client.ChangeResourceRecordSets(changeInput)
	return res, err
}

// BuildMultiValueRecordSets ...
func (r awsRoute) BuildMultiValueRecordSets(ips []string) []*route53.ResourceRecordSet {
	var sets []*route53.ResourceRecordSet
	for _, ip := range ips {
		awsRecordSet := &route53.ResourceRecordSet{
			Name:             aws.String(RecordName),
			Type:             aws.String("A"),
			SetIdentifier:    aws.String(tool.GenerateRandomString(5)),
			MultiValueAnswer: aws.Bool(true),
			ResourceRecords: []*route53.ResourceRecord{
				&route53.ResourceRecord{Value: aws.String(ip)},
			},
			TTL: aws.Int64(60),
		}
		sets = append(sets, awsRecordSet)
	}
	return sets
}

// FilterFailedRecords ...
func (r awsRoute) FilterFailedRecords(sets []*route53.ResourceRecordSet) []*route53.ResourceRecordSet {
	var failedSets []*route53.ResourceRecordSet
	for _, set := range sets {
		ip := set.ResourceRecords[0].Value
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		cmd := exec.CommandContext(ctx, "nc", "-vz", *ip, "8545")
		cmd.Start()
		if err := cmd.Wait(); err != nil {
			failedSets = append(failedSets, set)
		}
		cancel()
	}
	return failedSets
}
