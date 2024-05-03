package main

import (
	"go-aws/lambda/database"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoAwsStackProps struct {
	awscdk.StackProps
}

func NewGoAwsStack(scope constructs.Construct, id string, props *GoAwsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	usersTable := awsdynamodb.NewTable(stack, jsii.String("MyUserTable"), &awsdynamodb.TableProps{
		TableName: jsii.String(database.USERS_TABLE_NAME),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	myFunction := awslambda.NewFunction(stack, jsii.String("MyFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	apiGateway := awsapigateway.NewRestApi(
		stack,
		jsii.String("MyApiGateway"),
		&awsapigateway.RestApiProps{
			DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
				AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
				AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
				AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"),
			},
			CloudWatchRole: jsii.Bool(true),
			DeployOptions: &awsapigateway.StageOptions{
				LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
			},
		},
	)

	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)
	// define the routes
	apiGateway.Root().
		AddResource(jsii.String("register"), nil).
		AddMethod(jsii.String("POST"), integration, nil)

	apiGateway.Root().
		AddResource(jsii.String("login"), nil).
		AddMethod(jsii.String("POST"), integration, nil)

	usersTable.GrantReadWriteData(myFunction)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoAwsStack(app, "GoAwsStack", &GoAwsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
