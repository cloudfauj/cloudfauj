// Code generated by smithy-go-codegen DO NOT EDIT.

package elasticloadbalancingv2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Replaces the specified properties of the specified listener. Any properties that
// you do not specify remain unchanged. Changing the protocol from HTTPS to HTTP,
// or from TLS to TCP, removes the security policy and default certificate
// properties. If you change the protocol from HTTP to HTTPS, or from TCP to TLS,
// you must add the security policy and default certificate properties. To add an
// item to a list, remove an item from a list, or update an item in a list, you
// must provide the entire list. For example, to add an action, specify a list with
// the current actions plus the new action.
func (c *Client) ModifyListener(ctx context.Context, params *ModifyListenerInput, optFns ...func(*Options)) (*ModifyListenerOutput, error) {
	if params == nil {
		params = &ModifyListenerInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ModifyListener", params, optFns, c.addOperationModifyListenerMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ModifyListenerOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ModifyListenerInput struct {

	// The Amazon Resource Name (ARN) of the listener.
	//
	// This member is required.
	ListenerArn *string

	// [TLS listeners] The name of the Application-Layer Protocol Negotiation (ALPN)
	// policy. You can specify one policy name. The following are the possible
	// values:
	//
	// * HTTP1Only
	//
	// * HTTP2Only
	//
	// * HTTP2Optional
	//
	// * HTTP2Preferred
	//
	// *
	// None
	//
	// For more information, see ALPN policies
	// (https://docs.aws.amazon.com/elasticloadbalancing/latest/network/create-tls-listener.html#alpn-policies)
	// in the Network Load Balancers Guide.
	AlpnPolicy []string

	// [HTTPS and TLS listeners] The default certificate for the listener. You must
	// provide exactly one certificate. Set CertificateArn to the certificate ARN but
	// do not set IsDefault.
	Certificates []types.Certificate

	// The actions for the default rule.
	DefaultActions []types.Action

	// The port for connections from clients to the load balancer. You cannot specify a
	// port for a Gateway Load Balancer.
	Port *int32

	// The protocol for connections from clients to the load balancer. Application Load
	// Balancers support the HTTP and HTTPS protocols. Network Load Balancers support
	// the TCP, TLS, UDP, and TCP_UDP protocols. You can’t change the protocol to UDP
	// or TCP_UDP if dual-stack mode is enabled. You cannot specify a protocol for a
	// Gateway Load Balancer.
	Protocol types.ProtocolEnum

	// [HTTPS and TLS listeners] The security policy that defines which protocols and
	// ciphers are supported. For more information, see Security policies
	// (https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html#describe-ssl-policies)
	// in the Application Load Balancers Guide or Security policies
	// (https://docs.aws.amazon.com/elasticloadbalancing/latest/network/create-tls-listener.html#describe-ssl-policies)
	// in the Network Load Balancers Guide.
	SslPolicy *string
}

type ModifyListenerOutput struct {

	// Information about the modified listener.
	Listeners []types.Listener

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationModifyListenerMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpModifyListener{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpModifyListener{}, middleware.After)
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
	if err = addOpModifyListenerValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opModifyListener(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opModifyListener(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "elasticloadbalancing",
		OperationName: "ModifyListener",
	}
}
