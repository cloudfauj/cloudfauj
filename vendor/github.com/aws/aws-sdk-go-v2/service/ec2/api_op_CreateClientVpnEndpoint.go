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

// Creates a Client VPN endpoint. A Client VPN endpoint is the resource you create
// and configure to enable and manage client VPN sessions. It is the destination
// endpoint at which all client VPN sessions are terminated.
func (c *Client) CreateClientVpnEndpoint(ctx context.Context, params *CreateClientVpnEndpointInput, optFns ...func(*Options)) (*CreateClientVpnEndpointOutput, error) {
	if params == nil {
		params = &CreateClientVpnEndpointInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateClientVpnEndpoint", params, optFns, c.addOperationCreateClientVpnEndpointMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateClientVpnEndpointOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateClientVpnEndpointInput struct {

	// Information about the authentication method to be used to authenticate clients.
	//
	// This member is required.
	AuthenticationOptions []types.ClientVpnAuthenticationRequest

	// The IPv4 address range, in CIDR notation, from which to assign client IP
	// addresses. The address range cannot overlap with the local CIDR of the VPC in
	// which the associated subnet is located, or the routes that you add manually. The
	// address range cannot be changed after the Client VPN endpoint has been created.
	// The CIDR block should be /22 or greater.
	//
	// This member is required.
	ClientCidrBlock *string

	// Information about the client connection logging options. If you enable client
	// connection logging, data about client connections is sent to a Cloudwatch Logs
	// log stream. The following information is logged:
	//
	// * Client connection
	// requests
	//
	// * Client connection results (successful and unsuccessful)
	//
	// * Reasons
	// for unsuccessful client connection requests
	//
	// * Client connection termination
	// time
	//
	// This member is required.
	ConnectionLogOptions *types.ConnectionLogOptions

	// The ARN of the server certificate. For more information, see the AWS Certificate
	// Manager User Guide (https://docs.aws.amazon.com/acm/latest/userguide/).
	//
	// This member is required.
	ServerCertificateArn *string

	// The options for managing connection authorization for new client connections.
	ClientConnectOptions *types.ClientConnectOptions

	// Unique, case-sensitive identifier that you provide to ensure the idempotency of
	// the request. For more information, see How to Ensure Idempotency
	// (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	ClientToken *string

	// A brief description of the Client VPN endpoint.
	Description *string

	// Information about the DNS servers to be used for DNS resolution. A Client VPN
	// endpoint can have up to two DNS servers. If no DNS server is specified, the DNS
	// address configured on the device is used for the DNS server.
	DnsServers []string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// The IDs of one or more security groups to apply to the target network. You must
	// also specify the ID of the VPC that contains the security groups.
	SecurityGroupIds []string

	// Specify whether to enable the self-service portal for the Client VPN endpoint.
	// Default Value: enabled
	SelfServicePortal types.SelfServicePortal

	// Indicates whether split-tunnel is enabled on the AWS Client VPN endpoint. By
	// default, split-tunnel on a VPN endpoint is disabled. For information about
	// split-tunnel VPN endpoints, see Split-Tunnel AWS Client VPN Endpoint
	// (https://docs.aws.amazon.com/vpn/latest/clientvpn-admin/split-tunnel-vpn.html)
	// in the AWS Client VPN Administrator Guide.
	SplitTunnel *bool

	// The tags to apply to the Client VPN endpoint during creation.
	TagSpecifications []types.TagSpecification

	// The transport protocol to be used by the VPN session. Default value: udp
	TransportProtocol types.TransportProtocol

	// The ID of the VPC to associate with the Client VPN endpoint. If no security
	// group IDs are specified in the request, the default security group for the VPC
	// is applied.
	VpcId *string

	// The port number to assign to the Client VPN endpoint for TCP and UDP traffic.
	// Valid Values: 443 | 1194 Default Value: 443
	VpnPort *int32
}

type CreateClientVpnEndpointOutput struct {

	// The ID of the Client VPN endpoint.
	ClientVpnEndpointId *string

	// The DNS name to be used by clients when establishing their VPN session.
	DnsName *string

	// The current state of the Client VPN endpoint.
	Status *types.ClientVpnEndpointStatus

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationCreateClientVpnEndpointMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateClientVpnEndpoint{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateClientVpnEndpoint{}, middleware.After)
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
	if err = addIdempotencyToken_opCreateClientVpnEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateClientVpnEndpointValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateClientVpnEndpoint(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpCreateClientVpnEndpoint struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpCreateClientVpnEndpoint) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpCreateClientVpnEndpoint) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*CreateClientVpnEndpointInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *CreateClientVpnEndpointInput ")
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
func addIdempotencyToken_opCreateClientVpnEndpointMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpCreateClientVpnEndpoint{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opCreateClientVpnEndpoint(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "CreateClientVpnEndpoint",
	}
}
