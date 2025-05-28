provider "aws" {
    region = "us-east-1"
}

module "ec2" {
    source = "./ec2"
    vpc_cidr_block_ip = var.vpc_cidr_block_ip
    subnet_1_cidr_block_ip = var.subnet_1_cidr_block_ip
    subnet_2_cidr_block_ip = var.subnet_2_cidr_block_ip
    availability_zone_1= var.availability_zone_1
    availability_zone_2 = var.availability_zone_2
    instance_ami = var.instance_ami
    instance_type = var.instance_type
    sg_name = var.sg_name
    leetcode_server_instance_profile = module.iam.leetcode_instance_profile_name
} 

module "sqs" {
    source = "./sqs"
    leetcode_dead_letter_queue_name =  var.leetcode_dead_letter_queue_name
    max_message_size = var.max_message_size
    leetcode_queue = var.leetcode_queue_name
    delay_seconds = var.delay_seconds
    message_retention_seconds = var.message_retention_seconds
    receive_wait_time_seconds = var.receive_wait_time_seconds
}

module "iam" {
    source = "./iam"
    ec2_instance_role_name =  var.ec2_instance_role_name
    sqs_iam_send_policy_name = var.sqs_iam_send_policy_name
    leetcode_server_instance_profile_name = var.leetcode_server_instance_profile_name
    sqs_send_policy_resource = module.sqs.leetcode_queue_arn
}