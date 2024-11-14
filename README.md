# MyApp

MyApp is a serverless application that uses AWS Lambda, API Gateway, and DynamoDB to create, read, and update items. This project includes all necessary Terraform configurations, a Lambda function written in Go, and a Makefile to build, deploy, and clean up the application.

## Project Structure

- `apigateway.tf`: Terraform configuration for setting up API Gateway and integrating with Lambda.
- `dynamodb.tf`: Terraform configuration for setting up a DynamoDB table.
- `lambda.tf`: Terraform configuration for setting up Lambda function and related IAM roles.
- `output.tf`: Terraform configuration for output variables.
- `provider.tf`: Terraform configuration for specifying the AWS provider.
- `build.sh`: Shell script to build the Go binary and package it into a zip file.
- `Makefile`: Makefile to build, deploy, and clean the application.
- `main.go`: Go code for the Lambda function.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads)
- [Go](https://golang.org/dl/)
- [AWS CLI](https://aws.amazon.com/cli/)
- AWS account with appropriate permissions to create resources.

## Setup

1. **Clone the repository:**
    ```bash
    git clone https://github.com/ravirk96/myapp.git
    cd myapp
    ```

2. **Configure AWS CLI:**
    ```bash
    duo-sso --interactive  
    ```

3. **Build the Go binary and package it:**
    ```bash
    make build
    ```

4. **Deploy the application:**
    ```bash
    make deploy
    ```

5. **Get the API URL:**
    After deployment, Terraform will output the API URL. You can also find it by running:
    ```bash
    terraform output api_url
    ```

## API Endpoints

- **dev url**: Steps to get the Dev URL.
    - Open AWS console.
    - API Gateway > APIs > go-api (23j5ffs6ki) (Select your API) > Stages
    - Stage details > copy Invoke URL
    - use it in Postman for routing.

- **POST /item**: Create a new item.
    - Request Body:
    ```json
    {
      "ID": "item1",
      "Name": "Item 1",
      "Value": "100"
    }
    ```
- **PUT /item**: Update an existing item.
    - Request Body:
    ```json
    {
      "ID": "item1",
      "Name": "Updated Item 1",
      "Value": "200"
    }
    ```
- **GET /item/{id}**: Get an item by ID.
    - Response Body:
    ```json
    {
      "ID": "item1",
      "Name": "Item 1",
      "Value": "100"
    }
    ```

## Clean Up

To destroy the resources created by Terraform:
```bash
make clean