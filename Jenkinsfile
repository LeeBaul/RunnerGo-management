def createVersion() {
    // 定义一个版本号作为当次构建的版本，输出结果 20191210175842_69
    return new Date().format('yyyy-MM-dd HH:mm:ss')
}



pipeline {
    agent  any
    options{
        buildDiscarder(logRotator(numToKeepStr: '30'))        
        timeout(time: 10, unit: 'MINUTES')
        disableConcurrentBuilds()
    }


    environment {
        //GOARCH="amd64"
        //GOHOSTOS="linux"
        //GOOS="linux"
        GOPATH="${env.WORKSPACE}/cmd"
        GO111MODULE="on"
        GOPROXY="https://goproxy.cn"
        //GOSUMDB="off"
        GOCACHE="/root/.cache/go-build"
        def BUILDVERSION = createVersion()
        feishuWb="https://open.feishu.cn/open-apis/bot/v2/hook/e1ea72fe-d073-4304-bf9d-21fbb602986e"
        scriptPath="/data/docker/jenkins/fasongfeishu"
        git="172.17.101.191:2223/kunpeng/kp-management.git"
        appname="kp-management"
    }

   parameters {
        //选择标签选项    
        listGitBranches(
            branchFilter: 'refs/heads/(.*)',
            defaultValue: "",
            name: 'WEBHOOK_REF',
            type: 'PT_BRANCH_TAG',
            remoteURL: 'ssh://git@172.17.101.191:2223/kunpeng/kp-management.git',
            credentialsId: 'bcb57094-5e19-4f14-9c15-0046ed8a91fe',
            selectedValue: 'DEFAULT',
            sortMode: 'DESCENDING'
            )
     }

        triggers {
        // 由 Generic Webhook Trigger 提供的触发器
        GenericTrigger(
                // 参数, 在这里配置的参数会被配置为环境变量
                // 支持使用 JSONPath 和 XPath 进行提取
                // 提取来源为 webhook 发送的 request body 中的内容
                // GitLab 发送的 request body 内容格式可通过文档查看, 文档地址为你的 GitLab 访问地址 + /help/user/project/integrations/webhooks
                genericVariables: [
                        // 提取分支名称, 格式为 refs/heads/{branch}
                        [key: 'WEBHOOK_REF', value: '$.ref'],
                        // 提取用户显示名称
                        [key: 'WEBHOOK_USER_NAME', value: '$.user_name'],
                        // 提取最近提交 id
                        [key: 'WEBHOOK_RECENT_COMMIT_ID', value: '$.commits[-1].id'],
                        // 提取最近提交 message
                        [key: 'WEBHOOK_RECENT_COMMIT_MESSAGE', value: '$.commits[-1].message'],
                        // 如需更多参数可通过查看 request body 参数文档进行提取
                ],

                // 项目运行消息, 会显示在 Jenkins Blue 中的项目活动列表中
                causeString: '$WEBHOOK_USER_NAME 推送 commit 到 $WEBHOOK_REF 分支',
                //
                // token, 因为对外的 webhook url 是一致的, 当希望触发某个特定的 Job 时可以为每个 Job 配置不同的 token, 然后在 webhook 中配置该参数
                token: 'kp-management123',

                // 打印通过 genericVariables 配置的变量
                printContributedVariables: true,
                // 打印 request body 内容
                printPostContent: true,

                // 避免使用已触发工作的信息作出响应
                silentResponse: false,

                // 可选的正则表达式过滤, 比如希望仅在 master 分支上触发, 你可以进行如下配置
                regexpFilterText: '$WEBHOOK_REF',
                regexpFilterExpression: 'refs/heads/' + env.BRANCH_NAME
        )
    } 





    stages {
        stage('清理工作空间'){steps {cleanWs()}}
        stage('checkout') {
            steps {
                 wrap([$class: 'BuildUser']) {
 
                   script {
                       if ( env.BUILD_USER ) {
                             env.WEBHOOK_USER_NAME="${env.BUILD_USER}"
                       }
                    
                   }
                   print "++ $env.WEBHOOK_USER_NAME"
				}
                checkout([$class: 'GitSCM', 
                          branches: [[name: "${WEBHOOK_REF}"]], 
                          doGenerateSubmoduleConfigurations: false, 
                          extensions: [], 
                          gitTool: 'Default', 
                          submoduleCfg: [], 
                          userRemoteConfigs: [[credentialsId: 'bcb57094-5e19-4f14-9c15-0046ed8a91fe', url: "ssh://git@${git}",]]
                        ])
            }
        }
        // stage('编译golang'){
        //     when {
        //        expression {   return JOB_BASE_NAME != "fbs" }
        //     }
        //     steps {
        //          sh '''
        //          /usr/local/go/bin/go build -o ${appname} cmd/main.go
        //          '''
        //     }
        // }


        stage('编译golang'){
            steps {
                 sh '''
                 /usr/local/go18/bin/go build -o ${appname} cmd/main.go
                 '''
            }
        }


    stage('发布到开发环境') {
        when {
               expression {   return JOB_BASE_NAME == "dev" }
            }
            steps {
                 ansiblePlaybook( 
                	installation: 'ansible2.9.25', 
                	inventory: "${env.WORKSPACE}/hosts", 
                	limit: '172.17.101.189', 
                	playbook: "${env.WORKSPACE}/playbook.yaml",
                	extraVars: [
                		project_name:'Dev_Apipost-socket',
                		deploy_path:'/data/',
                		WORKSPACE:"${env.WORKSPACE}/${appname}",
                        config:"${env.WORKSPACE}/cmd/config.dev.yaml",
                        private:"${env.WORKSPACE}/cmd/private.pem",
                        app_name:"${appname}",
                        backfile: 'false',
                        dport: '28081',
                        devsever: "devApipost-socket"
                	])
            }
    }


    stage('发布到测试环境') {
        when {
              expression {   return JOB_BASE_NAME == "test" }
            }


            steps {
                 ansiblePlaybook( 
                	installation: 'ansible2.9.25', 
                	inventory: "${env.WORKSPACE}/hosts", 
                	limit: '172.17.101.188', 
                	playbook: "${env.WORKSPACE}/playbook.yaml",
                	extraVars: [
                		project_name:'kp-management',
                		deploy_path:'/data/Application/',
                        app_name:"${appname}",
                        config:"${env.WORKSPACE}/configs/test.yaml",
                		WORKSPACE:"${env.WORKSPACE}/${appname}",
                        backfile: 'false',
                        dport: '20001'
                	])
            }
    }

    stage('发布到生产环境') {
        when {
              expression {   return JOB_BASE_NAME == "master" }
            }
            steps { 
                 input message: "发布或停止"
                 ansiblePlaybook( 
                	installation: 'ansible2.9.25', 
                	inventory: "${env.WORKSPACE}/hosts", 
                	limit: '39.101.160.1471', 
                	playbook: "${env.WORKSPACE}/playbook.yaml",
                	extraVars: [
                		project_name:'Apipost-socket',
                		deploy_path:'/data/',
                        app_name:"${appname}",
                        config:"${env.WORKSPACE}/cmd/config.online.yaml",
                        private:"${env.WORKSPACE}/cmd/private.pem",
                		WORKSPACE:"${env.WORKSPACE}/${appname}",
                        backfile: 'true',
                        dport: '8088'
                	])
            }
    }

    stage('版本标签') {
        when {
                //branch 'master'
              expression {   return JOB_BASE_NAME == "master" }
            }
            steps {
            script{
                    def GIT_TAG = new Date().format("yyyy-MM-dd_HH-mm-ss");
                    //git branch: "${WEBHOOK_REF}", credentialsId: 'bcb57094-5e19-4f14-9c15-0046ed8a91fe', url: "ssh://git@${git}"
                    sh """
                        git config --system user.email "jay@apipost.cn"
                        git config --system user.name "guobaobao"
                        git checkout "${env.BRANCH_NAME}"
                        git tag -a ${GIT_TAG} -m 'jenkins'
                        git config --system push.default simple
                        git push origin --tags
                    """
            }
        }
    }

    }

    post {

        success {
            script {
             if (JOB_BASE_NAME =~ /(test|master|rc|fbs)/ )  {
            sh """
                "$scriptPath"  "【发布成功】完成发布" "green"  $WEBHOOK_USER_NAME  $JOB_NAME  $BUILD_NUMBER   "$BUILDVERSION"    $BUILD_URL   "$feishuWb"
            """
             }
            }
        }
        failure {
            script {
             if (JOB_BASE_NAME =~ /(test|master|rc|fbs)/ )  {
            sh """
                "$scriptPath"  "【发布失败】发布失败" "red"  "$WEBHOOK_USER_NAME"  "$JOB_NAME"  "$BUILD_NUMBER"   "$BUILDVERSION"    "$BUILD_URL"   "$feishuWb"
            """
                         }
            }
        }
        aborted {
             script {
             if (JOB_BASE_NAME =~ /(test|master|rc|fbs)/ )  {           
            sh """
                "$scriptPath"  "【发布失败】终止发布" "purple" "$WEBHOOK_USER_NAME"  "$JOB_NAME"  "$BUILD_NUMBER"   "$BUILDVERSION"    "$BUILD_URL"   "$feishuWb"
            """  
                         }
            }          
        }
        unstable {
            script {
             if (JOB_BASE_NAME =~ /(test|master|rc|fbs)/ )  {
            sh """
                "$scriptPath"  "【发布失败】发布不稳定" "red" "$WEBHOOK_USER_NAME"  "$JOB_NAME"  "$BUILD_NUMBER"   "$BUILDVERSION"    "$BUILD_URL"   "$feishuWb"
            """  
                         }
            }              
        }
    }

}


