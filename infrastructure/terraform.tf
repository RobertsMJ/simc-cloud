resource "aws_apigatewayv2_api" "http_api" {
  name          = "simc-cloud-http-api"
  protocol_type = "HTTP"
}

resource "aws_lambda_function" "echo" {
  function_name = "echo"
  handler       = "main"
  runtime       = "go1.x"
  filename      = "${path.module}/../lambdas/echo/main.zip"
  source_code_hash = filebase64sha256("${path.module}/../lambdas/echo/main.zip")
  role          = aws_iam_role.lambda_exec.arn
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_apigatewayv2_integration" "echo_integration" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.echo.invoke_arn
  integration_method = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "echo_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "POST /api/echo"
  target    = "integrations/${aws_apigatewayv2_integration.echo_integration.id}"
}

resource "aws_lambda_permission" "apigw_echo" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.echo.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}