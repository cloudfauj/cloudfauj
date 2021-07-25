// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// Retrieves a credential report for the account. For more information about the
// credential report, see Getting credential reports
// (https://docs.aws.amazon.com/IAM/latest/UserGuide/credential-reports.html) in
// the IAM User Guide.
func (c *Client) GetCredentialReport(ctx context.Context, params *GetCredentialReportInput, optFns ...func(*Options)) (*GetCredentialReportOutput, error) {
	if params == nil {
		params = &GetCredentialReportInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetCredentialReport", params, optFns, c.addOperationGetCredentialReportMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetCredentialReportOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetCredentialReportInput struct {
}

// Contains the response to a successful GetCredentialReport request.
type GetCredentialReportOutput struct {

	// Contains the credential report. The report is Base64-encoded.
	Content []byte

	// The date and time when the credential report was created, in ISO 8601 date-time
	// format (http://www.iso.org/iso/iso8601).
	GeneratedTime *time.Time

	// The format (MIME type) of the credential report.
	ReportFormat types.ReportFormatType

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (c *Client) addOperationGetCredentialReportMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpGetCredentialReport{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpGetCredentialReport{}, middleware.After)
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetCredentialReport(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opGetCredentialReport(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "iam",
		OperationName: "GetCredentialReport",
	}
}
