pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    environment {
        TAG                = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
        LATEST_TAG         = 'latest'
        DOCKER_BUILDKIT    = '1'
        DISCORD_WEBHOOK_URL = credentials("DISCORD_WEBHOOK_URL")
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'ls -lah'
                sh 'git status'
                sh 'pwd'
                sh 'echo "WORKSPACE is: $WORKSPACE"'

                script {
                    env.GIT_COMMIT_HASH    = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                    env.GIT_AUTHOR         = sh(script: "git log -1 --pretty=format:'%an'", returnStdout: true).trim()
                    env.GIT_COMMIT_MESSAGE = sh(script: "git log -1 --pretty=%s", returnStdout: true).trim()
                }
            }
        }

        stage('Docker Build (with BuildKit secret)') {
            steps {
                configFileProvider([configFile(fileId: 'deploy-convert-file-env', variable: 'DEPLOY_ENV_FILE')]) {
                    sh '''#!/usr/bin/env bash
                    set -Eeuo pipefail
                    set -a
                    . "$DEPLOY_ENV_FILE"
                    set +a

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
                    docker push "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
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

# Chạy script từ xa
cat <<'REMOTE' | ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" bash -s -- \
  "$REGISTRY_HOST" "$IMAGE_NAME" "$TAG" "$APP_NAME" "$HOST_PORT" "$APP_PORT"
set -Eeuo pipefail

REGISTRY_HOST="$1"; IMAGE_NAME="$2"; TAG="$3"
APP_NAME="$4"; HOST_PORT="$5"; APP_PORT="$6"

REMOTE_ENV_FILE="/opt/apps/${APP_NAME}/.env"
REMOTE_CREDS_FILE="/opt/secrets/firebase-creds.json"

echo "[REMOTE] Docker: $(docker --version || true)"
echo "[REMOTE] Using ENV:   $REMOTE_ENV_FILE"
echo "[REMOTE] Using CREDS: $REMOTE_CREDS_FILE"

# Kiểm tra đủ file trước khi chạy
[ -f "$REMOTE_ENV_FILE" ]   || { echo "❌ Missing $REMOTE_ENV_FILE"; exit 1; }
[ -f "$REMOTE_CREDS_FILE" ] || { echo "❌ Missing $REMOTE_CREDS_FILE"; exit 1; }

# Kéo image mới
docker pull "$REGISTRY_HOST/$IMAGE_NAME:$TAG"

# Dự phòng rollback: lưu image/container cũ (nếu có)
OLD_IMG="$(docker inspect -f '{{.Image}}' "$APP_NAME" 2>/dev/null || true)"
docker rm -f "$APP_NAME" >/dev/null 2>&1 || true

# Run container mới
set -x
docker run -d --name "$APP_NAME" --restart=always \
  --add-host=host.docker.internal:host-gateway \
  -p "$HOST_PORT:$APP_PORT" \
  --env-file "$REMOTE_ENV_FILE" \
  --mount type=bind,source="$REMOTE_CREDS_FILE",target="/app/firebase-creds.json",readonly \
  "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
set +x

sleep 3
docker ps --filter "name=$APP_NAME" --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'

# (Tuỳ chọn) Health-check đơn giản, /healthz hoặc /ping tuỳ app của bạn
if command -v curl >/dev/null 2>&1; then
  curl -fsS "http://127.0.0.1:${HOST_PORT}/healthz" >/dev/null || {
    echo "❌ Health-check failed. Printing logs..."
    docker logs --tail 200 "$APP_NAME" || true
    exit 1
  }
fi

echo "✅ Deploy OK"
REMOTE
'''
      }
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
