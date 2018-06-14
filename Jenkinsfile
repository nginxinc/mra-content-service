#!groovy
pipeline {
// agent needs to build go, so start with golang:1.10-alpine
    agent none
    environment {
// could be useful for now
        NG_BRANCH = env.BRANCH_NAME.toLowerCase()
    }
    options { disableConcurrentBuilds() }
    stages {
        stage ('BuildImage') {
          agent { docker { image 'golang:1.10-alpine' }  }
          steps {
          // this is the list of commands that will be run in the agent
            sh '''
              echo "Building ${NG_BRANCH} Number ${env.BUILD_NUMBER} - home: ${env.HOME}"
              go version
              echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >>  /etc/apk/repositories
              apk update
              apk add docker
              echo "tagging images with:registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}"
              docker build -t registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}/mra-content-service:${env.BUILD_NUMBER} .
              docker push registry.ctrl.nginx.com/mra/ngrefarch/mra-content-service/${NG_BRANCH}/mra-content-service:${env.BUILD_NUMBER}
              docker rmi $(docker images -f "dangling=true" -q) || true
            '''
          }
      }
      stage ('DeployContainerToK8s') {
        agent { docker { image 'lachlanevenson:k8s-kubectl:v1.9.8' } }
        steps {
          sh '''
            echo `kubectl version`
            echo done
          '''
        }
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