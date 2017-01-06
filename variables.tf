variable "access_key" {
  description = "The AWS access key"
  default     = "AKIAIDJCICP6O5PC3S7Q"
}

variable "secret_key" {
  description = "The AWS secret key"
  default     = "s9Jo6Tr7NdL1BrwJzKVddEah7jhKnR98qaMh32Ko"
}

variable "region" {
  description = "The AWS region"
  default     = "ap-southeast-1"
}

variable "ami" {
  type        = "map"
  description = "The AWS ami"

  default = {
    ap-southeast-1 = "ami-83bd63e0"
  }
}

variable "instance_type" {
  description = " The instance type"
  default     = "t2.micro"
}

variable "key_name" {
  description = " SSH key name"
  default     = "sureshpt-apse"
}

variable "instance_ips" {
  description = "The IPs to use for our instances"
  default     = ["10.0.1.20", "10.0.1.21"]
}

variable "owner_tag" {
  default = ["team1", "team2"]
}

variable "environment" {
  default = "development"
}
