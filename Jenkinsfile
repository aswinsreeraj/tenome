pipeline {
    agent {
        docker {
            image 'golang:1.26'
        }
    }

    stages {
        stage('Hello') {
            steps {
                sh 'echo hello'
            }
        }
    }
}