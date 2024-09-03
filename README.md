# 通用 HTML 转 PDF

## 1.原理

直接利用 chrome headless 后台进行转换

## 2.前提条件

运行环境中必须安装 chrome 或 chrome 的可执行程序

## 3.基础实践

```bash
chrome --headless --disable-gpu --print-to-pdf=/root/aa.pdf  https://www.baidu.com
```

存在以下问题

> 1. 默认存在页面页脚
> 2. 无法修改目标尺寸

## 4.高级实践

1. 利用 chromedp 开源库，调用 chrome 生成 pdf，该库主要用于爬虫及自动化测试，支持自定义页眉页脚，目标尺寸，是否打印背景色，边距调整

> chromedp 开源库：<https://github.com/chromedp/chromedp>

2. 利用 chromedp 构建一个 HttpServer，做一个网络版服务，【 go 版本为:1.16.4】

> 源码地址：<https://github.com/nie312122330/niexq-pdf2html>
>
> 二进制下载：<https://github.com/nie312122330/niexq-pdf2html/releases/tag/1.0.0>

3. 转换效果

   ![效果图](./imgs/baidu.png)

4. 包含两个 http 接口，都为 POST 请求

5. URL 直接转换为 PDF 文件 <http://127.0.0.1:19444/niexq-html2pdf/pub/url2Pdf.do>

> 请求参数为:

```json
{
  "httpUrl": "https://www.baidu.com",
  "landscape": false,
  "displayHeaderFooter": false,
  "printBackground": true,
  "paperWidth": 8.5,
  "paperHeight": 11,
  "marginTop": 0,
  "marginBottom": 0,
  "marginLeft": 0.4,
  "marginRight": 0.4,
  "pageRanges": ""
}
```

> 响应，响应的 data 字段就是生成的 PDF 文件，使用: <http://127.0.0.1:19444/{data}>就可直接访问

```json
{
  "code": 0,
  "msg": "",
  "count": 0,
  "pageCount": 0,
  "warn": "",
  "serverTime": "2021-07-27 07:48:54",
  "data": "pdfdir/2021/07/27/EE072E19879F4F7486A90318D4B47D8F.pdf",
  "extData": {}
}
```

2. html 文本转换为 PDF 文件
   <http://127.0.0.1:19444/niexq-html2pdf/pub/html2Pdf.do>

> 请求参数为:

```json
{
  "htmlTxt": "<!DOCTYPE html><html><body><h1> 聂xxxxx</h1><img src=\"https://sanzi-oss.widthsoft.com/fixdir/fix_icon/c_system/40.png\" /></body></html>",
  "landscape": false,
  "displayHeaderFooter": false,
  "printBackground": true,
  "paperWidth": 8.5,
  "paperHeight": 11,
  "marginTop": 0,
  "marginBottom": 0,
  "marginLeft": 0.4,
  "marginRight": 0.4,
  "pageRanges": ""
}
```

> 响应，响应的 data 字段就是生成的 PDF 文件，使用:<http://127.0.0.1:19444/{data}>就可直接访问

```json
{
  "code": 0,
  "msg": "",
  "count": 0,
  "pageCount": 0,
  "warn": "",
  "serverTime": "2021-07-27 07:48:54",
  "data": "pdfdir/2021/07/27/EE072E19879F4F7486A90318D4B47D8F.pdf",
  "extData": {}
}
```

## 5.容器化实践（CentOS）

1. 增加 google-chrome-repo

```bash
tee /root/buildimages/google-chrome.repo <<-'EOF'
[google-chrome]
name=google-chrome
baseurl=http://dl.google.com/linux/chrome/rpm/stable/$basearch
enabled=1
gpgcheck=1
gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub
EOF
```

2. 拷贝 windows 字体

> 1. 在/root/buildimages 中建立文件夹 /root/buildimages/winfonts
> 2. 在 windows 上 c:\windows\fonts 中的所有文件到 /root/buildimages/winfonts

3. 拷贝引用程序

> 1. 拷贝 pdf2html 到 /root/buildimages/pdf2html
> 2. 修改 app_conf.yaml 文件中的 chromeConf.execPath 为: /usr/bin/google-chrome
> 3. 拷贝 app_conf.yaml 到 /root/buildimages/app_conf.yaml

4. 编写 DockerFile

```bash
#编写DockerFile
tee /root/build-chromdp/Dockerfile <<-'EOF'
FROM centos:7.9.2009

RUN curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
RUN yum clean all;yum makecache
RUN yum install -y vim wget unzip openssh openssh-clients net-tools kde-l10n-Chinese glibc-common fontconfig psmisc bind-utils curl
RUN echo Asia/Shanghai > /etc/timezone && localedef -c -f UTF-8 -i zh_CN zh_CN.utf8
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime

COPY google-chrome.repo  /etc/yum.repos.d/google-chrome.repo
WORKDIR /root
RUN yum install -y google-chrome-stable --nogpgcheck
#字体相关
RUN mkdir -p /usr/share/fonts/winfonts
COPY winfonts  /usr/share/fonts/winfonts
RUN chmod -R 755 /usr/share/fonts/winfonts
#安装 ttmkfdir 来搜索目录中的字体信息
RUN yum -y install ttmkfdir
RUN ttmkfdir -e /usr/share/X11/fonts/encodings/encodings.dir
RUN sed -i 's#<dir>/usr/share/fonts</dir>#<dir>/usr/share/fonts</dir>\n\t<dir>/usr/share/fonts/winfonts</dir>#g' /etc/fonts/fonts.conf
RUN fc-cache

#僵尸进程问题处理  【https://github.com/krallin/tini】
RUN curl -o /bin/tini https://sanzi-oss.widthsoft.com/fixdir/build-docker-img/tini-v0.19.0
RUN chmod +x /bin/tini
ENTRYPOINT ["/bin/tini","--"]

#应用程序
COPY html2pdf  /root/
COPY app_conf.yaml  /root/app_conf.yaml
RUN chmod -R 777 /root/html2pdf

WORKDIR /root
CMD ["/root/html2pdf"]

EOF


docker build -t registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v1.0.0 .
docker push  registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v1.0.0 


#运行
docker pull registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v1.0.0
docker stop html2pdf;docker rm html2pdf
docker run -d --restart=always --name html2pdf  -u root --privileged --cgroupns host \
 -p 19444:19444 -m 1G  \
 -v /root/html2pdf/:/root/pdfdir \
registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v1.0.0


#统计总进程数
pstree -p |wc -l

docker exec html2pdf pstree -p |wc -l

#杀进程
if (($(docker exec html2pdf pstree -p |wc -l)>4000));then docker restart html2pdf; fi

```

5. 构建|测试|运行

```bash
#构建
cd /root/buildimages
docker build -t local/html2pdf:v1 .

#测试
docker run -it --rm --name test  local/html2pdf:v1 /bin/bash


#正式运行
docker run -d --restart=always --name html2pdf -p 19444:19444 local/html2pdf:v1


#调用
curl -X POST  -H "Content-Type: application/json;charset=\'UTF-8\'" -d '{"httpUrl":"https://www.baidu.com","landscape":false,"displayHeaderFooter":false,"printBackground":true,"paperWidth":8.5,"paperHeight":11,"marginTop":0,"marginBottom":0,"marginLeft":0.4,"marginRight":0.4,"pageRanges":""}' http://192.168.0.251:19444/niexq-html2pdf/pub/url2Pdf.do

#循环调用
for i in {1..5000}
 do
 curl -X POST  -H "Content-Type: application/json;charset=\'UTF-8\'" -d '{"httpUrl":"https://www.baidu.com","landscape":false,"displayHeaderFooter":false,"printBackground":true,"paperWidth":8.5,"paperHeight":11,"marginTop":0,"marginBottom":0,"marginLeft":0.4,"marginRight":0.4,"pageRanges":""}' http://192.168.0.253:19444/niexq-html2pdf/pub/url2Pdf.do
 done






```

6.非容器安装（CentOS）

````bash
#增加 google-chrome-repo
tee /etc/yum.repos.d/google-chrome.repo <<-'EOF'
[google-chrome]
name=google-chrome
baseurl=http://dl.google.com/linux/chrome/rpm/stable/$basearch
enabled=1
gpgcheck=1
gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub
EOF

#安装google-chrome
yum install -y google-chrome-stable --nogpgcheck 

#制作service






````

7、使用openanolis作为基础镜像

```bash


tee /root/build-chromdp/Dockerfile <<-'EOF'
FROM registry.cn-chengdu.aliyuncs.com/width-public/openanolis:23-20230704

RUN yum clean all;yum makecache
RUN yum install -y vim wget unzip openssh openssh-clients net-tools glibc-common fontconfig psmisc bind-utils curl
RUN echo Asia/Shanghai > /etc/timezone 
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime

#安装chrome依赖
RUN mkdir -p /root/lib-fonts
WORKDIR /root/lib-fonts
# https://rpmfind.net/linux/rpm2html/search.php?query=liberation-fonts-common&submit=Search+...&system=&arch=
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-narrow-fonts-1.07.6-15.el10.noarch.rpm
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-fonts-2.1.5-10.el10.noarch.rpm
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-mono-fonts-2.1.5-10.el10.noarch.rpm
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-sans-fonts-2.1.5-10.el10.noarch.rpm
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-serif-fonts-2.1.5-10.el10.noarch.rpm
RUN wget https://rpmfind.net/linux/centos-stream/10-stream/AppStream/x86_64/os/Packages/liberation-fonts-common-2.1.5-10.el10.noarch.rpm
RUN yum localinstall -y *.rpm --disablerepo=*

#下载windows字体
RUN wget https://sanzi-oss.widthsoft.com/fixdir/anolisos/winfonts.tar
RUN tar -xvf winfonts.tar && mv winfonts /usr/share/fonts/winfonts && chmod -R 755 /usr/share/fonts/winfonts

#安装 ttmkfdir 来搜索目录中的字体信息
WORKDIR /root
RUN yum -y install ttmkfdir
RUN ttmkfdir -e /usr/share/X11/fonts/encodings/encodings.dir
RUN sed -i 's#<dir>/usr/share/fonts</dir>#<dir>/usr/share/fonts</dir>\n\t<dir>/usr/share/fonts/winfonts</dir>#g' /etc/fonts/fonts.conf
RUN fc-cache

#安装chrome
WORKDIR /root
RUN curl -o /etc/yum.repos.d/google-chrome.repo https://sanzi-oss.widthsoft.com/fixdir/build-docker-img/anolisos/google-chrome.repo
RUN yum install -y google-chrome-stable --nogpgcheck

#僵尸进程问题处理  【https://github.com/krallin/tini】
WORKDIR /root
RUN curl -o /bin/tini https://sanzi-oss.widthsoft.com/fixdir/build-docker-img/tini-v0.19.0
RUN chmod +x /bin/tini
ENTRYPOINT ["/bin/tini","--"]

WORKDIR /root
RUN rm -rf /root/*
RUN wget https://sanzi-oss.widthsoft.com/fixdir/anolisos/html2pdf/app_conf.yaml
RUN wget https://sanzi-oss.widthsoft.com/fixdir/anolisos/html2pdf/html2pdf
RUN chmod +x *
CMD ["/root/html2pdf"]

EOF
docker build -t registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v2.0.0 .
docker run -it --rm registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v2.0.0

docker stop html2pdf;docker rm html2pdf
docker run -d --restart=always --name html2pdf  -u root --privileged --cgroupns host \
 -p 19444:19444 -m 1G  \
 -v /root/html2pdf/:/root/pdfdir \
registry.cn-chengdu.aliyuncs.com/width-public/html2pdf:v2.0.0

```



