variable "vpc_cidr_block_ip" {
    type = string
    default = "10.0.0.0/16"  
}

variable "subnet_1_cidr_block_ip" {
    type = string 
    default = "10.0.1.0/24"
}

variable "subnet_2_cidr_block_ip" {
    type = string
    default = "10.0.2.0/24"
}

variable "availability_zone_1" {
    type = string 
    default = "us-east-1a"
}

variable "availability_zone_2" {
    type = string
    default = "us-east-1b"
}

variable "instance_ami" {
    type = string
    default = "ami-0c614dee691cbbf37"
}

variable "instance_type" {
    type = string
    default = "t2.micro"
}

variable "sg_name" {
    type = string
    default = "leetcode_server_sg"
}

variable "leetcode_dead_letter_queue_name" {
    type = string
    default = "leetcode_dead_letter_queue"
}

variable "max_message_size" {
    type = number
    default = 2048
}

variable "leetcode_queue_name" {
    type = string
    default = "leetcode_queue"
}

variable "delay_seconds" {
    type = number
    default = 90
}

variable "message_retention_seconds" {
    type = number
    default = 86400
}

variable "receive_wait_time_seconds" {
    type = number
    default = 10
}

variable "ec2_instance_role_name" {
    type = string
    default = "LeetcodeEC2SendMessageRole"
}

variable "sqs_iam_send_policy_name" {
    type = string
    default = "LeetcodeSQSQueueSendMessagePolicy"
}

variable "leetcode_server_instance_profile_name" {
    type = string
    default = "leetcode_server_instance_profile"
}