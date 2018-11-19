package converters

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// TagsToMap converts a []*ec2.Tag into a map[string]string.
func TagsToMap(src []*ec2.Tag) map[string]string {
	// Create an array of exactly the length we require to hopefully avoid some
	// allocations while looping.
	tags := make(map[string]string)

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapToTags converts a map[string]string to a []*ec2.Tag
func MapToTags(src map[string]string) []*ec2.Tag {
	// Create an array of exactly the length we require to hopefully avoid some
	// allocations while looping.
	tags := make([]*ec2.Tag, 0, len(src))

	for k, v := range src {
		tag := &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}
