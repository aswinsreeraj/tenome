pipeline {
    agent none

    environment {
        GOCACHE = "${WORKSPACE}/.cache/go-build"
        GOPATH  = "${WORKSPACE}/.go"
    }

    stages {
        // stage('Format') {
        //     agent {
        //         docker {
        //             image 'golang:1.26'
        //         }
        //     }

        //     steps {
        //         sh 'test -z "$(gofmt -l .)"'
        //     }
        // }

        // stage('Vet') {
        //     agent {
        //         docker {
        //             image 'golang:1.26'
        //         }
        //     }
        //     steps {
        //         sh 'go vet ./...'
        //     }
        // }

        stage('Test') {
            agent {
                docker {
                    image 'golang:1.26'
                }
            }
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build') {
            agent {
                docker {
                    image 'golang:1.26'
                }
            }
            steps {
                sh 'go build -o tenome ./cmd/server/main.go'
            }
        }

        stage('Docker Build') {
            agent any
            steps {
                sh '''
                docker build \
                    -t tenome:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Deploy') {
            agent any
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