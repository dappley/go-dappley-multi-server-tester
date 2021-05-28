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

Ansible Playbook Pipeline
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
        stage('Build') {
            steps {
                sh 'go mod init github.com/heesooh/ansible-dappley'
                sh 'go mod tidy'
                sh 'go build'
            }
        }
        stage('Update Hosts') {
            steps {
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
                sh 'ansible-playbook setup.yml > setup.txt'
            }
        }
        stage('Send From Miner') {
            steps {
                sh 'ansible-playbook sendFromMiner.yml > sendFromMiner.txt'
            }
        }
        stage('Single Transaction No Tip') {
            steps {
                sh 'ansible-playbook single_transaction_no_tip.yml > single_transaction_no_tip.txt'
            }
        }
        stage('Single Transaction With Tip') {
            steps {
                sh 'ansible-playbook single_transaction_with_tip.yml > single_transaction_with_tip.txt'
            }
        }
        stage('Multi Transaction No Tip') {
            steps {
                sh 'ansible-playbook multi_transaction_no_tip.yml > multi_transaction_no_tip.txt'
            }
        }
        stage('Multi Transaction With Tip') {
            steps {
                sh 'ansible-playbook multi_transaction_with_tip.yml > multi_transaction_with_tip.txt'
            }
        }
        stage('Smart Contract Gas 1') {
            steps {
                sh 'ansible-playbook smart_contract_gas_1.yml > smart_contract_gas_1.txt'
            }
        }
        stage('Smart Contract Gas 2') {
            steps {
                sh 'ansible-playbook smart_contract_gas_2.yml > smart_contract_gas_2.txt'
            }
        }
        stage('Test Send') {
            steps {
                sh 'ansible-playbook send.yml > seend.txt'
            }
        }
        stage('Send Report') {
            steps {
                sh './ansible-dappley -function send_result -recipient blockchainwarning@omnisolu.com -senderEmail blockchainwarning@omnisolu.com -senderPasswd gabroq-bucfe0-pubqiC'
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