output "server_vpc_id" {
    value = aws_vpc.leetcode_server_vpc.id
}

output "leetcode_server_id" {
    value = aws_instance.leetcode_server.id
}

output "leetcode_server_ip" {
    value = aws_instance.leetcode_server.public_ip
}
output "leetcode_server_sg_id" {
    value = aws_security_group.leetcode_server_sg.id
}

output "leetcode_server_subnet_1_id" {
    value = aws_subnet.leetcode_server_subnet_1.id
}

output "leetcode_server_subnet_2_id" {
    value = aws_subnet.leetcode_server_subnet_2.id
}
