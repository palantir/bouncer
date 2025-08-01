// Code generated by smithy-go-codegen DO NOT EDIT.

package autoscaling

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Gets information about the instance refreshes for the specified Auto Scaling
// group from the previous six weeks.
//
// This operation is part of the [instance refresh feature] in Amazon EC2 Auto Scaling, which helps you
// update instances in your Auto Scaling group after you make configuration
// changes.
//
// To help you determine the status of an instance refresh, Amazon EC2 Auto
// Scaling returns information about the instance refreshes you previously
// initiated, including their status, start time, end time, the percentage of the
// instance refresh that is complete, and the number of instances remaining to
// update before the instance refresh is complete. If a rollback is initiated while
// an instance refresh is in progress, Amazon EC2 Auto Scaling also returns
// information about the rollback of the instance refresh.
//
// [instance refresh feature]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-instance-refresh.html
func (c *Client) DescribeInstanceRefreshes(ctx context.Context, params *DescribeInstanceRefreshesInput, optFns ...func(*Options)) (*DescribeInstanceRefreshesOutput, error) {
	if params == nil {
		params = &DescribeInstanceRefreshesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeInstanceRefreshes", params, optFns, c.addOperationDescribeInstanceRefreshesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeInstanceRefreshesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeInstanceRefreshesInput struct {

	// The name of the Auto Scaling group.
	//
	// This member is required.
	AutoScalingGroupName *string

	// One or more instance refresh IDs.
	InstanceRefreshIds []string

	// The maximum number of items to return with this call. The default value is 50
	// and the maximum value is 100 .
	MaxRecords *int32

	// The token for the next set of items to return. (You received this token from a
	// previous call.)
	NextToken *string

	noSmithyDocumentSerde
}

type DescribeInstanceRefreshesOutput struct {

	// The instance refreshes for the specified group, sorted by creation timestamp in
	// descending order.
	InstanceRefreshes []types.InstanceRefresh

	// A string that indicates that the response contains more items than can be
	// returned in a single response. To receive additional items, specify this string
	// for the NextToken value when requesting the next set of items. This value is
	// null when there are no more items to return.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeInstanceRefreshesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpDescribeInstanceRefreshes{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpDescribeInstanceRefreshes{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeInstanceRefreshes"); err != nil {
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
	if err = addOpDescribeInstanceRefreshesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeInstanceRefreshes(options.Region), middleware.Before); err != nil {
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

// DescribeInstanceRefreshesPaginatorOptions is the paginator options for
// DescribeInstanceRefreshes
type DescribeInstanceRefreshesPaginatorOptions struct {
	// The maximum number of items to return with this call. The default value is 50
	// and the maximum value is 100 .
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeInstanceRefreshesPaginator is a paginator for DescribeInstanceRefreshes
type DescribeInstanceRefreshesPaginator struct {
	options   DescribeInstanceRefreshesPaginatorOptions
	client    DescribeInstanceRefreshesAPIClient
	params    *DescribeInstanceRefreshesInput
	nextToken *string
	firstPage bool
}

// NewDescribeInstanceRefreshesPaginator returns a new
// DescribeInstanceRefreshesPaginator
func NewDescribeInstanceRefreshesPaginator(client DescribeInstanceRefreshesAPIClient, params *DescribeInstanceRefreshesInput, optFns ...func(*DescribeInstanceRefreshesPaginatorOptions)) *DescribeInstanceRefreshesPaginator {
	if params == nil {
		params = &DescribeInstanceRefreshesInput{}
	}

	options := DescribeInstanceRefreshesPaginatorOptions{}
	if params.MaxRecords != nil {
		options.Limit = *params.MaxRecords
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeInstanceRefreshesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeInstanceRefreshesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeInstanceRefreshes page.
func (p *DescribeInstanceRefreshesPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeInstanceRefreshesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxRecords = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.DescribeInstanceRefreshes(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// DescribeInstanceRefreshesAPIClient is a client that implements the
// DescribeInstanceRefreshes operation.
type DescribeInstanceRefreshesAPIClient interface {
	DescribeInstanceRefreshes(context.Context, *DescribeInstanceRefreshesInput, ...func(*Options)) (*DescribeInstanceRefreshesOutput, error)
}

var _ DescribeInstanceRefreshesAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opDescribeInstanceRefreshes(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeInstanceRefreshes",
	}
}
