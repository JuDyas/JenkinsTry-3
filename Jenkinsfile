pipeline {
    agent any

    environment {
        DOCKER_IMAGE = 'JuDyas/JenkinsTry'
        GIT_CREDENTIALS_ID = 'git-credentials'
        DOCKER_HOST = "unix:///var/run/docker.sock"
    }

    stages {

        stage('Checkout') {
            steps {
                sh 'docker ps'
                checkout scm
            }
        }

        stage('Calculate Version') {
             steps {
                 script {
                     sh "git fetch --all --prune"

                     def releaseBranch = sh(
                         script: """
                             git symbolic-ref --short HEAD 2>/dev/null || \
                             git for-each-ref --format='%(refname:short)' --points-at HEAD refs/remotes/origin/release/* || \
                             echo ''
                         """,
                         returnStdout: true
                     ).trim()

                     if (!releaseBranch.startsWith("release/v")) {
                         error "ERROR: Current branch '${releaseBranch}' doesn't match the required release/v* format!"
                     }

                     def majorVersion = releaseBranch.replaceAll("release/v", "")

                     def featureMerges = sh(script: "git log origin/${releaseBranch} --merges --grep=\"Merge branch 'feature/*'\" --oneline | wc -l", returnStdout: true).trim()
                     def bugfixMerges = sh(script: "git log origin/${releaseBranch} --merges --grep=\"Merge branch 'bugfix/*'\" --oneline | wc -l", returnStdout: true).trim()

                     env.MAJOR = majorVersion
                     env.MINOR = featureMerges
                     env.PATCH = bugfixMerges
                     env.VERSION = "${MAJOR}.${MINOR}.${PATCH}".toLowerCase()

                     echo "Calculated version: ${env.VERSION}"
                 }
             }
         }

        stage('Build Docker Image') {
            steps {
                script {
                    sh "docker build -t ${DOCKER_IMAGE}:${env.VERSION} ."
                }
            }
        }

        stage('Test in Builder') {
            steps {
                script {
                    sh "docker build -t builder-test --target builder -f Dockerfile ."
                    sh "docker run --rm builder-test go test ./..."
                }
            }
        }

        stage('Tag Release') {
            steps {
                sshagent([env.GIT_CREDENTIALS_ID]) {
                    sh """
                        git config user.name "jenkins"
                        git config user.email "jenkins@example.com"
                        git tag v${VERSION}
                        git push origin v${VERSION}
                    """
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    sh 'docker rm -f app-container || true'
                    sh "docker run -d -p 8081:8081 --name app-container ${DOCKER_IMAGE}:${env.VERSION}"
                    echo "Application deployed successfully. Running version: ${env.VERSION}"
                }
            }
        }
    }

    post {
        always {
            script {
                sh "docker image prune -f"
            }
        }
        success {
            echo 'Build completed successfully!'
        }
        failure {
            echo 'Build failed!'
        }
    }
}
