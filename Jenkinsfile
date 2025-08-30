pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm') // terminal color
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
      withCredentials([file(credentialsId: 'deploy-env', variable: 'DEPLOY_ENV')]) {
        sh '''#!/usr/bin/env bash
set -Eeuo pipefail
set -a; . "$DEPLOY_ENV"; set +a

# SSH vào remote server và xóa container cũ
ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" <<EOF
  echo "[INFO] Stopping and removing container if it exists..."
  docker stop be-server-convert-file || true
  docker rm -f be-server-convert-file || true
  echo "[INFO] Container removed"
EOF

# Tiến hành tạo container mới
docker pull "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
docker run -d --name "be-server-convert-file" --restart=always \
  -p "$HOST_PORT:$APP_PORT" \
  "$REGISTRY_HOST/$IMAGE_NAME:$TAG"
sleep 3
docker ps --filter name="be-server-convert-file" --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
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
