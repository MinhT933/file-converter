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
        // Define your environment variables here

        TAG = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
        LABEL = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
        VERSION = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
        BUILD_ID = "${env.BRANCH_NAME ?: 'main'}-${env.BUILD_NUMBER}"
        LATEST = "latest"
        // Remove DISCORD_WEBHOOK_URL from environment to use withCredentials
    }

    stages {

        stage('checkout'){
            steps {
                checkout scm
                script {
                    // lấy commit infor từ git
                    // %an : tác giả commit
                    // %s : thông điệp commit
                    // %ad : ngày commit
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
                    file(credentialsId: 'deploy-env',         variable: 'DEPLOY_ENV'),
                    file(credentialsId: 'envfile-portfolio',  variable: 'ENVFILE')
                    ]) {
                    sh '''#!/usr/bin/env bash
                        set -Eeuo pipefail
                        # Load deploy envs for this shell only
                        set -a
                        . "$DEPLOY_ENV"      # defines: REGISTRY_HOST, IMAGE_NAME, APP_NAME, ...
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
    always { 
        script {
            // Cleanup in script block instead of cleanWs()
            sh 'docker system prune -f || true'
            sh 'rm -rf ./* || true'
        }

    }
    success {
        script {
            try {
                def executionTime = String.format("%.2f", currentBuild.duration / 1000.0)
                def timestamp = new Date().format("yyyy-MM-dd HH:mm:ss", TimeZone.getTimeZone("UTC"))
                def branchName = env.BRANCH_NAME ?: 'unknown'
                def commitHash = env.GIT_COMMIT_HASH ?: 'unknown'
                def author = env.GIT_AUTHOR ?: 'unknown'
                def message = env.GIT_COMMIT_MESSAGE ?: 'unknown'

                echo "✅ Build Successful!"
                
                // Sử dụng withCredentials để secure Discord webhook
                withCredentials([string(credentialsId: 'DISCORD_WEBHOOK_URL', variable: 'WEBHOOK_URL')]) {
                    sh '''
                        curl -H "Content-Type: application/json" \
                        -X POST \
                        -d "{
                            \\"username\\": \\"Jenkins CI/CD\\",
                            \\"embeds\\": [{
                                \\"title\\": \\"✅ Build Successful!\\",
                                \\"description\\": \\"**Job:** ''' + env.JOB_NAME + '''\\\\n**Build:** #''' + env.BUILD_NUMBER + '''\\\\n**Branch:** ''' + branchName + '''\\\\n**Commit:** \`''' + commitHash + '''\`\\\\n**Author:** ''' + author + '''\\\\n**Message:** ''' + message + '''\\\\n**Execution Time:** ''' + executionTime + ''' sec\\\\n**Timestamp:** ''' + timestamp + '''\\\\n[View Build](''' + env.BUILD_URL + ''')\\",
                                \\"color\\": 65280,
                                \\"footer\\": {
                                    \\"text\\": \\"Jenkins CI/CD | Success ✅\\"
                                }
                            }]
                        }" \
                        "$WEBHOOK_URL"
                    '''
                }
            } catch (Exception e) {
                echo "Failed to send Discord notification: ${e.getMessage()}"
            }

            
        sh 'docker system prune -f'

        sh 'journalctl --vacuum-size=100M || true'
        }

    }
  failure {
            script {
                try {
                    def branchName = env.BRANCH_NAME ?: 'unknown'
                    def commitHash = env.GIT_COMMIT_HASH ?: 'unknown'
                    def author = env.GIT_AUTHOR ?: 'unknown'
                    def message = env.GIT_COMMIT_MESSAGE ?: 'unknown'
                    
                    echo "❌ Build Failed!"
                    
                    // Sử dụng withCredentials để secure Discord webhook
                    withCredentials([string(credentialsId: 'DISCORD_WEBHOOK_URL', variable: 'WEBHOOK_URL')]) {
                        sh '''
                            curl -H "Content-Type: application/json" \
                            -X POST \
                            -d "{
                                \\"username\\": \\"Jenkins CI/CD\\",
                                \\"embeds\\": [{
                                    \\"title\\": \\"❌ Build Failed!\\",
                                    \\"description\\": \\"**Job:** ''' + env.JOB_NAME + '''\\\\n**Build:** #''' + env.BUILD_NUMBER + '''\\\\n**Branch:** ''' + branchName + '''\\\\n**Commit:** \`''' + commitHash + '''\`\\\\n**Author:** ''' + author + '''\\\\n**Message:** ''' + message + '''\\\\n[View Build](''' + env.BUILD_URL + ''')\\",
                                    \\"color\\": 16711680,
                                    \\"footer\\": {
                                        \\"text\\": \\"Jenkins CI/CD | Failed ❌\\"
                                    }
                                }]
                            }" \
                            "$WEBHOOK_URL"
                        '''
                    }
                } catch (Exception e) {
                    echo "Failed to send Discord notification: ${e.getMessage()}"
                }
            }
        }
  }
}