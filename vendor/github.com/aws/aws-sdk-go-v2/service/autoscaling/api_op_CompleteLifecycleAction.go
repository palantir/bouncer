// Code generated by smithy-go-codegen DO NOT EDIT.

package autoscaling

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Completes the lifecycle action for the specified token or instance with the
// specified result.
//
// This step is a part of the procedure for adding a lifecycle hook to an Auto
// Scaling group:
//
//   - (Optional) Create a launch template or launch configuration with a user
//     data script that runs while an instance is in a wait state due to a lifecycle
//     hook.
//
//   - (Optional) Create a Lambda function and a rule that allows Amazon
//     EventBridge to invoke your Lambda function when an instance is put into a wait
//     state due to a lifecycle hook.
//
//   - (Optional) Create a notification target and an IAM role. The target can be
//     either an Amazon SQS queue or an Amazon SNS topic. The role allows Amazon EC2
//     Auto Scaling to publish lifecycle notifications to the target.
//
//   - Create the lifecycle hook. Specify whether the hook is used when the
//     instances launch or terminate.
//
//   - If you need more time, record the lifecycle action heartbeat to keep the
//     instance in a wait state.
//
//   - If you finish before the timeout period ends, send a callback by using the [CompleteLifecycleAction]
//     API call.
//
// For more information, see [Complete a lifecycle action] in the Amazon EC2 Auto Scaling User Guide.
//
// [CompleteLifecycleAction]: https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_CompleteLifecycleAction.html
// [Complete a lifecycle action]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/completing-lifecycle-hooks.html
func (c *Client) CompleteLifecycleAction(ctx context.Context, params *CompleteLifecycleActionInput, optFns ...func(*Options)) (*CompleteLifecycleActionOutput, error) {
	if params == nil {
		params = &CompleteLifecycleActionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CompleteLifecycleAction", params, optFns, c.addOperationCompleteLifecycleActionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CompleteLifecycleActionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CompleteLifecycleActionInput struct {

	// The name of the Auto Scaling group.
	//
	// This member is required.
	AutoScalingGroupName *string

	// The action for the group to take. You can specify either CONTINUE or ABANDON .
	//
	// This member is required.
	LifecycleActionResult *string

	// The name of the lifecycle hook.
	//
	// This member is required.
	LifecycleHookName *string

	// The ID of the instance.
	InstanceId *string

	// A universally unique identifier (UUID) that identifies a specific lifecycle
	// action associated with an instance. Amazon EC2 Auto Scaling sends this token to
	// the notification target you specified when you created the lifecycle hook.
	LifecycleActionToken *string

	noSmithyDocumentSerde
}

type CompleteLifecycleActionOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCompleteLifecycleActionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpCompleteLifecycleAction{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpCompleteLifecycleAction{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CompleteLifecycleAction"); err != nil {
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
	if err = addOpCompleteLifecycleActionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCompleteLifecycleAction(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCompleteLifecycleAction(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CompleteLifecycleAction",
	}
}
