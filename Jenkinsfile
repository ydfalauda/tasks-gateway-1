def app

pipeline {
    
    agent any
    /* lets create a more complex pipeline */

    stages {
        
        /* first build image */    
        stage('Build') {
            steps {
                script {
                    checkout scm
                    app = docker.build("claasv1/tasks-gateway")
                    docker.withRegistry('http://index-int.alauda.cn', 'index-int') {
                        app.push("${env.BUILD_NUMBER}")
                        app.push("latest")
                    }
                }

            }
            
        }

        stage('Create Temporary App') {
            steps {
                sh 'docker-compose -f docker-compose.1.yml -p jenkins up -d'
                echo 'will wait a little bit'
                sleep 10
            }
        }

        stage('Start testing') {
            try {
                steps {
                    sh 'docker pull index.alauda.cn/alaudaorg/tasks-integration:latest'
                    sh 'docker run -t --rm --network jenkins_default --link jenkins_gateway_1:gateway index.alauda.cn/alaudaorg/tasks-integration:latest'
                }
            } finally {
                sh 'docker-compose -f docker-compose.1.yml -p jenkins down'

            }
        }


    }
}