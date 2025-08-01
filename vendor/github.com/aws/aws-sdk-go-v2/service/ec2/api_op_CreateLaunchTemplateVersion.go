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

// Creates a new version of a launch template. You must specify an existing launch
// template, either by name or ID. You can determine whether the new version
// inherits parameters from a source version, and add or overwrite parameters as
// needed.
//
// Launch template versions are numbered in the order in which they are created.
// You can't specify, change, or replace the numbering of launch template versions.
//
// Launch templates are immutable; after you create a launch template, you can't
// modify it. Instead, you can create a new version of the launch template that
// includes the changes that you require.
//
// For more information, see [Modify a launch template (manage launch template versions)] in the Amazon EC2 User Guide.
//
// [Modify a launch template (manage launch template versions)]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/manage-launch-template-versions.html
func (c *Client) CreateLaunchTemplateVersion(ctx context.Context, params *CreateLaunchTemplateVersionInput, optFns ...func(*Options)) (*CreateLaunchTemplateVersionOutput, error) {
	if params == nil {
		params = &CreateLaunchTemplateVersionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateLaunchTemplateVersion", params, optFns, c.addOperationCreateLaunchTemplateVersionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateLaunchTemplateVersionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateLaunchTemplateVersionInput struct {

	// The information for the launch template.
	//
	// This member is required.
	LaunchTemplateData *types.RequestLaunchTemplateData

	// Unique, case-sensitive identifier you provide to ensure the idempotency of the
	// request. If a client token isn't specified, a randomly generated token is used
	// in the request to ensure idempotency.
	//
	// For more information, see [Ensuring idempotency].
	//
	// Constraint: Maximum 128 ASCII characters.
	//
	// [Ensuring idempotency]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html
	ClientToken *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// The ID of the launch template.
	//
	// You must specify either the launch template ID or the launch template name, but
	// not both.
	LaunchTemplateId *string

	// The name of the launch template.
	//
	// You must specify either the launch template ID or the launch template name, but
	// not both.
	LaunchTemplateName *string

	// If true , and if a Systems Manager parameter is specified for ImageId , the AMI
	// ID is displayed in the response for imageID . For more information, see [Use a Systems Manager parameter instead of an AMI ID] in the
	// Amazon EC2 User Guide.
	//
	// Default: false
	//
	// [Use a Systems Manager parameter instead of an AMI ID]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/create-launch-template.html#use-an-ssm-parameter-instead-of-an-ami-id
	ResolveAlias *bool

	// The version of the launch template on which to base the new version. Snapshots
	// applied to the block device mapping are ignored when creating a new version
	// unless they are explicitly included.
	//
	// If you specify this parameter, the new version inherits the launch parameters
	// from the source version. If you specify additional launch parameters for the new
	// version, they overwrite any corresponding launch parameters inherited from the
	// source version.
	//
	// If you omit this parameter, the new version contains only the launch parameters
	// that you specify for the new version.
	SourceVersion *string

	// A description for the version of the launch template.
	VersionDescription *string

	noSmithyDocumentSerde
}

type CreateLaunchTemplateVersionOutput struct {

	// Information about the launch template version.
	LaunchTemplateVersion *types.LaunchTemplateVersion

	// If the new version of the launch template contains parameters or parameter
	// combinations that are not valid, an error code and an error message are returned
	// for each issue that's found.
	Warning *types.ValidationWarning

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateLaunchTemplateVersionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateLaunchTemplateVersion{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateLaunchTemplateVersion{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateLaunchTemplateVersion"); err != nil {
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
	if err = addIdempotencyToken_opCreateLaunchTemplateVersionMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateLaunchTemplateVersionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateLaunchTemplateVersion(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpCreateLaunchTemplateVersion struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpCreateLaunchTemplateVersion) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpCreateLaunchTemplateVersion) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*CreateLaunchTemplateVersionInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *CreateLaunchTemplateVersionInput ")
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
func addIdempotencyToken_opCreateLaunchTemplateVersionMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpCreateLaunchTemplateVersion{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opCreateLaunchTemplateVersion(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateLaunchTemplateVersion",
	}
}
