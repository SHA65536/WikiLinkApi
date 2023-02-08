pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                // Build the application
                sh "go build ./cmd/linkapi"
            }
        }
        stage('Test') {
            steps {
                // Test the application
                sh "go test ./..."
            }
        }
        stage('Push') {
            steps {
                // Push to S3
                sh "aws s3 cp ./linkapi s3://cloudschoolproject-buildartifacts/linkapi"
            }
        }
    }
}