pipeline {
  agent any

  options {
    timestamps()
    ansiColor('xterm')
  }

  

  environment {
    TAG                 = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    LATEST_TAG          = 'latest'
    DOCKER_BUILDKIT     = '1'
    DISCORD_WEBHOOK_URL = credentials("DISCORD_WEBHOOK_URL")
  }

  stages {
    stage('Start Pipeline') {
      steps {
        discordSend(
          webhookURL: env.DISCORD_WEBHOOK_URL,
          description: """
          **Job:** ${env.JOB_NAME}
          **Build:** #${env.BUILD_NUMBER}
          **Branch:** ${env.BRANCH_NAME}
          **Commit:** `${env.GIT_COMMIT_HASH}`
          **Message:** ${env.GIT_COMMIT_MESSAGE}
          [View Build](${env.BUILD_URL})
          """,
          title: "🚀 Jenkins Pipeline Started!",
          footer: "Jenkins CI/CD | Started 🚀"
        )
      }
    }

    stage('Checkout') {
      steps {
        checkout scm
        sh 'ls -lah'
        sh 'git status'
        sh 'pwd'
        sh 'echo "WORKSPACE is: $WORKSPACE"'
        script {
          env.GIT_COMMIT_HASH    = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
          env.GIT_AUTHOR         = sh(script: "git log -1 --pretty=format:%an", returnStdout: true).trim()
          env.GIT_COMMIT_MESSAGE = sh(script: "git log -1 --pretty=%s", returnStdout: true).trim()
        }
      }
    }

    stage('Docker Build & Push Server') {
      steps {
        configFileProvider([configFile(fileId: 'deploy-convert-file-env', targetLocation: 'deploy.env')]) {
          sh '''#!/usr/bin/env bash
          set -Eeuo pipefail

          echo "==> Load env..."
          set -a; . deploy.env; set +a

          echo "==> Build image..."
          docker build -f Dockerfile.server \
            --build-arg APP_NAME="$APP_NAME" \
            -t "$IMAGE_NAME_SERVER:$TAG" .

          echo "==> Tag & Push..."
          docker tag "$IMAGE_NAME_SERVER:$TAG" "$REGISTRY_HOST/$IMAGE_NAME_SERVER:$TAG"
          docker tag "$IMAGE_NAME_SERVER:$TAG" "$REGISTRY_HOST/$IMAGE_NAME_SERVER:$LATEST_TAG"

          docker push "$REGISTRY_HOST/$IMAGE_NAME_SERVER:$TAG"
          docker push "$REGISTRY_HOST/$IMAGE_NAME_SERVER:$LATEST_TAG"
          '''
        }
      }
    }

    stage('Build & Push Worker'){
      steps{
          configFileProvider([configFile(fileId: 'deploy-convert-file-env', targetLocation: 'deploy.env')]) {
          sh '''#!/usr/bin/env bash
          set -Eeuo pipefail
          set -a; . deploy.env; set +a

          echo "==> Build worker image..."
          docker build -f Dockerfile.worker \
            --build-arg APP_NAME="$APP_NAME" \
            -t "$IMAGE_NAME_WORKER:$TAG" .

          echo "==> Tag & Push worker..."
          docker tag "$IMAGE_NAME_WORKER:$TAG" "$REGISTRY_HOST/$IMAGE_NAME_WORKER:$TAG"
          docker tag "$IMAGE_NAME_WORKER:$TAG" "$REGISTRY_HOST/$IMAGE_NAME_WORKER:$LATEST_TAG"

          docker push "$REGISTRY_HOST/$IMAGE_NAME_WORKER:$TAG"
          docker push "$REGISTRY_HOST/$IMAGE_NAME_WORKER:$LATEST_TAG"
          '''
        }

      }
    }


    stage('Deploy (SSH to remote)') {
      steps {
        sshagent(credentials: ['ssh-remote-dev']) {
         configFileProvider([configFile(fileId: 'deploy-convert-file-env', targetLocation: 'deploy.env')]) {
              sh '''#!/usr/bin/env bash
              set -a; . deploy.env; set +a
              echo "DEBUG Jenkins: TAG=$TAG SERVER=$IMAGE_NAME_SERVER WORKER=$IMAGE_NAME_WORKER"
              scp -o StrictHostKeyChecking=no deploy.sh ubuntu@192.168.1.100:/tmp/deploy.sh
              ssh -o StrictHostKeyChecking=no ubuntu@192.168.1.100 "TAG='${TAG}' IMAGE_NAME_SERVER='${IMAGE_NAME_SERVER}' IMAGE_NAME_WORKER='${IMAGE_NAME_WORKER}' bash /tmp/deploy.sh"
              '''
        }
      }
    }

  } 

  post {
    success {
      script {
        def executionTime = String.format("%.2f", currentBuild.duration / 1000.0)
        def timestamp = new Date().format("yyyy-MM-dd HH:mm:ss", TimeZone.getTimeZone("UTC"))
        echo "✅ Build Successful!"
        discordSend(
          webhookURL: env.DISCORD_WEBHOOK_URL,
          description: "**Job:** ${env.JOB_NAME}\n**Build:** #${env.BUILD_NUMBER}\n**Branch:** ${env.BRANCH_NAME}\n**Commit:** `${env.GIT_COMMIT_HASH}`\n**Author:** ${env.GIT_AUTHOR}\n**Message:** ${env.GIT_COMMIT_MESSAGE}\n**Execution Time:** ${executionTime} sec\n**Timestamp:** ${timestamp}\n[View Build](${env.BUILD_URL})",
          title: "✅ Build Successful!",
          footer: "Jenkins CI/CD | Success ✅"
        )
      }
    }
    failure {
      echo "❌ Build Failed! Details: ${currentBuild.result}"
      discordSend(
        webhookURL: env.DISCORD_WEBHOOK_URL,
        description: "**Job:** ${env.JOB_NAME}\n**Build:** #${env.BUILD_NUMBER}\n**Branch:** ${env.BRANCH_NAME}\n**Commit:** `${env.GIT_COMMIT_HASH}`\n**Author:** ${env.GIT_AUTHOR}\n**Message:** ${env.GIT_COMMIT_MESSAGE}\n[View Build](${env.BUILD_URL})",
        title: "❌ Build Failed!",
        footer: "Jenkins CI/CD | Failed ❌"
      )
    }
  }
}
