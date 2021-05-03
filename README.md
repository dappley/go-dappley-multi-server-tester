# ansible-dappley
ansible playbook

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