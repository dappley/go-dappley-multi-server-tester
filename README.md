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
        stage('Create instances 1') {
            steps {
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode1}]' > node1.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode2}]' > node2.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode3}]' > node3.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode4}]' > node4.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node5.txt"
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
                sh 'mkdir test_results test_results/createAccount test_results/getBalance test_results/listAddresses test_results/send test_results/sendFromMiner test_results/setup test_results/smartContract test_results/getPeerInfo test_results/getBlockchain test_results/estimateGas test_results/addProducer'
                sh 'mkdir accounts accounts/node1 accounts/node2 accounts/node3 accounts/node4 accounts/node5'
            }
        }
        stage('Update Hosts') {
            steps {
                sh './ansible-dappley -function update -number 5'
            }
        }
        stage('Initialize Hosts') {
            steps {
                sh './ansible-dappley -function initialize -number 5'
            }
        }
        stage('Setup Host Nodes 1') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/setup.yml > ./test_results/setup/setup_1.txt'
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
        stage('Generate Peer ID') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/get_PeerID.yml > ./test_results/setup/get_PeerID.txt'
            }
        }
        stage('Update Playbooks') {
            steps {
                sh './ansible-dappley -function update_address'
            }
        }
        stage('Get Peer Info'){
            steps{
                sh 'ansible-playbook ./playbooks/getPeerInfo/check_PeerID.yml > ./test_results/check_PeerID.txt'
            }
        }
        stage('Send Playbooks'){
            steps{
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
        stage('GetBlockchain Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/getBlockchain/getBlockchainInfo.yml > ./test_results/getBlockchain/getBlockchainInfo.txt'
                sh 'ansible-playbook ./playbooks/getBlockchain/invalid_input.yml > ./test_results/getBlockchain/invalid_input.txt'
            }
        }
        stage('EstimateGas Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/estimateGas/estimate_gas.yml > ./test_results/estimateGas/estimate_gas.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_address.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_smartContract.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_amount.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_tip.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_gasLimit.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_gasPrice.yml'
                sh 'ansible-playbook ./playbooks/estimateGas/invalid_data.yml'
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
        stage('Add Producer Playbooks') {
            steps {
                sh 'ansible-playbook ./playbooks/addProducer/add_producer.yml > ./test_results/addProducer/add_producer.txt'
                sh 'ansible-playbook ./playbooks/addProducer/add_when_max.yml > ./test_results/addProducer/add_when_max.txt'
                sh 'ansible-playbook ./playbooks/addProducer/invalid_address.yml > ./test_results/addProducer/invalid_address.txt'
            }
        }
        stage('Terminate Host Nodes') {
            steps {
                sh './ansible-dappley -function terminate'
            }
        }
        stage('Remove Files') {
            steps {
                sh 'rm -r node1.txt'
                sh 'rm -r node2.txt'
                sh 'rm -r node3.txt'
                sh 'rm -r node4.txt'
                sh 'rm -r node5.txt'
                sh 'rm -r instance_ids'
                sh 'rm -r hosts'
            }
        }
        stage('Create instances 2') {
            steps {
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode1}]' > node1.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode2}]' > node2.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode3}]' > node3.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode4}]' > node4.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node5.txt"
                sh "aws ec2 run-instances --image-id ami-02701bcdc5509e57b --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-03d39dd5dc5ddeef7 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode6}]' > node6.txt"
            }
        }
        stage('Update Hosts') {
            steps {
                sh './ansible-dappley -function update -number 6'
            }
        }
        stage('Initialize Hosts') {
            steps {
                sh './ansible-dappley -function initialize -number 6'
            }
        }
        stage('Setup Host Nodes 2-1') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/setup.yml > ./test_results/setup/setup_2.txt'
            }
        }
        stage('Setup Host Nodes 2-2') {
            steps {
                sh 'ansible-playbook ./playbooks/addProducer/setup.yml > ./test_results/setup/setup_3.txt'
            }
        }
        stage('6 Producer Test') {
            steps {
                sh 'ansible-playbook ./playbooks/addProducer/6_producer_test_1.yml'
                sh './ansible-dappley -function terminate -number 2'
                sh 'ansible-playbook ./playbooks/addProducer/6_producer_test_2.yml'
            }
        }
        stage('Send Report') {
            steps {
                sh 'zip -r test_results.zip test_results'
                sh './ansible-dappley -function send_result -recipient <RECIPIENT> -senderEmail <SENDER EMAIL> -senderPasswd <PASSWORD>'
            }
        }
        stage('Terminate Host Nodes') {
            steps {
                sh './ansible-dappley -function terminate'
                sh 'rm -r *'
            }
        }
    }
}
```