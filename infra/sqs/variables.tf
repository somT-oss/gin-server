variable "leetcode_queue" {
    type = string 
}

variable "delay_seconds" {
    type = number
}

variable "max_message_size" {
    type = number
}

variable "message_retention_seconds" {
    type = number
}

variable "receive_wait_time_seconds" {
    type = number
}

variable "leetcode_dead_letter_queue_name" {
    type = string
}
