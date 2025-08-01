// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Deregisters the specified AMI. A deregistered AMI can't be used to launch new
// instances.
//
// If a deregistered EBS-backed AMI matches a Recycle Bin retention rule, it moves
// to the Recycle Bin for the specified retention period. It can be restored before
// its retention period expires, after which it is permanently deleted. If the
// deregistered AMI doesn't match a retention rule, it is permanently deleted
// immediately. For more information, see [Recover deleted Amazon EBS snapshots and EBS-backed AMIs with Recycle Bin]in the Amazon EBS User Guide.
//
// When deregistering an EBS-backed AMI, you can optionally delete its associated
// snapshots at the same time. However, if a snapshot is associated with multiple
// AMIs, it won't be deleted even if specified for deletion, although the AMI will
// still be deregistered.
//
// Deregistering an AMI does not delete the following:
//
//   - Instances already launched from the AMI. You'll continue to incur usage
//     costs for the instances until you terminate them.
//
//   - For EBS-backed AMIs: Snapshots that are associated with multiple AMIs.
//     You'll continue to incur snapshot storage costs.
//
//   - For instance store-backed AMIs: The files uploaded to Amazon S3 during AMI
//     creation. You'll continue to incur S3 storage costs.
//
// For more information, see [Deregister an Amazon EC2 AMI] in the Amazon EC2 User Guide.
//
// [Deregister an Amazon EC2 AMI]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/deregister-ami.html
// [Recover deleted Amazon EBS snapshots and EBS-backed AMIs with Recycle Bin]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/recycle-bin.html
func (c *Client) DeregisterImage(ctx context.Context, params *DeregisterImageInput, optFns ...func(*Options)) (*DeregisterImageOutput, error) {
	if params == nil {
		params = &DeregisterImageInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DeregisterImage", params, optFns, c.addOperationDeregisterImageMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DeregisterImageOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Contains the parameters for DeregisterImage.
type DeregisterImageInput struct {

	// The ID of the AMI.
	//
	// This member is required.
	ImageId *string

	// Specifies whether to delete the snapshots associated with the AMI during
	// deregistration.
	//
	// If a snapshot is associated with multiple AMIs, it is not deleted, regardless
	// of this setting.
	//
	// Default: The snapshots are not deleted.
	DeleteAssociatedSnapshots *bool

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	noSmithyDocumentSerde
}

type DeregisterImageOutput struct {

	// The deletion result for each snapshot associated with the AMI, including the
	// snapshot ID and its success or error code.
	DeleteSnapshotResults []types.DeleteSnapshotReturnCode

	// Returns true if the request succeeds; otherwise, it returns an error.
	Return *bool

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDeregisterImageMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpDeregisterImage{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDeregisterImage{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DeregisterImage"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = addOpDeregisterImageValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDeregisterImage(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addInterceptBeforeRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addInterceptAttempt(stack, options); err != nil {
		return err
	}
	if err = addInterceptExecution(stack, options); err != nil {
		return err
	}
	if err = addInterceptBeforeSerialization(stack, options); err != nil {
		return err
	}
	if err = addInterceptAfterSerialization(stack, options); err != nil {
		return err
	}
	if err = addInterceptBeforeSigning(stack, options); err != nil {
		return err
	}
	if err = addInterceptAfterSigning(stack, options); err != nil {
		return err
	}
	if err = addInterceptTransmit(stack, options); err != nil {
		return err
	}
	if err = addInterceptBeforeDeserialization(stack, options); err != nil {
		return err
	}
	if err = addInterceptAfterDeserialization(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opDeregisterImage(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DeregisterImage",
	}
}
