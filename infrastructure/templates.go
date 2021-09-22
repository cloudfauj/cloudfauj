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
}

output "compute_subnets" {
  value = [aws_subnet.compute.id]
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

const appTfTpl = `data "terraform_remote_state" "env" {
  backend = "local"
  config = {
    path = "{{.env_tfstate_file}}"
  }
}

# Variables that need to be supplied during invokation
# Note that these have default empty values only to make TF destroy
# invokation easier.
variable "ingress_port" { default = 0 }
variable "cpu" { default = 256 }
variable "memory" { default = 512 }
variable "ecr_image" { default = "" }
variable "app_health_check_path" { default = "" }

locals {
  name = "{{.env_name}}-{{.app_name}}"
}

data "aws_region" "current" {}

# Security group
resource "aws_security_group" "main_app_sg" {
  name        = local.name
  description = "${local.name} application cluster traffic control"
  vpc_id      = data.terraform_remote_state.env.outputs.main_vpc_id
  tags        = local.common_tags

  ingress {
    description = "Application main ingress"
    from_port   = var.ingress_port
    to_port     = var.ingress_port
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

# Task Definition
resource "aws_ecs_task_definition" "main_app" {
  family                   = local.name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = data.terraform_remote_state.env.outputs.ecs_task_execution_role_arn
  cpu                      = var.cpu
  memory                   = var.memory

  container_definitions = jsonencode([
    {
      name  = "{{.app_name}}"
      image = var.ecr_image

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-create-group"  = "true"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-group"         = "{{.env_name}}"
          "awslogs-stream-prefix" = "{{.app_name}}"
        }
      }

      essential    = true
      portMappings = [{ containerPort = tonumber(var.ingress_port) }]
    }
  ])
}

# ECS Service
resource "aws_ecs_service" "main_app" {
  name                = "{{.app_name}}"
  cluster             = data.terraform_remote_state.env.outputs.compute_ecs_cluster_arn
  desired_count       = 1
  launch_type         = "FARGATE"
  task_definition     = aws_ecs_task_definition.main_app.arn
  scheduling_strategy = "REPLICA"

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  network_configuration {
    subnets          = data.terraform_remote_state.env.outputs.compute_subnets
    assign_public_ip = true
    security_groups  = [aws_security_group.main_app_sg.id]
  }

  // Only associate Load balancer if target group is supplied.
  dynamic "load_balancer" {
    for_each = [{{.target_group_resource}}]
    content {
      target_group_arn = load_balancer.value
      container_name   = "{{.app_name}}"
      container_port   = var.ingress_port
    }
  }
}

output "ecs_service" {
  value = aws_ecs_service.main_app.name
}

output "ecs_cluster_arn" {
  value = data.terraform_remote_state.env.outputs.compute_ecs_cluster_arn
}`

const appDnsTfTpl = `data "terraform_remote_state" "domain" {
  backend = "local"
  config = {
    path = "{{.domain_tfstate_file}}"
  }
}

locals {
  app_url = "${local.name}.{{.domain_name}}"
}

data "aws_lb" "apps_alb" {
  name = data.terraform_remote_state.env.outputs.apps_alb_name
}

resource "aws_route53_record" "app_url" {
  name    = local.app_url
  type    = "CNAME"
  zone_id = data.terraform_remote_state.domain.outputs.zone_id
  ttl     = 60
  records = [data.aws_lb.apps_alb.dns_name]
}

resource "aws_alb_target_group" "alb_to_ecs_service" {
  name        = replace(local.name, "_", "-")
  vpc_id      = data.terraform_remote_state.env.outputs.main_vpc_id
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  tags        = local.common_tags

  health_check {
    path = var.app_health_check_path
  }
}

resource "aws_lb_listener_rule" "app_router" {
  listener_arn = data.terraform_remote_state.env.outputs.main_alb_https_listener

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.alb_to_ecs_service.arn
  }
  condition {
    host_header {
      values = [local.app_url]
    }
  }
}`
