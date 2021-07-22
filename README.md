# go-dappley-multi-server-testing
The go-dappley-multi-server-tester pipeline task tests all the command line interface commands built inside the "go-dappley" program by establishing a blockchain network through connecting multiple AWS EC2 together. Each AWS EC2 instance works as a node for the blockchain network. Since keeping multiple test servers running at all time is wasteful, the pipeline creates, initializes and terminates all the AWS EC2 instances that are being used in this pipeline.

### Initial Testing Pipeline
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

### Multi Server Testing Pipeline
```
pipeline {
    agent any
    tools {
        go 'go-1.16.6'
    }
    environment {
        GO1166MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/go-dappley-multi-server-tester/'
            }
        }
        stage('Create instances 1') {
            steps {
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode1}]' > node1.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode2}]' > node2.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode3}]' > node3.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode4}]' > node4.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node5.txt"
            }
        }
        stage('Build') {
            steps {
                sh 'go mod init github.com/heesooh/go-dappley-multi-server-tester'
                sh 'go mod tidy'
                sh 'go build'
            }
        }
        stage('Create Directories') {
            steps {
                sh 'mkdir test_results test_results/createAccount test_results/getBalance test_results/listAddresses test_results/send test_results/sendFromMiner test_results/setup test_results/smartContract test_results/getPeerInfo test_results/getBlockchain test_results/estimateGas test_results/addProducer'
                sh 'mkdir accounts accounts/node1 accounts/node2 accounts/node3 accounts/node4 accounts/node5'
                sh 'mkdir playbook_error playbook_error/node1 playbook_error/node2 playbook_error/node3 playbook_error/node4 playbook_error/node5 playbook_error/node6'
            }
        }
        stage('Update Hosts 1') {
            steps {
                sh './go-dappley-multi-server-tester -function update -number 5'
            }
        }
        stage('Initialize Hosts 1') {
            steps {
                sh './go-dappley-multi-server-tester -function initialize -number 5'
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
                // sh 'ansible-playbook ./playbooks/setup/get_PeerID.yml > ./test_results/setup/get_PeerID.txt'
                run_playbooks_5('./playbooks/setup/get_PeerID.yml')
            }
        }
        stage('Update Playbooks') {
            steps {
                sh './go-dappley-multi-server-tester -function update_address'
            }
        }
        stage('Get Peer Info'){
            steps{
                run_playbooks_5('./playbooks/getPeerInfo/check_PeerID.yml')
                // > ./test_results/getPeerInfo/check_PeerID.txt'
            }
        }
        stage('Send Playbooks'){
            steps{
                run_playbooks_5('./playbooks/send/wrong_node.yml')
                // > ./test_results/send/wrong_node.txt'

                multi_playbooks('./playbooks/send/invalid_tip/invalid_tip_', 5)
                // > ./test_results/send/invalid_tip.txt'

                multi_playbooks('./playbooks/send/missing_flag/missing_flag_', 3)
                // > ./test_results/send/missing_flag.txt'

                multi_playbooks('./playbooks/send/invalid_file/invalid_file_', 2)
                // > ./test_results/send/invalid_file.txt'

                multi_playbooks('./playbooks/send/invalid_data/invalid_data_', 13)
                // > ./test_results/send/invalid_data.txt'

                multi_playbooks('./playbooks/send/invalid_amount/invalid_amount_', 6)
                // > ./test_results/send/invalid_amount.txt'

                multi_playbooks('./playbooks/send/invalid_address/invalid_address_', 6)
                // > ./test_results/send/invalid_address.txt'

                multi_playbooks('./playbooks/send/invalid_gas_limit/invalid_gas_limit_', 7)
                // > ./test_results/send/invalid_gas_limit.txt'

                multi_playbooks('./playbooks/send/invalid_gas_price/invalid_gas_price_', 6)
                // > ./test_results/send/invalid_gas_price.txt'

                run_playbooks_5('./playbooks/send/single_transaction_no_tip.yml')
                // > ./test_results/send/single_transaction_no_tip.txt'

                run_playbooks_5('./playbooks/send/single_transaction_with_tip.yml')
                // > ./test_results/send/single_transaction_with_tip.txt'

                run_playbooks_5('./playbooks/send/multi_transaction_no_tip.yml')
                // > ./test_results/send/multi_transaction_no_tip.txt'

                run_playbooks_5('./playbooks/send/multi_transaction_with_tip.yml')
                // > ./test_results/send/multi_transaction_with_tip.txt'
            }
        }
        stage('GetBlockchain Playbooks') {
            steps {
                run_playbooks_5('./playbooks/getBlockchain/getBlockchainInfo.yml')
                // > ./test_results/getBlockchain/getBlockchainInfo.txt'

                multi_playbooks('./playbooks/getBlockchain/invalid_input/invalid_input_', 12)
                // > ./test_results/getBlockchain/invalid_input.txt'
            }
        }
        stage('EstimateGas Playbooks') {
            steps {
                run_playbooks_5('./playbooks/estimateGas/estimate_gas.yml')
                // > ./test_results/estimateGas/estimate_gas.yml'

                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_address.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_smartContract.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_amount.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_tip.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_gasLimit.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_gasPrice.yml'
                // sh 'ansible-playbook ./playbooks/estimateGas/invalid_data.yml'
            }
        }
        stage('GetBalance Playbooks') {
            steps {
                run_playbooks_5('./playbooks/getBalance/invalid_address.yml')
                // > ./test_results/getBalance/invalid_address.txt'

                run_playbooks_5('./playbooks/getBalance/missing_argument.yml')
                // > ./test_results/getBalance/missing_argument.txt'
            }
        }
        stage('CreateAccount Playbooks') {
            steps {
                run_playbooks_5('./playbooks/createAccount/empty_password.yml')
                // > ./test_results/createAccount/empty_password.txt'

                run_playbooks_5('./playbooks/createAccount/invalid_password.yml')
                // > ./test_results/createAccount/invalid_password.txt'
            }
        }
        stage('ListAddresses Playbooks') {
            steps {
                run_playbooks_5('./playbooks/listAddresses/invalid_password.yml')
                // > ./test_results/listAddresses/invalid_password.txt'
            }
        }
        stage('SendFromMiner Playbooks') {
            steps {
                run_playbooks_5('./playbooks/sendFromMiner/missing_flag.yml')
                // > ./test_results/sendFromMiner/missing_flag.txt'

                multi_playbooks('./playbooks/sendFromMiner/invalid_amount/invalid_amount_', 4)
                // > ./test_results/sendFromMiner/invalid_amount.txt'

                multi_playbooks('./playbooks/sendFromMiner/invalid_address/invalid_address_', 3)
                // > ./test_results/sendFromMiner/invalid_address.txt'

                run_playbooks_5('./playbooks/sendFromMiner/single_transaction_from_miner.yml')
                // > ./test_results/sendFromMiner/single_transaction_from_miner.txt'
            }
        }
        stage('Smart Contract Playbooks') {
            steps {
                run_playbooks_5('./playbooks/smartContract/smart_contract_gas_1.yml')
                // > ./test_results/smartContract/smart_contract_gas_1.txt'

                run_playbooks_5('./playbooks/smartContract/smart_contract_gas_2.yml')
                // > ./test_results/smartContract/smart_contract_gas_2.txt'
            }
        }
        stage('Add Producer Playbooks') {
            steps {
                run_playbooks_5('./playbooks/addProducer/add_producer.yml')
                // > ./test_results/addProducer/add_producer.txt'

                run_playbooks_5('./playbooks/addProducer/add_when_max.yml')
                // > ./test_results/addProducer/add_when_max.txt'

                run_playbooks_5('./playbooks/addProducer/invalid_address.yml')
                // > ./test_results/addProducer/invalid_address.txt'
            }
        }
        stage('Terminate Host Nodes 1') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
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
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode1}]' > node1.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode2}]' > node2.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode3}]' > node3.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode4}]' > node4.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node5.txt"
                sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode5}]' > node6.txt"
            }
        }
        stage('Update Hosts 2') {
            steps {
                sh './go-dappley-multi-server-tester -function update -number 6'
            }
        }
        stage('Initialize Hosts 2') {
            steps {
                sh './go-dappley-multi-server-tester -function initialize -number 6'
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
                run_playbooks_6('./playbooks/addProducer/6_producer_test_1.yml')
                // > ./test_results/addProducer/6_producer_test_1.txt'

                sh './go-dappley-multi-server-tester -function terminate -number 2'
                sh 'ansible-playbook ./playbooks/addProducer/6_producer_test_2.yml  > ./test_results/addProducer/6_producer_test_2.txt'
            }
        }
        // stage('Send Report') {
        //     steps {
        //         sh 'zip -r test_results.zip test_results'
        //         sh './go-dappley-multi-server-tester -function send_result -senderEmail <Sender Email> -senderPasswd <Sender Email Password>'
        //     }
        // }
        stage('Terminate Host Nodes 2') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
                sh 'rm -r *'
            }
        }
    }
}

def multi_playbooks(String directory, int number) {
    for (int i = number; i > 0; i--) {
        run_playbooks_5("${directory}${i}.yml")
    }
}

def run_playbooks_5(String directory) {
    sh "ansible-playbook ${directory}"
    script {
        def ERROR1 = readFile './playbook_error/node1/error_1.txt'
        def ERROR2 = readFile './playbook_error/node2/error_2.txt'
        def ERROR3 = readFile './playbook_error/node3/error_3.txt'
        def ERROR4 = readFile './playbook_error/node4/error_4.txt'
        def ERROR5 = readFile './playbook_error/node5/error_5.txt'
        while (ERROR1.trim().equals("true") || ERROR2.trim().equals("true") || ERROR3.trim().equals("true") || ERROR4.trim().equals("true") || ERROR5.trim().equals("true")) {
            echo "Blockchain initialization failed!"
            sh 'ansible-playbook ./playbooks/clear.yml'
            sh 'sleep 30'
            sh "ansible-playbook ${directory}"
            // > ./test_results/getPeerInfo/check_PeerID.txt'
            ERROR1 = readFile './playbook_error/node1/error_1.txt'
            ERROR2 = readFile './playbook_error/node2/error_2.txt'
            ERROR3 = readFile './playbook_error/node3/error_3.txt'
            ERROR4 = readFile './playbook_error/node4/error_4.txt'
            ERROR5 = readFile './playbook_error/node5/error_5.txt'
        }
        echo "Blockchain initialized successfully."
    }
}

def run_playbooks_6(String directory) {
    sh "ansible-playbook ${directory}"
    script {
        def ERROR1 = readFile './playbook_error/node1/error_1.txt'
        def ERROR2 = readFile './playbook_error/node2/error_2.txt'
        def ERROR3 = readFile './playbook_error/node3/error_3.txt'
        def ERROR4 = readFile './playbook_error/node4/error_4.txt'
        def ERROR5 = readFile './playbook_error/node5/error_5.txt'
        def ERROR6 = readFile './playbook_error/node6/error_6.txt'
        while (ERROR1.trim().equals("true") || ERROR2.trim().equals("true") || ERROR3.trim().equals("true") || ERROR4.trim().equals("true") || ERROR5.trim().equals("true") || ERROR6.trim().equals("true")) {
            echo "Blockchain initialization failed!"
            sh 'ansible-playbook ./playbooks/clear.yml'
            sh 'sleep 30'
            sh "ansible-playbook ${directory}"
            // > ./test_results/getPeerInfo/check_PeerID.txt'
            ERROR1 = readFile './playbook_error/node1/error_1.txt'
            ERROR2 = readFile './playbook_error/node2/error_2.txt'
            ERROR3 = readFile './playbook_error/node3/error_3.txt'
            ERROR4 = readFile './playbook_error/node4/error_4.txt'
            ERROR5 = readFile './playbook_error/node5/error_5.txt'
            ERROR6 = readFile './playbook_error/node6/error_6.txt'
        }
        echo "Blockchain initialized successfully."
    }
}
```