#!groovy
pipeline {
// agent needs to build go, so start with golang:1.10-alpine
    agent { docker { image 'golang:1.10-alpine' }  }
    environment {
// could be useful for now
        NG_BRANCH = env.BRANCH_NAME.toLowerCase()
    }
    options { disableConcurrentBuilds() }
    stages {

        stage ('BuildImage') {
            steps {
                echo "Building ${NG_BRANCH} Number ${env.BUILD_NUMBER} - home: ${env.HOME}"
                whoami
                pwd
                go version
                echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >>  /etc/apk/repositories
                apk update
                apk add docker
                echo "tagging images with:registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}"
                docker build -t registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}/mra-content-service:${BUILD_NUMBER} -f docker/Dockerfile .
                docker push registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}/mra-content-service:${BUILD_NUMBER}
                docker rmi $(docker images -f "dangling=true" -q) || true
        }

   }
    post {
        always {
            print "Post operations"
        }
        success {
            print "SUCCESSFUL build ${env.JOB_NAME} [${env.BUILD_NUMBER}] (${env.BUILD_URL})"
          script {
            if (currentBuild.getPreviousBuild() &&
              currentBuild.getPreviousBuild().getResult().toString() != "SUCCESS") {
                emailext body: "Content Service Recovery: ${env.BUILD_URL}", recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'RequesterRecipientProvider']], subject: "BUILD ERROR:${env.BRANCH_NAME}"
            }
          }
        }
        failure {
            emailext body: "Content Service Error on: ${env.BUILD_URL}", recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'RequesterRecipientProvider']], subject: "BUILD ERROR:${env.BRANCH_NAME}"
        }
        unstable {
            print "UNSTABLE JOB build ${env.JOB_NAME} [${env.BUILD_NUMBER}] (${env.BUILD_URL})"
        }
    }
}