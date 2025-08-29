// Jenkinsfile
import groovy.json.JsonOutput

// Marker để chắc chắn Jenkins đang dùng file mới
def JF_MARKER = "v2025-08-29-6"

// ===== Helper: gửi thông báo tới webhook =====
// - Nếu URL là Discord -> gửi payload Discord chuẩn (trả 204 khi OK)
// - Nếu không phải Discord -> gửi payload "custom" có trường sender tối thiểu
def notifyWebhook(String title, String description, int color) {
  withCredentials([string(credentialsId: 'DISCORD_WEBHOOK_URL', variable: 'WEBHOOK_URL')]) {

    def isDiscord = WEBHOOK_URL.contains('discord.com') || WEBHOOK_URL.contains('discordapp.com')
    Map payloadMap

    if (isDiscord) {
      payloadMap = [
        content : "${title} ${description}",
        username: "Jenkins CI/CD",
        embeds  : [[ title: title, description: description, color: color ]]
      ]
    } else {
      payloadMap = [
        sender  : [ id: 1, login: 'jenkins', html_url: 'https://github.com/jenkins' ],
        username: 'Jenkins CI/CD',
        embeds  : [[ title: title, description: description, color: color ]]
      ]
    }

    def payload = JsonOutput.toJson(payloadMap)
    echo "Webhook payload (type=${isDiscord ? 'discord' : 'custom'}) => ${payload}"
    def escaped = payload.replace("'", "'\"'\"'")

    sh """#!/usr/bin/env bash
set -Eeuo pipefail
tmp=\$(mktemp)
code=\$(curl -sS -o "\$tmp" -w "%{http_code}" \\
  -H 'Content-Type: application/json' -X POST \\
  -d '${escaped}' "\$WEBHOOK_URL" || true)
echo "Webhook HTTP status: \$code"
if [ "\$code" -ge 400 ]; then
  echo "Webhook response body:"; sed -n '1,200p' "\$tmp"
fi
rm -f "\$tmp"
"""
  }
}

pipeline {
  agent any

  triggers { githubPush() }

  options {
    timeout(time: 1, unit: 'HOURS')
    timestamps()
    ansiColor('xterm')
  }

  environment {
    TAG      = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    LABEL    = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    VERSION  = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    BUILD_ID = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
    LATEST   = "latest"
  }

  stages {

    stage('Banner') {
      steps { echo "JF MARKER: ${JF_MARKER}" }
    }

    stage('Check Credentials (fail fast)') {
      steps {
        script { echo "→ Checking credentials..." }
        withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
          sh 'echo "OK: deploy-env (secret file exists)"'
        }
        sshagent(credentials: ['ssh-remote-dev']) {
          sh 'echo "OK: ssh-remote-dev (ssh key visible)"'
        }
        withCredentials([string(credentialsId: 'DISCORD_WEBHOOK_URL', variable: 'WEBHOOK_URL')]) {
          sh 'echo "OK: DISCORD_WEBHOOK_URL (secret text bound)"'
        }
      }
    }

    stage('Checkout') {
      steps {
        checkout scm
        script {
          env.GIT_COMMIT_HASH    = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
          env.GIT_AUTHOR         = sh(script: "git log -1 --pretty=format:'%an'", returnStdout: true).trim()
          env.GIT_COMMIT_MESSAGE = sh(script: "git log -1 --pretty=%s", returnStdout: true).trim()
          env.GIT_COMMIT_DATE    = sh(script: "git log -1 --pretty=format:'%ad'", returnStdout: true).trim()
        }
      }
    }

    stage('Docker Build') {
      steps {
        withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
          sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . "$DEPLOY_ENV"; set +a

echo "[BUILD] IMAGE_NAME=$IMAGE_NAME  TAG=$TAG"
DOCKER_BUILDKIT=1 docker build -t "$IMAGE_NAME:$TAG" .
'''
        }
      }
    }

    stage('Tag & Push') {
      steps {
        withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
          sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . "$DEPLOY_ENV"; set +a

echo "[PUSH] -> $REGISTRY_HOST/$IMAGE_NAME:$TAG"
docker tag  "$IMAGE_NAME:$TAG" "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
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

echo "[INFO] Target: $REMOTE_USER@$REMOTE_HOST"
echo "[INFO] Image : $REGISTRY_HOST/$IMAGE_NAME:$TAG"

# Chạy script trên remote qua stdin
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
docker ps --filter name="$APP_NAME" --format 'table {{.Names}}\\t{{.Image}}\\t{{.Status}}'
REMOTE
'''
          }
        }
      }
    }
  }

  post {
    always {
      script {
        sh 'docker system prune -f || true'
        sh 'rm -rf ./* || true'
      }
    }
    success {
      script {
        def desc = "Build #${env.BUILD_NUMBER} completed successfully for job ${env.JOB_NAME}"
        notifyWebhook("✅ Build Successful!", desc, 65280)
        sh 'journalctl --vacuum-size=100M || true'
      }
    }
    failure {
      script {
        def desc = "Build #${env.BUILD_NUMBER} failed for job ${env.JOB_NAME}"
        notifyWebhook("❌ Build Failed!", desc, 16711680)
      }
    }
  }
}
