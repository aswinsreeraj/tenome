pipeline {
    agent none

    environment {
        GOCACHE = "/tmp/go-build-cache"
        GOMODCACHE = "/tmp/go-mod-cache"
        // GOPATH  = "${WORKSPACE}/.go"
    }

    stages {
        stage('Debug') {
            agent {
                docker {
                    image 'golang:1.26'
                }
            }

            steps {
                sh '''
                pwd
                go version
                go env
                '''
            }
        }
        stage('CI') {
            agent {
                docker {
                    image 'golang:1.26'
                }
            }

            stages {
                stage('Format') {
                    steps {
                        sh '''
                            gofmt -l . | grep -v "^.go/" | tee fmt.out
                            test ! -s fmt.out
                        '''
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
                docker rm -f tenome || true

                docker run -d \
                    --name tenome \
                    -p 8050:8050 \
                    -e DB_PATH=/data/crawler.db \
                    -e REDIS_ADDR=host.docker.internal:6379 \
                    -v /opt/tenome/data:/data \
                    tenome:${BUILD_NUMBER}
                '''
                //     docker compose down
                //     docker compose up -d
                // '''
            }
        }
    }

    // post {
    //     success {
    //         archiveArtifacts artifacts: 'tenome'
    //     }
    // }
}