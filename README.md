# ansible-dappley
ansible playbook

Initial Testing Pipeline
```
pipeline {
    agent any
    tools {
        go 'go-1.16.3'
    }
    environment {
        GO1163MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/ansible-dappley/'
            }
        }
    }
}
```

Blockchain Testing Pipeline
```
pipeline {
    agent any
    tools {
        go 'go-1.16.3'
    }
    environment {
        GO1163MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/ansible-dappley/'
            }
        }
        stage('Create Nodes') {
            steps {
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.xlarge --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode1}]' > node1.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.xlarge --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode2}]' > node2.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.xlarge --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode3}]' > node3.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.xlarge --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode4}]' > node4.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.xlarge --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node5.txt"
            }
        }
        stage('Update Hosts') {
            steps {
                sh 'go build ansible-dappley.go'
                sh './ansible-dappley -function update'
            }
        }
        stage('Initialize Hosts') {
            steps {
                sh './ansible-dappley -function initialize'
            }
        }
        stage('Setup Host Nodes') {
            steps {
                sh 'ansible-playbook setup.yml'
            }
        }
        stage('Run Test') {
            steps {
                sh 'ansible-playbook test.yml'
            }
        }
        stage('Terminate Host Nodes') {
            steps {
                sh './ansible-dappley -function terminate'
            }
        }
        stage('Close') {
            steps {
                sh 'rm -r *'
            }
        }
    }
}
```