package infrastructure

const tfCoreConfigTpl = `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "{{.aws_provider_version}}"
    }
  }
}

provider "aws" {
  region = "{{.aws_region}}"
}

locals {
  common_tags = {
    manager = "cloudfauj"
  }
}`

const domainDnsTfConfigTpl = `resource "aws_route53_zone" "dns_manager" {
  name          = "{{.domain_name}}"
  tags          = local.common_tags
  force_destroy = true
  comment       = "Public Hosted Zone for {{.domain_name}} managed by Cloudfauj"
}

// Add the DNS validation records to the Hosted zone.
// Note that this alone is not enough to validate the ACM cert.
// The NS records of the main R53 zone must be configured with the domain
// provider manually by the user.
// After that, ACM will be able to validate & issue the certificate.
resource "aws_route53_record" "acm_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.primary_wildcard_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true

  name    = each.value.name
  records = [each.value.record]
  ttl     = 60
  type    = each.value.type
  zone_id = aws_route53_zone.dns_manager.zone_id
}

output "zone_id" {
  value = aws_route53_zone.dns_manager.zone_id
}

output "name_servers" {
  value = aws_route53_zone.dns_manager.name_servers
}`

const domainCertTfConfigTpl = `resource "aws_acm_certificate" "primary_wildcard_cert" {
  domain_name               = "{{.domain_name}}"
  subject_alternative_names = ["*.{{.domain_name}}"]
  validation_method         = "DNS"
  tags                      = local.common_tags
}

output "ssl_cert_arn" {
  value = aws_acm_certificate.primary_wildcard_cert.arn
}

output "apex_domain" {
  value = "{{.domain_name}}"
}`

const envDomainStateTfTpl = `data "terraform_remote_state" "domain" {
  backend = "local"
  config = {
    path = "%s"
  }
}`

const envOrchestratorTfTpl = `# ECS Fargate cluster
resource "aws_ecs_cluster" "compute_cluster" {
  name               = "%s"
  tags               = local.common_tags
  capacity_providers = ["FARGATE"]
}

output "compute_ecs_cluster_arn" {
  value = aws_ecs_cluster.compute_cluster.arn
}`

const envNetworkTfTpl = `data "aws_availability_zones" "available" {}

# VPC
resource "aws_vpc" "main_vpc" {
  cidr_block = "{{.vpc_cidr}}"
  tags       = local.common_tags
}

# Internet gateway
resource "aws_internet_gateway" "main_vpc_igw" {
  vpc_id = aws_vpc.main_vpc.id
  tags   = local.common_tags
}

# Routing
resource "aws_default_route_table" "main_vpc_default_rt" {
  default_route_table_id = aws_vpc.main_vpc.default_route_table_id
  tags                   = local.common_tags

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main_vpc_igw.id
  }
}

# Subnets
resource "aws_subnet" "compute" {
  vpc_id     = aws_vpc.main_vpc.id
  cidr_block = cidrsubnet(aws_vpc.main_vpc.cidr_block, 4, 1)
  tags       = local.common_tags
}

resource "aws_subnet" "apps_alb" {
  count             = 2
  vpc_id            = aws_vpc.main_vpc.id
  cidr_block        = cidrsubnet(aws_vpc.main_vpc.cidr_block, 9, count.index + 1)
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name    = "{{.env_name}}-alb-${count.index}"
    manager = local.common_tags.manager
  }
}

# ECS Task IAM role shared by all applications in the environment
resource "aws_iam_role" "ecs_task_exec_role" {
  name               = "{{.env_name}}-ecs-task-exec-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_exec_role_assume_role.json
  tags               = local.common_tags
}

data "aws_iam_policy_document" "ecs_task_exec_role_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_exec_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = aws_iam_role.ecs_task_exec_role.name
}

resource "aws_iam_role_policy" "ecs_task_execution_role_custom_policy" {
  name   = "{{.env_name}}-ecs-task-exec-role-custom-policy"
  role   = aws_iam_role.ecs_task_exec_role.id
  policy = data.aws_iam_policy_document.ecs_task_exec_role_custom_policy.json
}

data "aws_iam_policy_document" "ecs_task_exec_role_custom_policy" {
  statement {
    effect    = "Allow"
    actions   = ["logs:CreateLogGroup"]
    resources = ["*"]
  }
}

output "ecs_task_execution_role_arn" {
  value = aws_iam_role.ecs_task_exec_role.arn
}

output "main_vpc_id" {
  value = aws_vpc.main_vpc.id
}`

const envAlbTfTpl = `resource "aws_security_group" "env_apps_alb" {
  name        = "{{.env_name}}-apps-alb"
  description = "{{.env_name}} applications ALB traffic control"
  vpc_id      = aws_vpc.main_vpc.id

  tags = {
    Name    = "{{.env_name}}-alb"
    manager = local.common_tags.manager
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_alb" "env_apps" {
  name               = "{{.env_name}}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.env_apps_alb.id]
  tags               = local.common_tags
  subnets            = aws_subnet.apps_alb.*.id
}

resource "aws_alb_listener" "env_apps_https" {
  load_balancer_arn = aws_alb.env_apps.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"
  certificate_arn   = data.terraform_remote_state.domain.outputs.ssl_cert_arn

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Service Unavailable"
      status_code  = "503"
    }
  }
}

output "apps_alb_arn" {
  value = aws_alb.env_apps.arn
}

output "apps_alb_name" {
  value = aws_alb.env_apps.name
}

output "main_alb_https_listener" {
  value = aws_alb_listener.env_apps_https.arn
}`
