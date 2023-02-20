pipeline {
    agent any

    tools { go '1.20' }

    stages {
        stage('Build') {
            steps {
                // Build the application
                sh "go build ./cmd/wikilinkapi"
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
                sh "aws s3 cp ./linkapi s3://cloudschoolproject-buildartifacts/wikilinkapi_v_$BUILD_NUMBER"
            }
        }
    }
}