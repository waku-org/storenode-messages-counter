#!/usr/bin/env groovy
pipeline {
  agent { label 'linux' }

  parameters {
    string(
      name: 'IMAGE_TAG',
      defaultValue: params.IMAGE_TAG ?: 'latest',
      description: 'Optional Docker image tag to push.'
    )
  }

  options {
    disableConcurrentBuilds()
    /* manage how many builds we keep */
    buildDiscarder(logRotator(
      numToKeepStr: '20',
      daysToKeepStr: '30',
    ))
  }

  environment {
    DOCKER_REGISTRY = 'harbor.status.im'
    IMAGE_NAME = 'wakuorg/storenode-messages-counter'
    IMAGE_DEFAULT_TAG = "${env.GIT_COMMIT.take(8)}"
  }

  stages {
    stage('Build') {
      steps { script {
        image = docker.build(
          "${DOCKER_REGISTRY}/${IMAGE_NAME}:${IMAGE_DEFAULT_TAG}",
          "--build-arg='commit=${GIT_COMMIT}' .",
        )
      } }
    }

    stage('Deploy') {
      when { expression { params.IMAGE_TAG != '' } }
      steps { script {
        withDockerRegistry([
          credentialsId: 'harbor-telemetry-robot',
          url: 'https://${DOCKER_REGISTRY}',
        ]) {
          image.push(params.IMAGE_TAG)
        }
      } }
    }
  }

  post {
    cleanup { cleanWs() }
  }
}
