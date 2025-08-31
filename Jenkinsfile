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

# --- NEW: fallback nếu chưa có APP_PORT/HEALTH_PATH trong deploy.env ---
APP_PORT="${APP_PORT:-${PORT_HTTP:-8080}}"
HEALTH_PATH="${HEALTH_PATH:-/healthz}"
export APP_PORT HEALTH_PATH

# Chạy script từ xa, TRUYỀN THÊM tham số HEALTH_PATH (đối số thứ 7)
cat <<'REMOTE' | ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" bash -s -- \
  "$REGISTRY_HOST" "$IMAGE_NAME" "$TAG" "$APP_NAME" "$HOST_PORT" "$APP_PORT" "$HEALTH_PATH"
set -Eeuo pipefail

REGISTRY_HOST="$1"; IMAGE_NAME="$2"; TAG="$3"
APP_NAME="$4"; HOST_PORT="$5"; APP_PORT="$6"; HEALTH_PATH="$7"

REMOTE_ENV_FILE="/opt/apps/${APP_NAME}/.env"
REMOTE_CREDS_FILE="/opt/secrets/firebase-creds.json"

echo "[REMOTE] Docker: $(docker --version || true)"
echo "[REMOTE] Using ENV:   $REMOTE_ENV_FILE"
echo "[REMOTE] Using CREDS: $REMOTE_CREDS_FILE"
echo "[REMOTE] Ports: host:$HOST_PORT -> app:$APP_PORT ; health: $HEALTH_PATH"

[ -f "$REMOTE_ENV_FILE" ]   || { echo "❌ Missing $REMOTE_ENV_FILE"; exit 1; }
[ -f "$REMOTE_CREDS_FILE" ] || { echo "❌ Missing $REMOTE_CREDS_FILE"; exit 1; }

docker pull "$REGISTRY_HOST/$IMAGE_NAME:$TAG"

docker rm -f "$APP_NAME" >/dev/null 2>&1 || true

set -x
docker run -d --name "$APP_NAME" --restart=always \
  --add-host=host.docker.internal:host-gateway \
  -p "$HOST_PORT:$APP_PORT" \
  --env-file "$REMOTE_ENV_FILE" \
  --mount type=bind,source="$REMOTE_CREDS_FILE",target="/app/firebase-creds.json",readonly \
  "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
set +x

sleep 2
docker ps --filter "name=$APP_NAME" --format 'table {{.Names}}\t{{.Ports}}\t{{.Status}}'

# --- NEW: Health-check dùng HEALTH_PATH + retry ---
ok=""
for i in $(seq 1 30); do
  code=$(curl -s -o /dev/null -w "%{http_code}" "http://127.0.0.1:${HOST_PORT}${HEALTH_PATH}" || true)
  if [ "$code" -ge 200 ] && [ "$code" -lt 400 ]; then
    echo "✅ Health OK ($code)"
    ok="yes"
    break
  fi
  echo "…waiting app ready ($i/30), last_code=$code"
  sleep 1
done

if [ -z "$ok" ]; then
  echo "❌ Health-check failed. Printing logs..."
  docker logs --tail 200 "$APP_NAME" || true
  exit 1
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
