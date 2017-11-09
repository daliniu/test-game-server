#!/usr/bin/env python
# -*- coding:utf-8 -*-
# __author__:jk
# __version__: 0.0.1

import time
import os
import argparse
import sys
import shutil
import stat
import json
import re
import tarfile
from ftplib import FTP, FTP_TLS

GOROOT = '/usr/local/Cellar/go/1.8/libexec'  # go安装路径
PROJECT = '/Users/lcg/code/wolf-backend'  # 代码路径
TARGET = '/Users/lcg/data'
# GOROOT = '/data/services/go'  # go安装路径
# PROJECT = '/data/deploy/wolf-backend'  # 代码路径
# TARGET = '/data/goapp/wolf'
GOPATH = '%s:%s/vendor' % (PROJECT, PROJECT)
GO = '%s/bin/go' % GOROOT

# FTP信息
FTP_IP = 'ftp.cloudstone.qq.com'
FTP_PORT = 9054
FTP_USER = 'fworks_1180'
FTP_PASS = 'fanWW@2017'

aliasMapping = {
     'admin-socket-srv':'adminsocket',
     'admin-statistics-srv':'statistics',
     'aggregation-srv':'aggregation',
     'agora-dynamickey-srv':'agora',
     'api-admin-srv':'apiadmin',
     'api-login-srv':'apilogin',
     'api-srv':'api',
     'api-third-srv':'apithird',
     'bag-srv':'bag',
     'cms-srv':'cms',
     'connector-srv':'connector',
     'credit-srv':'credit',
     'game-room-srv':'gameroom',
     'game-rule-srv':'gamerule',
     'game-stats-srv':'gamestats',
     'group-srv':'group',
     'idgo-srv':'idgo',
     'idip-srv':'idip',
     'mall-srv':'mall',
     'match-srv':'match',
     'mission-srv':'mission',
     'push-srv':'push',
     'qcloud-im-srv':'qcloudim',
     'relation-srv':'relation',
     'room-manager-srv':'roommanager',
     'security-srv':'security',
     'tss-srv':'tss',
     'user-report-srv':'report',
     'user-srv':'user',
     'user-status-srv':'userstatus'
}

def pack(branch, services, noupload):
    currentDir = os.path.dirname(os.path.realpath(__file__)) # 当前目录的路径
    now = time.strftime('%Y%m%d%H%M%S')  # 当前的时间
    gitID = updateProgram(branch)

    targetDir = '%s/bin_%s_%s_%s' % (TARGET, branch, now, gitID)  # go install路径
    print '目标:', targetDir


    if services[0] == 'all':
        services = aliasMapping.keys()
    compileServices(targetDir, services)

    compileConfigs(targetDir)
    prepareOtherFiles(targetDir)

    appfile = generateTarFile(targetDir, branch, now, gitID)
    if noupload == False:
        upload(appfile)

def updateProgram(branch):
    print '更新代码'  # 更新代码

    os.chdir(PROJECT)
    releaseBranch = 'release-%s' % branch
    os.system('git pull')
    os.system('git checkout %s' % releaseBranch)
    os.system('git pull origin %s' % releaseBranch)
    return os.popen('git rev-parse --short HEAD').read().rstrip("\n")

def compileServices(targetDir, services):
    print '编译所有服务'

    for service in services:
        if aliasMapping.has_key(service) == False:
            print '别名不存在:', service
            continue
        alias = aliasMapping[service]
        print '更新 %s -> %s' % (service, alias)

        GOBIN = os.path.join(targetDir, 'ww', alias, 'bin')
        os.system('GOROOT=%s GOBIN=%s GOPATH=%s %s install %s' % (GOROOT, GOBIN, GOPATH, GO, service))
        shutil.move(os.path.join(GOBIN, service), os.path.join(GOBIN, alias))

def prepareOtherFiles(targetDir):
    print '复制 C++动态库lib'
    shutil.copytree(os.path.join(PROJECT, 'lib'), os.path.join(targetDir, 'ww', 'lib'))

    print '处理tss配置文件'
    tssTargetDir = os.path.join(targetDir, 'ww', 'tss', 'conf')
    if not os.path.exists(tssTargetDir):
        os.makedirs(tssTargetDir)
    tssSdkConfPathFile = os.path.join(PROJECT, 'src/share/tencent/tsssdk/tss_sdk_conf_path.xml.deploy')
    tssSdkConfFile = os.path.join(PROJECT, 'src/share/tencent/tsssdk/tss_sdk_conf.xml.deploy')
    shutil.copy(tssSdkConfPathFile, os.path.join(tssTargetDir, 'tss_sdk_conf_path.xml'))
    shutil.copy(tssSdkConfFile, os.path.join(tssTargetDir, 'tss_sdk_conf.xml'))

    print '处理app.json'
    shutil.copy(os.path.join(PROJECT, 'app.json'), os.path.join(targetDir, 'ww'))

def compileConfigs(targetDir):
    print '生成配置文件'

    configDir = os.path.join(PROJECT, 'config')
    targetConfigDir = os.path.join(targetDir, 'ww', 'config')
    if not os.path.exists(targetConfigDir):
        os.makedirs(targetConfigDir)

    for file in os.listdir(configDir):
        if file[-4:-1] + file[-1] != "json":
            continue
        if file != 'env.json':
            print file
            syncConfigFile(os.path.join(configDir, file), os.path.join(targetConfigDir, file))

def syncConfigFile(filePath, targetFilePath):
    defaultEnv = 'test'
    conf = json.load(open(filePath))
    fileName = os.path.basename(filePath)[0:-5]

    confStrs = {'tpre': '', 'audit': '', 'pro': ''}
    for env in ['test', 'tpre', 'audit', 'pro']:
        if conf.has_key(env) == False:
            conf[env] = {
                'wx': conf[defaultEnv]['wx'],
                'qq': conf[defaultEnv]['wx'],
                'guest': conf[defaultEnv]['wx']
            }

        for zone in conf[env]:
            if (env == defaultEnv and zone == 'wx') == False:
                conf[env][zone] = mergeConfig(env, conf[env][zone], conf[defaultEnv]['wx'])

        if env != 'test':
            confStrs[env] = replaceResourcePath(env, json.dumps(conf[env], ensure_ascii=True))

    f = open(targetFilePath, 'w')
    f.write("{\"tpre\":%s,\"audit\":%s,\"pro\":%s}" % (confStrs['tpre'], confStrs['audit'], confStrs['pro']))
    f.close()

def mergeConfig(env, c1, c2):
    for key in c2.keys():
        if c1.has_key(key) == False:
            c1[key] = c2[key]
            continue

        if isinstance(c2[key], dict):
            c1[key] = mergeConfig(env, c1[key], c2[key]);
        elif isinstance(c2[key], list):
            if c1.has_key(key) == False or len(c1[key]) == 0:
                c1[key] = c2[key]
    return c1

def replaceResourcePath(env, str):
    if env == 'pro':
        str, num = re.subn(r'web\.test\.ifanju\.com', 'image.ww.qq.com/pro/web', str)
        str, num = re.subn(r'ifanju-wolf\.cdn\.ifanju\.com', 'image.ww.qq.com/pro', str)
    elif env == 'tpre' or env == 'audit':
        str, num = re.subn(r'web\.test\.ifanju\.com', 'image.ww.qq.com/test/web', str)
        str, num = re.subn(r'ifanju-wolf\.cdn\.ifanju\.com', 'image.ww.qq.com/test', str)
    return str

def generateTarFile(targetDir, branch, now, gitID):
    appfile = 'ww_%s_%s_%s.tar.gz' % (branch, now, gitID)
    print '打包%s' % (appfile)

    os.chdir(targetDir)
    tar = tarfile.open(appfile, 'w:gz')
    for files in os.listdir(targetDir):
        tar.add(files)
    tar.close()

    print '计算md5值'
    os.system('md5sum %s' % (appfile))
    return appfile

def upload(appfile):
    print '上传FTP'
    # os.system('lftp -c 'open ftp.cloudstone.qq.com -p 9054 -u fworks_1180,fanWW@2017;cd server;put %s'' % (appfile))
    os.system('lftp -c "open %s -p %d -u %s,%s;cd server;put %s"' % (FTP_IP, FTP_PORT, FTP_USER, FTP_PASS, appfile))

if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog='部署打包', description='使用该脚本来打包编译服务端代码并上传到腾讯的FTP服务器')
    parser.add_argument('--noupload', default=False, type=bool, help='无需上传')
    parser.add_argument('branch', help='服务端程序版本分支')
    parser.add_argument('services', nargs='+', help='要打包的服务名')
    args = parser.parse_args()

    pack(args.branch, args.services, args.noupload)