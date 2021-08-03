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

### Test Cases
Send Test Cases

    - Invalid Address
        1. From an invalid address
        2. No from address
        3. To an invalid address
        4. No to address
        5. From and to an invalid address
        6. No from and to address

    - Invalid Amount
        1. Send more than the account1 balance
        2. Send negative value
        3. Send zero amount
        4. No amount
        5. Amount is string
        6. Amount is special characters

    - Invalid Tip
        1. Negative tip amount
        2. Over balance tip
        3. Equal to balance
        4. Zero tip amount
        5. No tip

    - Invalid Data
        1. No address in data [ ,"500"]
        2. No address in data ["500"]
        3. No amount in data ["account_2_address", ]
        4. No amount in data ["account_2_address"]
        5. Negative amount in data
        6. More than miner balance
        7. Amount in data is string
        8. Amount in data is special characters
        9. Invalid arguments
        10. Invalid address in data
        11. Empty address in data
        12. Empty amount in data
        13. No data

    - Invalid File
        1. Invalid file
        2. No argument

    - Invalid Gas Limit
        1. Over balance gas limit
        2. Zero gas limit
        3. Insufficient gas limit
        4. Negative gas limit
        5. Gas limit is string
        6. Gas limit is special chracter
        7. No gas limit

    - Invalid Gas Price
        1. Over balance gas price
        2. Zero gas price
        3. Negative gas price
        4. String gas price
        5. Special character gas price
        6. No gas price

    - Missing Flag
        1. No flag
        2. -from and -gasPrice only
        3. -from, -gasPrice and -amount

    - Send from wrong node
        1. Deploy transaction from the node without the account.dat file


SendFromMiner Test Cases

    - Invalid amount
        1. Send more than miner's balance
        2. Send negative amount
        3. Send zero amount
        4. Send string 
        5. Send special characters
        6. Send no amount

    - Invalid argument
        1. Send to an invalid address
        2. Send to an address without accounts.dat
        3. Send to no address

    - Missing flags
        1. No address or amount arguments
        2. No flags


CreateAccount Test Cases

    - Empty password
        1. First input is empty
        2. Second input is empty
        3. Both inputs are empty

    - Invalid password
        1. Unmatching passwords
        2. Integer password
        3. Special characters in password
        4. Empty space password


GetBalance Test Cases

    - Invalid address
        1. Get balance from an invalid address
        2. Get balance from a valid address without accounts.dat

    - Missing arguments
        1. No arguments
        2. No flag


ListAddress Test Cases

    - Invalid password
        1. ListAddress when there is no accounts.dat
        2. ListAddress with -privateKey when there is no accounts.dat
        3. ListAddress with incorrect password


Smart Contract Test Cases

    - Smart contract gas 1
        1. Smart contract deployment with gas price of 1

    - Smart contract gas 2
        1. Smart contract deployment with gas price of 2


Get Blocks & BlockchainInfo Test Cases

    - Compare blockchainInfo & blocks returning value

    - Invalid Input
        1. maxCount is 0
        2. maxCount is negative
        4. maxCount is string
        5. maxCount is special character
        6. maxCount is empty
        7. startBlockHashes is invalid
        8. startBlockHashes is valid but doesn't exist
        9. startBlockHashes is 0
        10. startBlockHashes is negative
        11. startBlockHashes is random string
        12. startBlockHashes is special character
        13. startBlockHashes is emtpy

EstimateGas Test Cases

    - Estimate gas price and deploy smart contract witht the estimated gasLimit

    - Invalid address
        1. from address is invalid
        2. from address is valid but no accounts.dat
        3. from address is empty

    - Invalid smartContract
        1. smart contract address is invalid
        2. smart contract address is valid but does not exist
        3. smart contract address is empty

    - Invalid amount
        1. amount is 0
        2. amount is negative
        3. amount more than the balance 
        4. amount is string
        5. amount is special character
        6. amount is empty
        7. amount is equal to balance

    - Invalid tip
        1. tip is 0
        2. tip is negative
        3. tip more than the balance 
        4. tip is string
        5. tip is special character
        6. tip is empty
        7. tip is equal to balance

    - Invalid gasLimit
        1. gasLimit is 0
        2. gasLimit is negative
        3. gasLimit is more than the balance
        4. gasLimit is insufficient
        5. gasLimit is string
        6. gasLimit is special character
        7. gasLimit is equal to balance
        8. gasLimit is empty

    - Invalid gasPrice
        1. gasPrice is 0
        2. gasPrice is negative
        3. gasPrice is more than the balance
        4. gasPrice is string
        5. gasPrice is special character
        6. gasPrice is empty

    - Invalid data
        1. data amount is 0
        2. data amount is negative
        3. data amount is over balance
        4. data amount is string
        5. data amount is special character
        6. data amount is empty 1
        7. data amount is empty 2
        8. data address is invalid
        9. data address is valid but no accounts.dat
        10. data address is empty 1
        11. data address is empty 2
        12. data arguments are invalid
        13. data is empty


Get Peer Info Test Case

    - Check if the peer ID exists in the log file on each server.


Add Producer Test Cases

    - Invalid address
        1. address is invalid
        2. address is integer
        3. address is special characters
        4. address is empty

    - Add producer when max producer number is already reached

    - Initialize a 6 node blockchain with 4 producers,
      then add 2 more producers to the blockchain. Shut down 2
      servers and check if the blockchain if still mining or not.
      If the blockchain height continues to grow, then there is an error.


Remove Producer Test Cases

    - Initialize a 6 node blockchain with 6 produers,
      then delete 1 producer and check if the blockchain continues.
      Delete another producer and shutdown 2 of the blockchain nodes.
      If the blockchain stops mining, then there is an error.

    - Invalid input
        1. height is less than the current blockahin height
        2. height is string
        3. height is speical characters
        4. height is decimal
        5. height is empty

    - Run deleteProducer method from the node that isn't running any producer in the blockchain OR 
      Run deleteProducer twice from the node that is running a proudcer in the blockchain


Change Producer Test Cases

    - Change one of the produer and check if the balance of the miner is increasing after change.

    - Change one of the producer's address to an address that is already inside the blockchain.
      If the producers' addresses remain the same then pass, otherwise fail.

    - Change 3 producers' address to a valid, but non-existing (non of the nodes in the blockchain contains the address)
      producer addresses. If the blockchain continues to mine, then fail.

    - Invalid address
        1. address is invalid
        2. address is non-sense string
        3. address is integer
        4. address is special characters
        5. address is empty

    - Invalid height
        1. height is less than the current blockchain height
        2. height is string
        3. height is special character
        4. height is decimal
        5. height is empty

    - Run changeProducer command twice
        1. Run change producer twice with different addresses but the same height. The blockchain is expected to discard
           the input of the first changeProdcer command and change the producer address to the input of the second changeProducer
           command.

        2. Run change producer twice with different addresses but the same height. The address of the second changeProducer
           command is invalid. The blockchain is expected to discard the input of the second changeProdcer command and change
           the producer address to the input of the first changeProducer command.