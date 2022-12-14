def NETWORK = UUID.randomUUID().toString()
def CHROME_CONTAINER = "chrome-${NETWORK}"
def APP_CONTAINER = "app-${NETWORK}"

pipeline {
  agent none

  environment {
    WEBDRIVER_HOST = "${CHROME_CONTAINER}"
    APP_URL = "http://${APP_CONTAINER}:3000"
  }

  stages {
    stage('Bootstrap services') {
      agent any
      steps {
        script {
          sh "docker network create ${NETWORK} || true"
          CHROME = docker.image('seleniarm/standalone-chromium:latest').run("--name ${CHROME_CONTAINER} --network ${NETWORK} --shm-size=2g")
        }
      }
    }

    stage('Build and Test') {
      agent {
        docker {
            image 'golang:1.19'
            args "-u root --network ${NETWORK} --name ${APP_CONTAINER}"
        }
      }

      steps {
        dir('webapp') {
            sh 'go test'
        }
      }
    }

    stage('Linters') {
      agent any
      steps {
        script {
          // sonarqube depends on the previous stage workspace, so these two steps cannot be run in parallel
          def scanner = tool 'sonarqube'
          withSonarQubeEnv('sonarqube') {
            sh "${scanner}/bin/sonar-scanner -Dsonar.projectKey=practical -Dsonar.source=webapp"
          }
        }
      }
    }

  }
  post {
    always {
      node(null) {
        script {
          if (CHROME) {
            CHROME.stop()
          }

          sh "docker network rm ${NETWORK} || true"
        }
      }
    }
  }
}
