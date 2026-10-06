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

            stage('Deploy') {
                steps {
                    sh '''
                        set -e
                        docker rm -f extract stage 2>/dev/null || true

                        # take the new binary out of the image we just built
                        docker create --name extract devops-test:${GIT_SHORT}
                        docker cp extract:/app/app ./app.new
                        docker rm extract

                        # put it into the volume
                        docker create --name stage -v app-bin:/app alpine
                        docker cp ./app.new stage:/app/app.new
                        docker rm stage

                        # keep the old binary, swap, restart
                        docker run --rm -v app-bin:/app alpine sh -c "cp /app/app /app/app.prev && mv /app/app.new /app/app"
                        docker restart simple-journey-test

                        sleep 3
                        if docker run --rm --network container:simple-journey-test curlimages/curl -sf localhost:8080 | grep "version=${GIT_SHORT}"; then
                            echo "deploy ok"
                        else
                            echo "health check failed, rolling back"
                            docker run --rm -v app-bin:/app alpine sh -c "cp /app/app.prev /app/app"
                            docker restart simple-journey-test
                            exit 1
                        fi
                        '''
                }
            }}
}
