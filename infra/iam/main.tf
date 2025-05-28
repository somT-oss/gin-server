resource "aws_iam_role" "ec2_instance_role" {
  name = var.ec2_instance_role_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "sqs_send_policy" {
  name = var.sqs_iam_send_policy_name
  role = aws_iam_role.ec2_instance_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid = "AllowSqsSendMessage"
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage"
        ]
        Resource = var.sqs_send_policy_resource
      }
    ]
  })
}


resource "aws_iam_instance_profile" "leetcode_server_instance_profile" {
  name = var.leetcode_server_instance_profile_name
  role = aws_iam_role.ec2_instance_role.name
}
