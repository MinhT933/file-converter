// === Helper: gửi thông báo webhook (Discord / proxy nội bộ) ===
// - Sinh JSON bằng Groovy (JsonOutput) để tránh lỗi escape
// - Thêm field `sender` theo schema API của bạn
def notifyDiscord(String title, String description, int color) {
  // Tạo payload JSON an toàn
  def payload = groovy.json.JsonOutput.toJson([
    sender  : 'jenkins',
    username: 'Jenkins CI/CD',
    embeds  : [[
      title      : title,
      description: description,
      color      : color
    ]]
  ])

  // Escape single-quote để nhét JSON vào chuỗi shell
  def escaped = payload.replace("'", "'\"'\"'")

  withCredentials([string(credentialsId: 'DISCORD_WEBHOOK_URL', variable: 'WEBHOOK_URL')]) {
    sh """#!/usr/bin/env bash
set -Eeuo pipefail
curl -sS --fail-with-body \\
  -H 'Content-Type: application/json' \\
  -X POST \\
  -d '${escaped}' "\$WEBHOOK_URL" || true
"""
  }
}

pipeline {
  agent any

  triggers {
    githubPush()
  }

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

    stage('Docker Build (with BuildKit secret)') {
      steps {
        withCredentials([
          file(credentialsId: 'deploy-env',        variable: 'DEPLOY_ENV'),
          file(credentialsId: 'envfile-portfolio', variable: 'ENVFILE')
        ]) {
          sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a
. "$DEPLOY_ENV"   # defines: REGISTRY_HOST, IMAGE_NAME, APP_NAME, REMOTE_HOST, ...
set +a

DOCKER_BUILDKIT=1 docker build \
  --secret id=dotenv,src="$ENVFILE" \
  -t "$IMAGE_NAME:$TAG" .
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

# Đổ script qua stdin cho remote host
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
    always {
      script {
        echo "JF MARKER: v2025-08-29-1" // đánh dấu để chắc chắn Jenkins dùng file mới
        // Cleanup
        sh 'docker system prune -f || true'
        sh 'rm -rf ./* || true'
      }
    }

    success {
      script {
        def desc = "Build #${env.BUILD_NUMBER} completed successfully for job ${env.JOB_NAME}"
        notifyDiscord("✅ Build Successful!", desc, 65280)
        sh 'journalctl --vacuum-size=100M || true'
      }
    }

    failure {
      script {
        def desc = "Build #${env.BUILD_NUMBER} failed for job ${env.JOB_NAME}"
        notifyDiscord("❌ Build Failed!", desc, 16711680)
      }
    }
  }
}
