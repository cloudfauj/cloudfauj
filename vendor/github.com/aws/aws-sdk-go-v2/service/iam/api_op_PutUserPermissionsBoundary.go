// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Adds or updates the policy that is specified as the IAM user's permissions
// boundary. You can use an Amazon Web Services managed policy or a customer
// managed policy to set the boundary for a user. Use the boundary to control the
// maximum permissions that the user can have. Setting a permissions boundary is an
// advanced feature that can affect the permissions for the user. Policies that are
// used as permissions boundaries do not provide permissions. You must also attach
// a permissions policy to the user. To learn how the effective permissions for a
// user are evaluated, see IAM JSON policy evaluation logic
// (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_evaluation-logic.html)
// in the IAM User Guide.
func (c *Client) PutUserPermissionsBoundary(ctx context.Context, params *PutUserPermissionsBoundaryInput, optFns ...func(*Options)) (*PutUserPermissionsBoundaryOutput, error) {
	if params == nil {
		params = &PutUserPermissionsBoundaryInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutUserPermissionsBoundary", params, optFns, c.addOperationPutUserPermissionsBoundaryMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutUserPermissionsBoundaryOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutUserPermissionsBoundaryInput struct {

	// The ARN of the policy that is used to set the permissions boundary for the user.
	//
	// This member is required.
	PermissionsBoundary *string

	// The name (friendly name, not ARN) of the IAM user for which you want to set the
	// permissions boundary.
	//
	// This member is required.
	UserName *string
}

type PutUserPermissionsBoundaryOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationPutUserPermissionsBoundaryMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpPutUserPermissionsBoundary{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpPutUserPermissionsBoundary{}, middleware.After)
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
	if err = addOpPutUserPermissionsBoundaryValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutUserPermissionsBoundary(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opPutUserPermissionsBoundary(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "iam",
		OperationName: "PutUserPermissionsBoundary",
	}
}
