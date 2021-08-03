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
        stage('Build') {
            steps {
                sh 'go mod init github.com/heesooh/go-dappley-multi-server-tester'
                sh 'go mod tidy'
                sh 'go build'
            }
        }
        stage('Create Directories') {
            steps {
                sh 'mkdir test_results test_results/createAccount test_results/getBalance test_results/listAddresses test_results/send test_results/sendFromMiner test_results/setup test_results/smartContract test_results/getPeerInfo test_results/getBlockchain test_results/estimateGas test_results/addProducer test_results/removeProducer test_results/changeProducer'
                sh 'mkdir accounts accounts/node1 accounts/node2 accounts/node3 accounts/node4 accounts/node5'
                sh 'mkdir playbook_error playbook_error/node1 playbook_error/node2 playbook_error/node3 playbook_error/node4 playbook_error/node5 playbook_error/node6'
            }
        }
        stage('Create instances 1') {
            steps {
                create_n_setup_instances(5)
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
                run_playbook('./playbooks/setup/get_PeerID.yml', './test_results/setup/get_PeerID.txt', 5)
            }
        }
        stage('Update Playbooks') {
            steps {
                sh './go-dappley-multi-server-tester -function update_address'
            }
        }
        stage('Get Peer Info'){
            steps{
                run_playbook('./playbooks/getPeerInfo/check_PeerID.yml', './test_results/getPeerInfo/check_PeerID.txt', 5)
            }
        }
        stage('Send Playbooks'){
            steps{
                run_playbook('./playbooks/send/wrong_node.yml', './test_results/send/wrong_node.txt', 5)
                multi_playbooks('./playbooks/send/invalid_tip/invalid_tip_', './test_results/send/invalid_tip_', 5, 5)
                multi_playbooks('./playbooks/send/missing_flag/missing_flag_', './test_results/send/missing_flag_', 3, 5)
                multi_playbooks('./playbooks/send/invalid_file/invalid_file_', './test_results/send/invalid_file_', 2, 5)
                multi_playbooks('./playbooks/send/invalid_data/invalid_data_', './test_results/send/invalid_data_', 13, 5)
                multi_playbooks('./playbooks/send/invalid_amount/invalid_amount_', './test_results/send/invalid_amount_', 6, 5)
                multi_playbooks('./playbooks/send/invalid_address/invalid_address_', './test_results/send/invalid_address_', 6, 5)
                multi_playbooks('./playbooks/send/invalid_gas_limit/invalid_gas_limit_', './test_results/send/invalid_gas_limit_', 7, 5)
                multi_playbooks('./playbooks/send/invalid_gas_price/invalid_gas_price_', './test_results/send/invalid_gas_price_', 6, 5)
                run_playbook('./playbooks/send/single_transaction_no_tip.yml', './test_results/send/single_transaction_no_tip.txt', 5)
                run_playbook('./playbooks/send/single_transaction_with_tip.yml', './test_results/send/single_transaction_with_tip.txt', 5)
                run_playbook('./playbooks/send/multi_transaction_no_tip.yml', './test_results/send/multi_transaction_no_tip.txt', 5)
                run_playbook('./playbooks/send/multi_transaction_with_tip.yml', './test_results/send/multi_transaction_with_tip.txt', 5)
            }
        }
        stage('GetBlockchain Playbooks') {
            steps {
                run_playbook('./playbooks/getBlockchain/getBlockchainInfo.yml', './test_results/getBlockchain/getBlockchainInfo.txt', 5)
                multi_playbooks('./playbooks/getBlockchain/invalid_input/invalid_input_', './test_results/getBlockchain/invalid_input_', 12, 5)
            }
        }
        stage('EstimateGas Playbooks') {
            steps {
                run_playbook('./playbooks/estimateGas/estimate_gas.yml', './test_results/estimateGas/estimate_gas.txt', 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_tip/invalid_tip_', './test_results/estimateGas/invalid_tip/invalid_tip_', 7, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_data/invalid_data_', './test_results/estimateGas/invalid_data/invalid_data_', 13, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_amount/invalid_amount_', './test_results/estimateGas/invalid_data/invalid_data_', 7, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_address/invalid_address_', './test_results/estimateGas/invalid_data/invalid_data_', 3, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_gasPrice/invalid_gasPrice_', './test_results/estimateGas/invalid_data/invalid_data_', 6, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_gasLimit/invalid_gasLimit_', './test_results/estimateGas/invalid_data/invalid_data_', 8, 5)
                // multi_playbooks('./playbooks/estimateGas/invalid_smartContract/invalid_smartContract_', './test_results/estimateGas/invalid_data/invalid_data_', 3, 5)
            }
        }
        stage('GetBalance Playbooks') {
            steps {
                run_playbook('./playbooks/getBalance/invalid_address.yml', './test_results/getBalance/invalid_address.txt', 5)
                run_playbook('./playbooks/getBalance/missing_argument.yml', './test_results/getBalance/missing_argument.txt', 5)
            }
        }
        stage('CreateAccount Playbooks') {
            steps {
                run_playbook('./playbooks/createAccount/empty_password.yml', './test_results/createAccount/empty_password.txt', 5)
                run_playbook('./playbooks/createAccount/invalid_password.yml', './test_results/createAccount/invalid_password.txt', 5)
            }
        }
        stage('ListAddresses Playbooks') {
            steps {
                run_playbook('./playbooks/listAddresses/invalid_password.yml', './test_results/listAddresses/invalid_password.txt', 5)
            }
        }
        stage('SendFromMiner Playbooks') {
            steps {
                run_playbook('./playbooks/sendFromMiner/missing_flag.yml', './test_results/sendFromMiner/missing_flag.txt', 5)
                multi_playbooks('./playbooks/sendFromMiner/invalid_amount/invalid_amount_', './test_results/sendFromMiner/invalid_amount_', 4, 5)
                multi_playbooks('./playbooks/sendFromMiner/invalid_address/invalid_address_', './test_results/sendFromMiner/invalid_address_', 3, 5)
                run_playbook('./playbooks/sendFromMiner/single_transaction_from_miner.yml', './test_results/sendFromMiner/single_transaction_from_miner.txt', 5)
            }
        }
        stage('Smart Contract Playbooks') {
            steps {
                run_playbook('./playbooks/smartContract/smart_contract_gas_1.yml', './test_results/smartContract/smart_contract_gas_1.txt', 5)
                run_playbook('./playbooks/smartContract/smart_contract_gas_2.yml', './test_results/smartContract/smart_contract_gas_2.txt', 5)
            }
        }
        stage('Add Producer Playbooks') {
            steps {
                run_playbook('./playbooks/addProducer/add_producer.yml', './test_results/addProducer/add_producer.txt', 5)
                run_playbook('./playbooks/addProducer/add_when_max.yml', './test_results/addProducer/add_when_max.txt', 5)
                run_playbook('./playbooks/addProducer/invalid_address.yml', './test_results/addProducer/invalid_address.txt', 5)
            }
        }
        stage('Terminate Host Nodes 1') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
                remove_files()
            }
        }
        stage('Create instances 2') {
            steps {
                create_n_setup_instances(6)
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
        stage('6 Producer Test - addProducer') {
            steps {
                run_playbook('./playbooks/addProducer/6_producer_test_1.yml', './test_results/addProducer/6_producer_test_1.txt', 6)
                sh './go-dappley-multi-server-tester -function terminate -number 2'
                run_playbook('./playbooks/addProducer/6_producer_test_2.yml', './test_results/addProducer/6_producer_test_2.txt', 6)
            }
        }
        stage('Terminate Host Nodes 2') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
                remove_files()
            }
        }
        stage('Create instances 3') {
            steps {
                create_n_setup_instances(6)
            }
        }
        stage('Setup Host Nodes 3-1') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/setup.yml > ./test_results/setup/setup_4.txt'
            }
        }
        stage('Setup Host Nodes 3-2') {
            steps {
                sh 'ansible-playbook ./playbooks/removeProducer/setup.yml > ./test_results/setup/setup_5.txt'
            }
        }
        stage('6 Producer Test - removeProducer') {
            steps {
                run_playbook('./playbooks/removeProducer/6_producer_test_1.yml', './test_results/removeProducer/6_producer_test_1.txt', 6)
                sh './go-dappley-multi-server-tester -function terminate -number 2'
                run_playbook('./playbooks/removeProducer/6_producer_test_2.yml', './test_results/removeProducer/6_producer_test_2.txt', 6)
            }
        }
        stage('Terminate Host Nodes 3') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
                remove_files()
            }
        }
        stage('Create instances 4') {
            steps {
                create_n_setup_instances(6)
            }
        }
        stage('Setup Host Nodes 4-1') {
            steps {
                sh 'ansible-playbook ./playbooks/setup/setup.yml > ./test_results/setup/setup_4.txt'
            }
        }
        stage('Setup Host Nodes 4-2') {
            steps {
                sh 'ansible-playbook ./playbooks/changeProducer/setup.yml > ./test_results/setup/setup_5.txt'
            }
        }
        stage('6 Producer Test - changeProducer') {
            steps {
                run_playbook('./playbooks/changeProducer/6_producer_test_1.yml', './test_results/removeProducer/6_producer_test_1.txt', 6)
                run_playbook('./playbooks/changeProducer/6_producer_test_2.yml', './test_results/removeProducer/6_producer_test_2.txt', 6)
                run_playbook('./playbooks/changeProducer/6_producer_test_3.yml', './test_results/removeProducer/6_producer_test_3.txt', 6)
            }
        }
        stage('Send Report') {
            steps {
                sh 'zip -r test_results.zip test_results'
                sh './go-dappley-multi-server-tester -function send_result -senderEmail <Sender Email> -senderPasswd <Sender Password>'
            }
        }
        stage('Terminate Host Nodes 4') {
            steps {
                sh './go-dappley-multi-server-tester -function terminate'
                sh 'rm -r *'
            }
        }
    }
}

def create_n_setup_instances(int nodes) {
    for (int i = 1; i <= nodes; i++) {
        sh "aws ec2 run-instances --image-id ami-09e67e426f25ce0d7 --instance-type m5.large --count 1 --key-name jenkins --security-group-ids sg-0805d37807366482d --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=JenkinsTestNode${i}}]' > node${i}.txt"
    }
    sh "./go-dappley-multi-server-tester -function update -number ${nodes}"
    sh "./go-dappley-multi-server-tester -function initialize -number ${nodes}"
    sh 'sleep 20'
}

def check_initialization_status(int nodes) {
    def is_error = false
    for (int i = 1; i <= nodes; i++) {
        def error = readFile "./playbook_error/node${i}/error_${i}.txt"
        is_error = is_error || error.trim().equals("true")
    }
    return is_error
}

def run_playbook(String playbook_directory, String test_result_directory, int nodes) {
    sh "ansible-playbook ${playbook_directory} > ${test_result_directory}"

    def is_error = check_initialization_status(nodes)
    while (is_error == true) {
        echo "Blockchain initialization failed!"
        sh 'ansible-playbook ./playbooks/clear.yml'
        sh 'sleep 30'
        sh "ansible-playbook ${playbook_directory} > ${test_result_directory}"
        is_error = check_initialization_status(nodes)
    }
    echo "Blockchain initialized successfully."
}

def multi_playbooks(String playbook_directory, String test_result_directory, int number, int nodes) {
    for (int i = number; i > 0; i--) {
        run_playbook("${playbook_directory}${i}.yml", "${test_result_directory}${i}.txt", nodes)
    }
}

def remove_files() {
    sh 'rm -r node1.txt'
    sh 'rm -r node2.txt'
    sh 'rm -r node3.txt'
    sh 'rm -r node4.txt'
    sh 'rm -r node5.txt'
    sh 'rm -r instance_ids'
    sh 'rm -r hosts'
}
```