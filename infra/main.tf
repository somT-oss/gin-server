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
} 
