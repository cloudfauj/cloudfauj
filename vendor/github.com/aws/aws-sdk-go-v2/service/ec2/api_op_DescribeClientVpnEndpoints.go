// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Describes one or more Client VPN endpoints in the account.
func (c *Client) DescribeClientVpnEndpoints(ctx context.Context, params *DescribeClientVpnEndpointsInput, optFns ...func(*Options)) (*DescribeClientVpnEndpointsOutput, error) {
	if params == nil {
		params = &DescribeClientVpnEndpointsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeClientVpnEndpoints", params, optFns, c.addOperationDescribeClientVpnEndpointsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeClientVpnEndpointsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeClientVpnEndpointsInput struct {

	// The ID of the Client VPN endpoint.
	ClientVpnEndpointIds []string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// One or more filters. Filter names and values are case-sensitive.
	//
	// * endpoint-id
	// - The ID of the Client VPN endpoint.
	//
	// * transport-protocol - The transport
	// protocol (tcp | udp).
	Filters []types.Filter

	// The maximum number of results to return for the request in a single page. The
	// remaining results can be seen by sending another request with the nextToken
	// value.
	MaxResults *int32

	// The token to retrieve the next page of results.
	NextToken *string
}

type DescribeClientVpnEndpointsOutput struct {

	// Information about the Client VPN endpoints.
	ClientVpnEndpoints []types.ClientVpnEndpoint

	// The token to use to retrieve the next page of results. This value is null when
	// there are no more results to return.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationDescribeClientVpnEndpointsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeClientVpnEndpoints{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeClientVpnEndpoints{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeClientVpnEndpoints(options.Region), middleware.Before); err != nil {
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
	return nil
}

// DescribeClientVpnEndpointsAPIClient is a client that implements the
// DescribeClientVpnEndpoints operation.
type DescribeClientVpnEndpointsAPIClient interface {
	DescribeClientVpnEndpoints(context.Context, *DescribeClientVpnEndpointsInput, ...func(*Options)) (*DescribeClientVpnEndpointsOutput, error)
}

var _ DescribeClientVpnEndpointsAPIClient = (*Client)(nil)

// DescribeClientVpnEndpointsPaginatorOptions is the paginator options for
// DescribeClientVpnEndpoints
type DescribeClientVpnEndpointsPaginatorOptions struct {
	// The maximum number of results to return for the request in a single page. The
	// remaining results can be seen by sending another request with the nextToken
	// value.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeClientVpnEndpointsPaginator is a paginator for
// DescribeClientVpnEndpoints
type DescribeClientVpnEndpointsPaginator struct {
	options   DescribeClientVpnEndpointsPaginatorOptions
	client    DescribeClientVpnEndpointsAPIClient
	params    *DescribeClientVpnEndpointsInput
	nextToken *string
	firstPage bool
}

// NewDescribeClientVpnEndpointsPaginator returns a new
// DescribeClientVpnEndpointsPaginator
func NewDescribeClientVpnEndpointsPaginator(client DescribeClientVpnEndpointsAPIClient, params *DescribeClientVpnEndpointsInput, optFns ...func(*DescribeClientVpnEndpointsPaginatorOptions)) *DescribeClientVpnEndpointsPaginator {
	if params == nil {
		params = &DescribeClientVpnEndpointsInput{}
	}

	options := DescribeClientVpnEndpointsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeClientVpnEndpointsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeClientVpnEndpointsPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next DescribeClientVpnEndpoints page.
func (p *DescribeClientVpnEndpointsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeClientVpnEndpointsOutput, error) {
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

	result, err := p.client.DescribeClientVpnEndpoints(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken && prevToken != nil && p.nextToken != nil && *prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opDescribeClientVpnEndpoints(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "DescribeClientVpnEndpoints",
	}
}