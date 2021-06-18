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
        stage('Create Directories') {
            steps {
                sh 'mkdir test_results test_results/createAccount test_results/getBalance test_results/listAddresses test_results/send test_results/sendFromMiner test_results/setup test_results/smartContract'
                sh 'mkdir /var/lib/jenkins/workspace/go-dappley-ansible-accounts /var/lib/jenkins/workspace/go-dappley-ansible-accounts/node1 /var/lib/jenkins/workspace/go-dappley-ansible-accounts/node2 /var/lib/jenkins/workspace/go-dappley-ansible-accounts/node3 /var/lib/jenkins/workspace/go-dappley-ansible-accounts/node4 /var/lib/jenkins/workspace/go-dappley-ansible-accounts/node5'
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
                sh 'ansible-playbook ./playbooks/setup/setup.yml > ./test_results/setup/setup.txt'
            }
        }
        stage('Generate Accoounts') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/accounts_generator.yml > ./test_results/setup/accounts_generator.txt'
            }
        }
        stage('Update Seeds and Ports') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/update_seed_port.yml > ./test_results/setup/update_seed_port.txt'
            }
        }
        stage('Update Playbooks') {
            steps {
                sh './ansible-dappley -function update_address'
            }
        }
        stage('Send Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/send/wrong_node.yml > ./test_results/send/wrong_node.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_tip.yml > ./test_results/send/invalid_tip.txt'
                sh 'ansible-playbook ./playbooks/send/missing_flag.yml > ./test_results/send/missing_flag.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_file.yml > ./test_results/send/invalid_file.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_data.yml > ./test_results/send/invalid_data.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_amount.yml > ./test_results/send/invalid_amount.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_address.yml > ./test_results/send/invalid_address.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_gas_limit.yml > ./test_results/send/invalid_gas_limit.txt'
                sh 'ansible-playbook ./playbooks/send/invalid_gas_price.yml > ./test_results/send/invalid_gas_price.txt'
                sh 'ansible-playbook ./playbooks/send/single_transaction_no_tip.yml > ./test_results/send/single_transaction_no_tip.txt'
                sh 'ansible-playbook ./playbooks/send/single_transaction_with_tip.yml > ./test_results/send/single_transaction_with_tip.txt'
                sh 'ansible-playbook ./playbooks/send/multi_transaction_no_tip.yml > ./test_results/send/multi_transaction_no_tip.txt'
                sh 'ansible-playbook ./playbooks/send/multi_transaction_with_tip.yml > ./test_results/send/multi_transaction_with_tip.txt'
            }
        }
        stage('GetBalance Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/getBalance/invalid_address.yml > ./test_results/getBalance/invalid_address.txt'
                sh 'ansible-playbook ./playbooks/getBalance/missing_argument.yml > ./test_results/getBalance/missing_argument.txt'
            }
        }
        stage('CreateAccount Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/createAccount/empty_password.yml > ./test_results/createAccount/empty_password.txt'
                sh 'ansible-playbook ./playbooks/createAccount/invalid_password.yml > ./test_results/createAccount/invalid_password.txt'

            }
        }
        stage('ListAddresses Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/listAddresses/invalid_password.yml > ./test_results/listAddresses/invalid_password.txt'
            }
        }
        stage('SendFromMiner Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/sendFromMiner/missing_flag.yml > ./test_results/sendFromMiner/missing_flag.txt'
                sh 'ansible-playbook ./playbooks/sendFromMiner/invalid_amount.yml > ./test_results/sendFromMiner/invalid_amount.txt'
                sh 'ansible-playbook ./playbooks/sendFromMiner/invalid_address.yml > ./test_results/sendFromMiner/invalid_address.txt'
                sh 'ansible-playbook ./playbooks/sendFromMiner/single_transaction_from_miner.yml > ./test_results/sendFromMiner/single_transaction_from_miner.txt'
            }
        }

        stage('Smart Contract Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/smartContract/smart_contract_gas_1.yml > ./test_results/smartContract/smart_contract_gas_1.txt'
                sh 'ansible-playbook ./playbooks/smartContract/smart_contract_gas_2.yml > ./test_results/smartContract/smart_contract_gas_2.txt'
            }
        }
        stage('Send Report') {
            steps {
                sh 'zip -r test_results.zip test_results'
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
                sh 'rm -r ../go-dappley-ansible-accounts'   
            }
        }
    }
}
```