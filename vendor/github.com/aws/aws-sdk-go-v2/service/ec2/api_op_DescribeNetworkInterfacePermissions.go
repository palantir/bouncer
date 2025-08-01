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

// Describes the permissions for your network interfaces.
func (c *Client) DescribeNetworkInterfacePermissions(ctx context.Context, params *DescribeNetworkInterfacePermissionsInput, optFns ...func(*Options)) (*DescribeNetworkInterfacePermissionsOutput, error) {
	if params == nil {
		params = &DescribeNetworkInterfacePermissionsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeNetworkInterfacePermissions", params, optFns, c.addOperationDescribeNetworkInterfacePermissionsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeNetworkInterfacePermissionsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Contains the parameters for DescribeNetworkInterfacePermissions.
type DescribeNetworkInterfacePermissionsInput struct {

	// One or more filters.
	//
	//   - network-interface-permission.network-interface-permission-id - The ID of the
	//   permission.
	//
	//   - network-interface-permission.network-interface-id - The ID of the network
	//   interface.
	//
	//   - network-interface-permission.aws-account-id - The Amazon Web Services
	//   account ID.
	//
	//   - network-interface-permission.aws-service - The Amazon Web Services service.
	//
	//   - network-interface-permission.permission - The type of permission (
	//   INSTANCE-ATTACH | EIP-ASSOCIATE ).
	Filters []types.Filter

	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. If this
	// parameter is not specified, up to 50 results are returned by default. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	MaxResults *int32

	// The network interface permission IDs.
	NetworkInterfacePermissionIds []string

	// The token returned from a previous paginated request. Pagination continues from
	// the end of the items returned by the previous request.
	NextToken *string

	noSmithyDocumentSerde
}

// Contains the output for DescribeNetworkInterfacePermissions.
type DescribeNetworkInterfacePermissionsOutput struct {

	// The network interface permissions.
	NetworkInterfacePermissions []types.NetworkInterfacePermission

	// The token to include in another request to get the next page of items. This
	// value is null when there are no more items to return.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeNetworkInterfacePermissionsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeNetworkInterfacePermissions{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeNetworkInterfacePermissions{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeNetworkInterfacePermissions"); err != nil {
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeNetworkInterfacePermissions(options.Region), middleware.Before); err != nil {
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

// DescribeNetworkInterfacePermissionsPaginatorOptions is the paginator options
// for DescribeNetworkInterfacePermissions
type DescribeNetworkInterfacePermissionsPaginatorOptions struct {
	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. If this
	// parameter is not specified, up to 50 results are returned by default. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeNetworkInterfacePermissionsPaginator is a paginator for
// DescribeNetworkInterfacePermissions
type DescribeNetworkInterfacePermissionsPaginator struct {
	options   DescribeNetworkInterfacePermissionsPaginatorOptions
	client    DescribeNetworkInterfacePermissionsAPIClient
	params    *DescribeNetworkInterfacePermissionsInput
	nextToken *string
	firstPage bool
}

// NewDescribeNetworkInterfacePermissionsPaginator returns a new
// DescribeNetworkInterfacePermissionsPaginator
func NewDescribeNetworkInterfacePermissionsPaginator(client DescribeNetworkInterfacePermissionsAPIClient, params *DescribeNetworkInterfacePermissionsInput, optFns ...func(*DescribeNetworkInterfacePermissionsPaginatorOptions)) *DescribeNetworkInterfacePermissionsPaginator {
	if params == nil {
		params = &DescribeNetworkInterfacePermissionsInput{}
	}

	options := DescribeNetworkInterfacePermissionsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeNetworkInterfacePermissionsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeNetworkInterfacePermissionsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeNetworkInterfacePermissions page.
func (p *DescribeNetworkInterfacePermissionsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeNetworkInterfacePermissionsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.DescribeNetworkInterfacePermissions(ctx, &params, optFns...)
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

// DescribeNetworkInterfacePermissionsAPIClient is a client that implements the
// DescribeNetworkInterfacePermissions operation.
type DescribeNetworkInterfacePermissionsAPIClient interface {
	DescribeNetworkInterfacePermissions(context.Context, *DescribeNetworkInterfacePermissionsInput, ...func(*Options)) (*DescribeNetworkInterfacePermissionsOutput, error)
}

var _ DescribeNetworkInterfacePermissionsAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opDescribeNetworkInterfacePermissions(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeNetworkInterfacePermissions",
	}
}
