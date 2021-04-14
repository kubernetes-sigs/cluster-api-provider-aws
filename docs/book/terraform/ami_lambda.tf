provider "aws" {
  region = "us-east-1"
  profile = "heptio-oss"
}


resource "aws_lambda_function" "amilist" {
  function_name    = "amilist"
  filename         = "amilist.zip"
  handler          = "amilist"
  source_code_hash = "data.archive_file.zip.output_base64sha256"
  role             = aws_iam_role.iam_for_lambda.arn
  runtime          = "go1.x"
  memory_size      = 256
  timeout          = 60
}

resource "aws_s3_bucket" "cluster-api-aws" {
  bucket = "cluster-api-aws-amis.sigs.k8s.io"
  acl    = "public-read"
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
  website {
    index_document = "amis.json"
  }
}

resource "aws_cloudfront_origin_access_identity" "cluster_api" {
  comment = "Cluster API AMIs"
}

resource "aws_cloudfront_distribution" "cluster-api-aws" {
  origin {
    origin_id = "s3"
    domain_name = aws_s3_bucket.cluster-api-aws.bucket_regional_domain_name
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.cluster_api.cloudfront_access_identity_path
    }
  }
  enabled = true
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  default_cache_behavior {
    allowed_methods  = ["HEAD","GET"]
    cached_methods   = ["HEAD","GET"]
    target_origin_id = "s3"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }

      headers      = ["Origin", "Access-Control-Request-Headers", "Access-Control-Request-Method"]

    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 60
    default_ttl            = 3600
    max_ttl                = 3600
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "amilist-lambda"

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })

  inline_policy {
    name = "inline-policy"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Action = [
            "ec2:DescribeImages",
          ]
          Resource = "*"
          Effect = "Allow"
        },
        {
          Action = [
            "s3:*",
          ]
          Resource = "arn:aws:s3:::cluster-api-aws-amis.sigs.k8s.io/*"
          Effect = "Allow"
        },
      ]
    })
  }
}

resource "aws_s3_bucket_policy" "cluster-api" {
  bucket = aws_s3_bucket.cluster-api-aws.id
  policy = data.aws_iam_policy_document.s3_policy.json
}

data "aws_iam_policy_document" "s3_policy" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.cluster-api-aws.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.cluster_api.iam_arn]
    }
  }
}

resource "aws_cloudwatch_event_rule" "every_hour" {
    name = "every-hour"
    description = "Fires every hour"
    schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "update_amis_every_hour" {
    rule = aws_cloudwatch_event_rule.every_hour.name
    target_id = "amilist"
    arn = aws_lambda_function.amilist.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_amilist" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.amilist.function_name
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.every_hour.arn
}

data "archive_file" "zip" {
  type        = "zip"
  source_dir = "out"
  output_path = "amilist.zip"
}
