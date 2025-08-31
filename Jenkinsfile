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

    stage('Docker Build (with BuildKit secret)') {
      steps {
        configFileProvider([configFile(fileId: 'deploy-convert-file-env', variable: 'DEPLOY_ENV_FILE')]) {
          sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . "$DEPLOY_ENV_FILE"; set +a
docker build -f Dockerfile.server --build-arg APP_NAME="$APP_NAME" -t "$IMAGE_NAME:$TAG" .
'''
        }
      }
    }

    stage('Tag & Push') {
      steps {
        configFileProvider([configFile(fileId: 'deploy-convert-file-env', targetLocation: 'deploy.env')]) {
          sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . deploy.env; set +a
docker tag "$IMAGE_NAME:$TAG" "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
docker tag "$IMAGE_NAME:$TAG" "$REGISTRY_HOST/$IMAGE_NAME:$LATEST_TAG"
docker push "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
docker push "$REGISTRY_HOST/$IMAGE_NAME:$LATEST_TAG"
'''
        }
      }
    }

    stage('Deploy (SSH to remote)') {
      steps {
        sshagent(credentials: ['ssh-remote-dev']) {
          configFileProvider([configFile(fileId: 'deploy-convert-file-env', targetLocation: 'deploy.env')]) {
            sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . deploy.env; set +a

: "${STACK_DIR:?missing STACK_DIR}"
SERVICE="${SERVICE:-app}"
HOST_PORT="${HOST_PORT:-8081}"
HEALTH_PATH="${HEALTH_PATH:-/healthz}"
COMPOSE_FILE_REMOTE="${COMPOSE_FILE_REMOTE:-docker-compose.prod.yml}"

ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" bash -s -- \
  "$STACK_DIR" "$SERVICE" "$HOST_PORT" "$HEALTH_PATH" "$COMPOSE_FILE_REMOTE" <<'REMOTE'
set -Eeuo pipefail
STACK_DIR="$1"; SERVICE="$2"; HOST_PORT="$3"; HEALTH_PATH="$4"; COMPOSE_FILE="$5"

cd "$STACK_DIR" || { echo "❌ No such dir: $STACK_DIR"; exit 1; }
[ -f "$COMPOSE_FILE" ] || { echo "❌ Missing $COMPOSE_FILE"; ls -l; exit 1; }
export COMPOSE_FILE="$COMPOSE_FILE"

docker compose pull "$SERVICE"
docker compose up -d --no-deps "$SERVICE"
docker compose ps

ok=""; for i in $(seq 1 30); do
  code=$(curl -s -o /dev/null -w "%{http_code}" "http://127.0.0.1:${HOST_PORT}${HEALTH_PATH}" || true)
  [ "$code" -ge 200 ] && [ "$code" -lt 400 ] && { echo "✅ Health OK ($code)"; ok=1; break; }
  echo "…waiting ($i/30), code=$code"; sleep 1
done
[ -n "$ok" ] || { echo "❌ Health failed"; docker compose logs --no-color --tail 200 "$SERVICE" || true; exit 1; }
REMOTE
'''
          }
        }
      }
    }
  }  // end stages

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
