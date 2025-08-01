// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	smithywaiter "github.com/aws/smithy-go/waiter"
	"time"
)

// Describes one or more of your VPN customer gateways.
//
// For more information, see [Amazon Web Services Site-to-Site VPN] in the Amazon Web Services Site-to-Site VPN User
// Guide.
//
// [Amazon Web Services Site-to-Site VPN]: https://docs.aws.amazon.com/vpn/latest/s2svpn/VPC_VPN.html
func (c *Client) DescribeCustomerGateways(ctx context.Context, params *DescribeCustomerGatewaysInput, optFns ...func(*Options)) (*DescribeCustomerGatewaysOutput, error) {
	if params == nil {
		params = &DescribeCustomerGatewaysInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeCustomerGateways", params, optFns, c.addOperationDescribeCustomerGatewaysMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeCustomerGatewaysOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Contains the parameters for DescribeCustomerGateways.
type DescribeCustomerGatewaysInput struct {

	// One or more customer gateway IDs.
	//
	// Default: Describes all your customer gateways.
	CustomerGatewayIds []string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// One or more filters.
	//
	//   - bgp-asn - The customer gateway's Border Gateway Protocol (BGP) Autonomous
	//   System Number (ASN).
	//
	//   - customer-gateway-id - The ID of the customer gateway.
	//
	//   - ip-address - The IP address of the customer gateway device's external
	//   interface.
	//
	//   - state - The state of the customer gateway ( pending | available | deleting |
	//   deleted ).
	//
	//   - type - The type of customer gateway. Currently, the only supported type is
	//   ipsec.1 .
	//
	//   - tag : - The key/value combination of a tag assigned to the resource. Use the
	//   tag key in the filter name and the tag value as the filter value. For example,
	//   to find all resources that have a tag with the key Owner and the value TeamA ,
	//   specify tag:Owner for the filter name and TeamA for the filter value.
	//
	//   - tag-key - The key of a tag assigned to the resource. Use this filter to find
	//   all resources assigned a tag with a specific key, regardless of the tag value.
	Filters []types.Filter

	noSmithyDocumentSerde
}

// Contains the output of DescribeCustomerGateways.
type DescribeCustomerGatewaysOutput struct {

	// Information about one or more customer gateways.
	CustomerGateways []types.CustomerGateway

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeCustomerGatewaysMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeCustomerGateways{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeCustomerGateways{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeCustomerGateways"); err != nil {
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeCustomerGateways(options.Region), middleware.Before); err != nil {
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

// CustomerGatewayAvailableWaiterOptions are waiter options for
// CustomerGatewayAvailableWaiter
type CustomerGatewayAvailableWaiterOptions struct {

	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behavior.
	//
	// Passing options here is functionally equivalent to passing values to this
	// config's ClientOptions field that extend the inner client's APIOptions directly.
	APIOptions []func(*middleware.Stack) error

	// Functional options to be passed to all operations invoked by this client.
	//
	// Function values that modify the inner APIOptions are applied after the waiter
	// config's own APIOptions modifiers.
	ClientOptions []func(*Options)

	// MinDelay is the minimum amount of time to delay between retries. If unset,
	// CustomerGatewayAvailableWaiter will use default minimum delay of 15 seconds.
	// Note that MinDelay must resolve to a value lesser than or equal to the MaxDelay.
	MinDelay time.Duration

	// MaxDelay is the maximum amount of time to delay between retries. If unset or
	// set to zero, CustomerGatewayAvailableWaiter will use default max delay of 120
	// seconds. Note that MaxDelay must resolve to value greater than or equal to the
	// MinDelay.
	MaxDelay time.Duration

	// LogWaitAttempts is used to enable logging for waiter retry attempts
	LogWaitAttempts bool

	// Retryable is function that can be used to override the service defined
	// waiter-behavior based on operation output, or returned error. This function is
	// used by the waiter to decide if a state is retryable or a terminal state.
	//
	// By default service-modeled logic will populate this option. This option can
	// thus be used to define a custom waiter state with fall-back to service-modeled
	// waiter state mutators.The function returns an error in case of a failure state.
	// In case of retry state, this function returns a bool value of true and nil
	// error, while in case of success it returns a bool value of false and nil error.
	Retryable func(context.Context, *DescribeCustomerGatewaysInput, *DescribeCustomerGatewaysOutput, error) (bool, error)
}

// CustomerGatewayAvailableWaiter defines the waiters for CustomerGatewayAvailable
type CustomerGatewayAvailableWaiter struct {
	client DescribeCustomerGatewaysAPIClient

	options CustomerGatewayAvailableWaiterOptions
}

// NewCustomerGatewayAvailableWaiter constructs a CustomerGatewayAvailableWaiter.
func NewCustomerGatewayAvailableWaiter(client DescribeCustomerGatewaysAPIClient, optFns ...func(*CustomerGatewayAvailableWaiterOptions)) *CustomerGatewayAvailableWaiter {
	options := CustomerGatewayAvailableWaiterOptions{}
	options.MinDelay = 15 * time.Second
	options.MaxDelay = 120 * time.Second
	options.Retryable = customerGatewayAvailableStateRetryable

	for _, fn := range optFns {
		fn(&options)
	}
	return &CustomerGatewayAvailableWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for CustomerGatewayAvailable waiter. The
// maxWaitDur is the maximum wait duration the waiter will wait. The maxWaitDur is
// required and must be greater than zero.
func (w *CustomerGatewayAvailableWaiter) Wait(ctx context.Context, params *DescribeCustomerGatewaysInput, maxWaitDur time.Duration, optFns ...func(*CustomerGatewayAvailableWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for CustomerGatewayAvailable waiter and
// returns the output of the successful operation. The maxWaitDur is the maximum
// wait duration the waiter will wait. The maxWaitDur is required and must be
// greater than zero.
func (w *CustomerGatewayAvailableWaiter) WaitForOutput(ctx context.Context, params *DescribeCustomerGatewaysInput, maxWaitDur time.Duration, optFns ...func(*CustomerGatewayAvailableWaiterOptions)) (*DescribeCustomerGatewaysOutput, error) {
	if maxWaitDur <= 0 {
		return nil, fmt.Errorf("maximum wait time for waiter must be greater than zero")
	}

	options := w.options
	for _, fn := range optFns {
		fn(&options)
	}

	if options.MaxDelay <= 0 {
		options.MaxDelay = 120 * time.Second
	}

	if options.MinDelay > options.MaxDelay {
		return nil, fmt.Errorf("minimum waiter delay %v must be lesser than or equal to maximum waiter delay of %v.", options.MinDelay, options.MaxDelay)
	}

	ctx, cancelFn := context.WithTimeout(ctx, maxWaitDur)
	defer cancelFn()

	logger := smithywaiter.Logger{}
	remainingTime := maxWaitDur

	var attempt int64
	for {

		attempt++
		apiOptions := options.APIOptions
		start := time.Now()

		if options.LogWaitAttempts {
			logger.Attempt = attempt
			apiOptions = append([]func(*middleware.Stack) error{}, options.APIOptions...)
			apiOptions = append(apiOptions, logger.AddLogger)
		}

		out, err := w.client.DescribeCustomerGateways(ctx, params, func(o *Options) {
			baseOpts := []func(*Options){
				addIsWaiterUserAgent,
			}
			o.APIOptions = append(o.APIOptions, apiOptions...)
			for _, opt := range baseOpts {
				opt(o)
			}
			for _, opt := range options.ClientOptions {
				opt(o)
			}
		})

		retryable, err := options.Retryable(ctx, params, out, err)
		if err != nil {
			return nil, err
		}
		if !retryable {
			return out, nil
		}

		remainingTime -= time.Since(start)
		if remainingTime < options.MinDelay || remainingTime <= 0 {
			break
		}

		// compute exponential backoff between waiter retries
		delay, err := smithywaiter.ComputeDelay(
			attempt, options.MinDelay, options.MaxDelay, remainingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error computing waiter delay, %w", err)
		}

		remainingTime -= delay
		// sleep for the delay amount before invoking a request
		if err := smithytime.SleepWithContext(ctx, delay); err != nil {
			return nil, fmt.Errorf("request cancelled while waiting, %w", err)
		}
	}
	return nil, fmt.Errorf("exceeded max wait time for CustomerGatewayAvailable waiter")
}

func customerGatewayAvailableStateRetryable(ctx context.Context, input *DescribeCustomerGatewaysInput, output *DescribeCustomerGatewaysOutput, err error) (bool, error) {

	if err == nil {
		v1 := output.CustomerGateways
		var v2 []string
		for _, v := range v1 {
			v3 := v.State
			if v3 != nil {
				v2 = append(v2, *v3)
			}
		}
		expectedValue := "available"
		match := len(v2) > 0
		for _, v := range v2 {
			if string(v) != expectedValue {
				match = false
				break
			}
		}

		if match {
			return false, nil
		}
	}

	if err == nil {
		v1 := output.CustomerGateways
		var v2 []string
		for _, v := range v1 {
			v3 := v.State
			if v3 != nil {
				v2 = append(v2, *v3)
			}
		}
		expectedValue := "deleted"
		var match bool
		for _, v := range v2 {
			if string(v) == expectedValue {
				match = true
				break
			}
		}

		if match {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	if err == nil {
		v1 := output.CustomerGateways
		var v2 []string
		for _, v := range v1 {
			v3 := v.State
			if v3 != nil {
				v2 = append(v2, *v3)
			}
		}
		expectedValue := "deleting"
		var match bool
		for _, v := range v2 {
			if string(v) == expectedValue {
				match = true
				break
			}
		}

		if match {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

// DescribeCustomerGatewaysAPIClient is a client that implements the
// DescribeCustomerGateways operation.
type DescribeCustomerGatewaysAPIClient interface {
	DescribeCustomerGateways(context.Context, *DescribeCustomerGatewaysInput, ...func(*Options)) (*DescribeCustomerGatewaysOutput, error)
}

var _ DescribeCustomerGatewaysAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opDescribeCustomerGateways(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeCustomerGateways",
	}
}
