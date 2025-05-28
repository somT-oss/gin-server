resource "aws_sqs_queue" "leetcode_queue" {
  name                      = var.leetcode_queue
  delay_seconds             = var.delay_seconds
  max_message_size          = var.max_message_size
  message_retention_seconds = var.message_retention_seconds 
  receive_wait_time_seconds = var.receive_wait_time_seconds
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.leetcode_dead_letter_queue.arn
    maxReceiveCount     = 4
  })

  tags = {
    Environment = "dev"
  }
}

resource "aws_sqs_queue" "leetcode_dead_letter_queue" {
  name = var.leetcode_dead_letter_queue_name
}