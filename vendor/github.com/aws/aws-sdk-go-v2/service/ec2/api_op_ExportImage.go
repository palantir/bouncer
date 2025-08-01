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

// Exports an Amazon Machine Image (AMI) to a VM file. For more information, see [Exporting a VM directly from an Amazon Machine Image (AMI)]
// in the VM Import/Export User Guide.
//
// [Exporting a VM directly from an Amazon Machine Image (AMI)]: https://docs.aws.amazon.com/vm-import/latest/userguide/vmexport_image.html
func (c *Client) ExportImage(ctx context.Context, params *ExportImageInput, optFns ...func(*Options)) (*ExportImageOutput, error) {
	if params == nil {
		params = &ExportImageInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ExportImage", params, optFns, c.addOperationExportImageMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ExportImageOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ExportImageInput struct {

	// The disk image format.
	//
	// This member is required.
	DiskImageFormat types.DiskImageFormat

	// The ID of the image.
	//
	// This member is required.
	ImageId *string

	// The Amazon S3 bucket for the destination image. The destination bucket must
	// exist.
	//
	// This member is required.
	S3ExportLocation *types.ExportTaskS3LocationRequest

	// Token to enable idempotency for export image requests.
	ClientToken *string

	// A description of the image being exported. The maximum length is 255 characters.
	Description *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// The name of the role that grants VM Import/Export permission to export images
	// to your Amazon S3 bucket. If this parameter is not specified, the default role
	// is named 'vmimport'.
	RoleName *string

	// The tags to apply to the export image task during creation.
	TagSpecifications []types.TagSpecification

	noSmithyDocumentSerde
}

type ExportImageOutput struct {

	// A description of the image being exported.
	Description *string

	// The disk image format for the exported image.
	DiskImageFormat types.DiskImageFormat

	// The ID of the export image task.
	ExportImageTaskId *string

	// The ID of the image.
	ImageId *string

	// The percent complete of the export image task.
	Progress *string

	// The name of the role that grants VM Import/Export permission to export images
	// to your Amazon S3 bucket.
	RoleName *string

	// Information about the destination Amazon S3 bucket.
	S3ExportLocation *types.ExportTaskS3Location

	// The status of the export image task. The possible values are active , completed
	// , deleting , and deleted .
	Status *string

	// The status message for the export image task.
	StatusMessage *string

	// Any tags assigned to the export image task.
	Tags []types.Tag

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationExportImageMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpExportImage{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpExportImage{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ExportImage"); err != nil {
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
	if err = addIdempotencyToken_opExportImageMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpExportImageValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opExportImage(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpExportImage struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpExportImage) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpExportImage) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*ExportImageInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *ExportImageInput ")
	}

	if input.ClientToken == nil {
		t, err := m.tokenProvider.GetIdempotencyToken()
		if err != nil {
			return out, metadata, err
		}
		input.ClientToken = &t
	}
	return next.HandleInitialize(ctx, in)
}
func addIdempotencyToken_opExportImageMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpExportImage{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opExportImage(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ExportImage",
	}
}
