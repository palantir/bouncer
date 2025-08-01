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

// Modifies the configuration of an existing route server.
//
// Amazon VPC Route Server simplifies routing for traffic between workloads that
// are deployed within a VPC and its internet gateways. With this feature, VPC
// Route Server dynamically updates VPC and internet gateway route tables with your
// preferred IPv4 or IPv6 routes to achieve routing fault tolerance for those
// workloads. This enables you to automatically reroute traffic within a VPC, which
// increases the manageability of VPC routing and interoperability with third-party
// workloads.
//
// Route server supports the follow route table types:
//
//   - VPC route tables not associated with subnets
//
//   - Subnet route tables
//
//   - Internet gateway route tables
//
// Route server does not support route tables associated with virtual private
// gateways. To propagate routes into a transit gateway route table, use [Transit Gateway Connect].
//
// For more information see [Dynamic routing in your VPC with VPC Route Server] in the Amazon VPC User Guide.
//
// [Dynamic routing in your VPC with VPC Route Server]: https://docs.aws.amazon.com/vpc/latest/userguide/dynamic-routing-route-server.html
// [Transit Gateway Connect]: https://docs.aws.amazon.com/vpc/latest/tgw/tgw-connect.html
func (c *Client) ModifyRouteServer(ctx context.Context, params *ModifyRouteServerInput, optFns ...func(*Options)) (*ModifyRouteServerOutput, error) {
	if params == nil {
		params = &ModifyRouteServerInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ModifyRouteServer", params, optFns, c.addOperationModifyRouteServerMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ModifyRouteServerOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ModifyRouteServerInput struct {

	// The ID of the route server to modify.
	//
	// This member is required.
	RouteServerId *string

	// A check for whether you have the required permissions for the action without
	// actually making the request and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// Specifies whether to persist routes after all BGP sessions are terminated.
	//
	//   - enable: Routes will be persisted in FIB and RIB after all BGP sessions are
	//   terminated.
	//
	//   - disable: Routes will not be persisted in FIB and RIB after all BGP sessions
	//   are terminated.
	//
	//   - reset: If a route server has persisted routes due to all BGP sessions
	//   having ended, reset will withdraw all routes and reset route server to an empty
	//   FIB and RIB.
	PersistRoutes types.RouteServerPersistRoutesAction

	// The number of minutes a route server will wait after BGP is re-established to
	// unpersist the routes in the FIB and RIB. Value must be in the range of 1-5.
	// Required if PersistRoutes is enabled .
	//
	// If you set the duration to 1 minute, then when your network appliance
	// re-establishes BGP with route server, it has 1 minute to relearn it's adjacent
	// network and advertise those routes to route server before route server resumes
	// normal functionality. In most cases, 1 minute is probably sufficient. If,
	// however, you have concerns that your BGP network may not be capable of fully
	// re-establishing and re-learning everything in 1 minute, you can increase the
	// duration up to 5 minutes.
	PersistRoutesDuration *int64

	// Specifies whether to enable SNS notifications for route server events. Enabling
	// SNS notifications persists BGP status changes to an SNS topic provisioned by
	// Amazon Web Services.
	SnsNotificationsEnabled *bool

	noSmithyDocumentSerde
}

type ModifyRouteServerOutput struct {

	// Information about the modified route server.
	RouteServer *types.RouteServer

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationModifyRouteServerMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpModifyRouteServer{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpModifyRouteServer{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ModifyRouteServer"); err != nil {
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
	if err = addOpModifyRouteServerValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opModifyRouteServer(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opModifyRouteServer(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ModifyRouteServer",
	}
}
