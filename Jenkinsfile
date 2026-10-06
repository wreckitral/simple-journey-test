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
        }
}
