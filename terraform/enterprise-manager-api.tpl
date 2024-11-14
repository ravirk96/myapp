openapi: 3.0.0
info:
  title: Enterprise Manager API
  version: 1.0.0
paths:
  /item:
    post:
      summary: Create an item
      operationId: createItem
      x-amazon-apigateway-integration:
        uri: "${lambda_arn}"
        httpMethod: POST
        type: "AWS_PROXY"
      responses:
        '200':
          description: Successful operation
    put:
      summary: Update an item
      operationId: updateItem
      x-amazon-apigateway-integration:
        uri: "${lambda_arn}"
        httpMethod: POST
        type: "AWS_PROXY"
      responses:
        '200':
          description: Successful operation
  /item/{id}:
    get:
      summary: Get an item by ID
      operationId: getItem
      x-amazon-apigateway-integration:
        uri: "${lambda_arn}"
        httpMethod: POST
        type: "AWS_PROXY"
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Successful operation