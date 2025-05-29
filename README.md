# Leetcode Two Sum Solver

This is a project that replicates the leetcode style code execution engine on a small scale to solve the popular leetcode [two sum problem](https://leetcode.com/problems/two-sum/). 

Leetcode code execution engine is an impressive system that parses user code submission and returns an appropriate response after testing the users submission against a list of problems. Some probelms require you to write optimal code, this code engine also tests optimal code, catching inefficiencies in Time Limit Exceeded (TLE) and Memory Limit Exceeded (MLE) errors.


## Architecture
The architecture of this application involves two main components, a server and a queue. An endpoint on the server is used to send source code and the language type to a queue service, in this case Amazon Simple Queue Service (SQS). The queue serves a buffer to hold code submissions. Another endpoint is used to read the message from the queue and process the code in the same server within a docker container.


## Prerequisites
- An AWS account

- Terraform

- Golang version >= 1.23

- Ansible
 
## Local Testing
1 Clone the repository.

2 Run ```go mod download``` to install go dependencies.

3 Create an SQS queue in AWS.

4 Create a .env and populate it with these fields.
```
    SQS_URL=https://sqs.us-east-1.amazonaws.com
    ACCOUNT_ID=<aws_account_id>
    QUEUE_NAME=<sqs_queue_name>
```

5 Run ```go run main.go``` to start up the local development server.

To send a message to the queue, MAKE SURE to send a POST request to the ```/send/``` endpoint
### JSON request body example
```
{
    "source_code": "def two_sum(nums, target):\n    lookup = {}\n    for i, num in enumerate(nums):\n        if target - num in lookup:\n            return [lookup[target - num], i]\n        lookup[num] = i\n    return []",
    "language": "python"
}
```
**PS**: You can convert the code in source code to actual code and test it out in your code editor.


## Deploy to AWS
1 Clone the repository

2 cd into ```infra``` folder

3 Run ```terraform init```

4 Configure AWS on your PC with the ```aws configure``` command

5 Run ```terraform apply```

6 Run ```terraform ouput```. This command returns the ip address of the newly deployed server

7 Update the ```your_inventory.ini``` file ```ansible/inventory/``` with this code
```
<ec2-user@ip_address>
``` 

6 Run ```ansible-playbook ansible/leetcode_server_playbook.yaml -i ansible/inventory/your_inventory.ini --private-key leetcode_server.pem``` to install dependencies in the server

7 ssh into the server with ```ssh ec2-user@ip_address -i leetcode_server.pem```.

8 cd into ```leetcode-two-sum-server```

9 Create a .env and populate it with these fields.
```
    SQS_URL=https://sqs.us-east-1.amazonaws.com
    ACCOUNT_ID=<aws_account_id>
    QUEUE_NAME=leetcode_queue
```

10 Run ```sudo chown -R ec2-user:ec2-user ~```. This changes the ownership of the directory to the current user in the server. This will help to create the files we need when we want to test the code.

11 Run the server with "./main"

12 Access the application on ```<ip_address>:8080/```

**PS**: If you mess up the setup after step 6, make sure to delete the pulled code directory in the server and rerun step 6.


# Contributions
This project is open to contributions, there are some areas that need fine tuning. I will leave issues open for contributions.

Thank you.