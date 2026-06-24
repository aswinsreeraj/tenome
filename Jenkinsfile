pipeline {
    agent {
        docker {
            image 'golang:1.26'
        }
    }

    environment {
        GOCACHE = "${WORKSPACE}/.cache/go-build"
        GOPATH  = "${WORKSPACE}/.go"
    }

    stages {
        stage('Format') {
            steps {
                sh 'test -z "$(gofmt -l .)"'
            }
        }

        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o tenome ./cmd/server/main.go'
            }
        }

        stage('Docker Build') {
            steps {
                sh '''
                docker build \
                    -t tenome:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    docker compose down
                    docker compose up -d
                '''
                // docker rm -f tenome || true

                // docker run -d \
                //     --name tenome \
                //     -p 8050:8050 \
                //     -e DB_PATH=/data/crawler.db \
                //     -e REDIS_ADDR=host.docker.internal:6379 \
                //     -v /opt/tenome/data:/app/data \
                //     tenome:${BUILD_NUMBER}
                // '''
            }
        }
    }

    post {
        success {
            archiveArtifacts artifacts: 'tenome'
        }
    }
}