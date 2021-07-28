// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a data feed for Spot Instances, enabling you to view Spot Instance usage
// logs. You can create one data feed per account. For more information, see Spot
// Instance data feed
// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-data-feeds.html) in
// the Amazon EC2 User Guide for Linux Instances.
func (c *Client) CreateSpotDatafeedSubscription(ctx context.Context, params *CreateSpotDatafeedSubscriptionInput, optFns ...func(*Options)) (*CreateSpotDatafeedSubscriptionOutput, error) {
	if params == nil {
		params = &CreateSpotDatafeedSubscriptionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateSpotDatafeedSubscription", params, optFns, c.addOperationCreateSpotDatafeedSubscriptionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateSpotDatafeedSubscriptionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Contains the parameters for CreateSpotDatafeedSubscription.
type CreateSpotDatafeedSubscriptionInput struct {

	// The name of the Amazon S3 bucket in which to store the Spot Instance data feed.
	// For more information about bucket names, see Rules for bucket naming
	// (https://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html#bucketnamingrules)
	// in the Amazon S3 Developer Guide.
	//
	// This member is required.
	Bucket *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun *bool

	// The prefix for the data feed file names.
	Prefix *string
}

// Contains the output of CreateSpotDatafeedSubscription.
type CreateSpotDatafeedSubscriptionOutput struct {

	// The Spot Instance data feed subscription.
	SpotDatafeedSubscription *types.SpotDatafeedSubscription

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationCreateSpotDatafeedSubscriptionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateSpotDatafeedSubscription{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateSpotDatafeedSubscription{}, middleware.After)
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
	if err = addOpCreateSpotDatafeedSubscriptionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateSpotDatafeedSubscription(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateSpotDatafeedSubscription(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "CreateSpotDatafeedSubscription",
	}
}