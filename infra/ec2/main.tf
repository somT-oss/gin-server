resource "tls_private_key" "ssh_key" {
    algorithm = "RSA"
    rsa_bits = 4096
}


resource "aws_key_pair" "ssh_key_pair" {
    key_name = "leetcode_server_key"
    public_key = tls_private_key.ssh_key.public_key_openssh

    provisioner "local-exec" {
        command = "echo '${tls_private_key.ssh_key.private_key_pem}' > ./leetcode_server.pem"
    }
}


resource "aws_vpc" "leetcode_server_vpc" {
    cidr_block = var.vpc_cidr_block_ip

    tags = {
        Name = "leetcode_server_vpc"
    }
}

resource "aws_subnet" "leetcode_server_subnet_1" {
    vpc_id = aws_vpc.leetcode_server_vpc.id
    cidr_block = var.subnet_1_cidr_block_ip
    availability_zone = var.availability_zone_1
    map_public_ip_on_launch = true
    

    tags = {
        Name = "leetcode_server_subnet_1"
    }
    depends_on = [aws_vpc.leetcode_server_vpc]
}


resource "aws_subnet" "leetcode_server_subnet_2" {
    vpc_id = aws_vpc.leetcode_server_vpc.id
    cidr_block = var.subnet_2_cidr_block_ip
    availability_zone = var.availability_zone_2
    map_public_ip_on_launch = true
    

    tags = {
        Name = "leetcode_server_subnet_2"
    }
    depends_on = [aws_vpc.leetcode_server_vpc]
}


resource "aws_internet_gateway" "leetcode_server_ig" {
    vpc_id = aws_vpc.leetcode_server_vpc.id

    tags = {
        Name = "leetcode_server_ig"
    }
}

resource "aws_route_table" "leetcode_server_rt" {
    vpc_id = aws_vpc.leetcode_server_vpc.id

    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.leetcode_server_ig.id
    }

    tags = {
        Name = "leetcode_server_rt"
    }
}

resource "aws_route_table_association" "leetcode_server_rta" {
    subnet_id = aws_subnet.leetcode_server_subnet_1.id
    route_table_id = aws_route_table.leetcode_server_rt.id
}

resource "aws_security_group" "leetcode_server_sg" {
    name = var.sg_name
    vpc_id = aws_vpc.leetcode_server_vpc.id

    # Allow HTTP connections to port 80
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    # Allow SSH connection to port 22
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    # Allow all outbound traffic
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1" # All protocols
        cidr_blocks =  ["0.0.0.0/0"]
    }

    tags = {
      Name = "leetcode_server_sg"
    }
    depends_on = [aws_vpc.leetcode_server_vpc]
}


resource "aws_instance" "leetcode_server" {
    ami = var.instance_ami
    instance_type = var.instance_type
    associate_public_ip_address = true
    subnet_id = aws_subnet.leetcode_server_subnet_1.id
    vpc_security_group_ids = [aws_security_group.leetcode_server_sg.id]
    key_name = aws_key_pair.ssh_key_pair.key_name

    credit_specification {
        cpu_credits = "standard"
    }

    tags = {
      Name = "leetcode_server1"
    }

    user_data = <<-EOF
    #!/bin/bash
    sudo yum update -y
    sudo yum install -y nginx
    sudo systemctl enable nginx
    sudo systemctl start nginx
    EOF


    depends_on = [ aws_security_group.leetcode_server_sg ]
}