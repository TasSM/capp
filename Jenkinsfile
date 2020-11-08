pipeline {
    agent any
    parameters {
        string(name: 'DOCKER_REGISTRY', defaultValue: 'docker-registry.labnet', description: 'Location of your docker repository')
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git credentialsId: 'github-ssh', url: 'git@github.com:TasSM/capp.git'
            }
        }
        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${params.DOCKER_REGISTRY}/capp:latest ."
            }
        }
        stage('Push Docker Image to Registry') {
            steps {
                sh "docker push ${params.DOCKER_REGISTRY}/capp:latest"
            }
        }
        stage('Cleanup') {
            steps {
                sh "docker image rm ${params.DOCKER_REGISTRY}/capp:latest"
            }
        }
    }
}