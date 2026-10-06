pipeline {
    agent any
        stages {
            stage('Checkout') {
                steps {
                    checkout scm
                        sh 'ls -l'
                }
            }

            stage('Test') {
                steps {
                    sh '''
                        docker run --rm \
                        -v jenkins_home:/var/jenkins_home \
                        -w "$WORKSPACE" \
                        golang:1.27 go test ./...

                        '''
                }
            }

            stage('Build Image') {
                steps {
                    script {
                        env.GIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    }

                    sh 'docker build --build-arg VERSION=${GIT_SHORT} -t devops-test:${GIT_SHORT} .'
                }
            }

            stage('Push') {
                steps {
                    withCredentials([usernamePassword(
                                credentialsId: 'ghcr-token',
                                usernameVariable: 'GHCR_USER',
                                passwordVariable: 'GHCR_TOKEN'
                                )]) {
                        sh '''
                            echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
                            docker tag devops-test:${GIT_SHORT} ghcr.io/wreckitral/devops-test:${GIT_SHORT}
                            docker push ghcr.io/wreckitral/devops-test:${GIT_SHORT}
                            docker logout ghcr.io

                            '''
                    }
                }

            }
        }
}
