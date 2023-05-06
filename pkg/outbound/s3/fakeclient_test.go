// These tests does not verify whether or not our fake client correctly mimics the real client.
// Rather, they verify that the fake client is, at least, consistent; i.e., the set of methods
// work together and produce what we assume the real client would produce.
// They serve as a guideline to measure our assumptions about the real client, such that we can
// find gaps in our understanding, and hence, help us fix our assumptions about what a fake is.
package s3

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3Operations(t *testing.T) {
	client := newMockedS3Client()

	// Test PutObject operation
	bucketName := "test-bucket"
	objectKey := "test-object"
	objectContent := "test-content"
	_, err := client.PutObject(&s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   bytes.NewReader([]byte(objectContent)),
	})
	assert.NoError(t, err)

	// Test GetObject operation
	getObjectOutput, err := client.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(getObjectOutput.Body)
	assert.Equal(t, objectContent, buf.String())

	// Test HeadBucket operation
	_, err = client.HeadBucket(&s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	assert.NoError(t, err)

	// Test DeleteObject operation
	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	assert.NoError(t, err)

	// Test that the object no longer exists
	_, err = client.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	assert.Error(t, err)
}
