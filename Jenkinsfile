pipeline {
  agent any

  options {
    timestamps()
    ansiColor('xterm') // terminal color
  }

  environment {
    TAG             = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    LATEST_TAG      = 'latest'
    DOCKER_BUILDKIT = '1'
    DISCORD_WEBHOOK_URL = credentials("DISCORD_WEBHOOK_URL")
  }

  stages {

    stage('Checkout') {
      steps {
        checkout scm
        script {
          // Lấy commit info từ git
          env.GIT_COMMIT_HASH    = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
          env.GIT_AUTHOR         = sh(script: "git log -1 --pretty=format:'%an'", returnStdout: true).trim()
          env.GIT_COMMIT_MESSAGE = sh(script: "git log -1 --pretty=%s", returnStdout: true).trim()
        }
      }
    }

    stage('Docker Build (with BuildKit secret)') {
      steps {
        withCredentials([
          file(credentialsId: 'deploy-env',         variable: 'DEPLOY_ENV'),
        ]) {
          sh '''#!/usr/bin/env bash
            set -Eeuo pipefail
            # Load deploy envs for this shell only
            set -a
            . "$DEPLOY_ENV"      # defines: REGISTRY_HOST, IMAGE_NAME, APP_NAME, ...
            set +a
            '''
        }
      }
    }

    stage('Tag & Push') {
      steps {
        withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
          sh '''#!/usr/bin/env bash
            set -Eeuo pipefail
            set -a
            . "$DEPLOY_ENV"
            set +a

            docker tag "$IMAGE_NAME:$TAG" "$REGISTRY_HOST/$IMAGE_NAME:$TAG"

            docker push "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
            '''
        }
      }
    }

 stage('Deploy (SSH to remote)') {
  steps {
    sshagent(credentials: ['ssh-remote-dev']) {
      withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
        sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . "$DEPLOY_ENV"; set +a

echo "[INFO] Jenkins node: $(hostname) / user: $(whoami)"
echo "[INFO] SSH client: $(ssh -V 2>&1 || true)"
echo "[INFO] Target: $REMOTE_USER@$REMOTE_HOST"
echo "[INFO] Image:  $REGISTRY_HOST/$IMAGE_NAME:$TAG"

# Đổ script qua stdin cho ssh, truyền tham số qua argv
cat <<'REMOTE' | ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" bash -s -- \
  "$REGISTRY_HOST" "$IMAGE_NAME" "$TAG" "$APP_NAME" "$HOST_PORT" "$APP_PORT"
set -Eeuo pipefail
REGISTRY_HOST="$1"; IMAGE_NAME="$2"; TAG="$3"
APP_NAME="$4"; HOST_PORT="$5"; APP_PORT="$6"

echo "[REMOTE] Docker: $(docker --version || true)"
docker pull "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
docker rm -f "$APP_NAME" || true
docker run -d --name "$APP_NAME" --restart=always \
  -p "$HOST_PORT:$APP_PORT" \
  "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
sleep 3
docker ps --filter name="$APP_NAME" --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
REMOTE
'''
      }
    }
  }
}



  }

  post {
    always { cleanWs() }
    success {
        script {
            def executionTime = String.format("%.2f", currentBuild.duration / 1000.0) // Time in seconds
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
        echo "❌ Build Failed!"
        discordSend(
            webhookURL: env.DISCORD_WEBHOOK_URL,
            description: "**Job:** ${env.JOB_NAME}\n**Build:** #${env.BUILD_NUMBER}\n**Branch:** ${env.BRANCH_NAME}\n**Commit:** `${env.GIT_COMMIT_HASH}`\n**Author:** ${env.GIT_AUTHOR}\n**Message:** ${env.GIT_COMMIT_MESSAGE}\n[View Build](${env.BUILD_URL})",
            title: "❌ Build Failed!",
            footer: "Jenkins CI/CD | Failed ❌"
        )
    }
  }
}