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

// Modify the auto-placement setting of a Dedicated Host. When auto-placement is
// enabled, any instances that you launch with a tenancy of host but without a
// specific host ID are placed onto any available Dedicated Host in your account
// that has auto-placement enabled. When auto-placement is disabled, you need to
// provide a host ID to have the instance launch onto a specific host. If no host
// ID is provided, the instance is launched onto a suitable host with
// auto-placement enabled.
//
// You can also use this API action to modify a Dedicated Host to support either
// multiple instance types in an instance family, or to support a specific instance
// type only.
func (c *Client) ModifyHosts(ctx context.Context, params *ModifyHostsInput, optFns ...func(*Options)) (*ModifyHostsOutput, error) {
	if params == nil {
		params = &ModifyHostsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ModifyHosts", params, optFns, c.addOperationModifyHostsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ModifyHostsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ModifyHostsInput struct {

	// The IDs of the Dedicated Hosts to modify.
	//
	// This member is required.
	HostIds []string

	// Specify whether to enable or disable auto-placement.
	AutoPlacement types.AutoPlacement

	// Indicates whether to enable or disable host maintenance for the Dedicated Host.
	// For more information, see [Host maintenance]in the Amazon EC2 User Guide.
	//
	// [Host maintenance]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/dedicated-hosts-maintenance.html
	HostMaintenance types.HostMaintenance

	// Indicates whether to enable or disable host recovery for the Dedicated Host.
	// For more information, see [Host recovery]in the Amazon EC2 User Guide.
	//
	// [Host recovery]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/dedicated-hosts-recovery.html
	HostRecovery types.HostRecovery

	// Specifies the instance family to be supported by the Dedicated Host. Specify
	// this parameter to modify a Dedicated Host to support multiple instance types
	// within its current instance family.
	//
	// If you want to modify a Dedicated Host to support a specific instance type
	// only, omit this parameter and specify InstanceType instead. You cannot specify
	// InstanceFamily and InstanceType in the same request.
	InstanceFamily *string

	// Specifies the instance type to be supported by the Dedicated Host. Specify this
	// parameter to modify a Dedicated Host to support only a specific instance type.
	//
	// If you want to modify a Dedicated Host to support multiple instance types in
	// its current instance family, omit this parameter and specify InstanceFamily
	// instead. You cannot specify InstanceType and InstanceFamily in the same request.
	InstanceType *string

	noSmithyDocumentSerde
}

type ModifyHostsOutput struct {

	// The IDs of the Dedicated Hosts that were successfully modified.
	Successful []string

	// The IDs of the Dedicated Hosts that could not be modified. Check whether the
	// setting you requested can be used.
	Unsuccessful []types.UnsuccessfulItem

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationModifyHostsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpModifyHosts{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpModifyHosts{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ModifyHosts"); err != nil {
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
	if err = addOpModifyHostsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opModifyHosts(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opModifyHosts(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ModifyHosts",
	}
}
