// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudformation

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Updates a generated template. This can be used to change the name, add and
// remove resources, refresh resources, and change the DeletionPolicy and
// UpdateReplacePolicy settings. You can check the status of the update to the
// generated template using the DescribeGeneratedTemplate API action.
func (c *Client) UpdateGeneratedTemplate(ctx context.Context, params *UpdateGeneratedTemplateInput, optFns ...func(*Options)) (*UpdateGeneratedTemplateOutput, error) {
	if params == nil {
		params = &UpdateGeneratedTemplateInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateGeneratedTemplate", params, optFns, c.addOperationUpdateGeneratedTemplateMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateGeneratedTemplateOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateGeneratedTemplateInput struct {

	// The name or Amazon Resource Name (ARN) of a generated template.
	//
	// This member is required.
	GeneratedTemplateName *string

	// An optional list of resources to be added to the generated template.
	AddResources []types.ResourceDefinition

	// An optional new name to assign to the generated template.
	NewGeneratedTemplateName *string

	// If true , update the resource properties in the generated template with their
	// current live state. This feature is useful when the resource properties in your
	// generated a template does not reflect the live state of the resource properties.
	// This happens when a user update the resource properties after generating a
	// template.
	RefreshAllResources *bool

	// A list of logical ids for resources to remove from the generated template.
	RemoveResources []string

	// The configuration details of the generated template, including the
	// DeletionPolicy and UpdateReplacePolicy .
	TemplateConfiguration *types.TemplateConfiguration

	noSmithyDocumentSerde
}

type UpdateGeneratedTemplateOutput struct {

	// The Amazon Resource Name (ARN) of the generated template. The format is
	// arn:${Partition}:cloudformation:${Region}:${Account}:generatedtemplate/${Id} .
	// For example,
	// arn:aws:cloudformation:us-east-1:123456789012:generatedtemplate/2e8465c1-9a80-43ea-a3a3-4f2d692fe6dc
	// .
	GeneratedTemplateId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateGeneratedTemplateMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpUpdateGeneratedTemplate{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpUpdateGeneratedTemplate{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateGeneratedTemplate"); err != nil {
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
	if err = addOpUpdateGeneratedTemplateValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateGeneratedTemplate(options.Region), middleware.Before); err != nil {
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
	return nil
}

func newServiceMetadataMiddleware_opUpdateGeneratedTemplate(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateGeneratedTemplate",
	}
}
