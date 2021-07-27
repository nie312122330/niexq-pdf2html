# 通用HTML转PDF

## 1.原理

直接利用chrome headless后台进行转换

## 2.前提条件

运行环境中必须安装chrome或chrome的可执行程序

## 3.基础实践

```bash
chrome --headless --disable-gpu --print-to-pdf=/root/aa.pdf  https://www.baidu.com
```

存在以下问题

> 1. 默认存在页面页脚
> 2. 无法修改目标尺寸

## 4.高级实践

利用chromedp开源库，调用chrome生成pdf，该库主要用于爬虫及自动化测试，支持自定义页眉页脚，目标尺寸，是否打印背景色，边距调整

> chromedp开源库：<https://github.com/chromedp/chromedp>

利用chromedp构建一个HttpServer，做一个网络版服务，【 go版本为:1.16.4】

> 源码地址：<https://github.com/nie312122330/niexq-pdf2html>
>
> 二进制下载：<https://github.com/nie312122330/niexq-pdf2html/releases/tag/1.0.0>

包含两个http接口，都为POST请求

1.URL直接转换为PDF文件 <http://127.0.0.1:19444/niexq-html2pdf/pub/url2Pdf.do>

请求参数为:

```json
{
    "httpUrl":"https://www.baidu.com",
    "landscape":false,
    "displayHeaderFooter":false,
    "printBackground":true,
    "paperWidth":8.5,
    "paperHeight":11,
    "marginTop":0,
    "marginBottom":0,
    "marginLeft":0.4,
    "marginRight":0.4,
    "pageRanges":""
}
```

响应，响应的data字段就是生成的PDF文件，使用: <http://127.0.0.1:19444/{data}>就可直接访问

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

2.html文本转换为PDF文件 <http://127.0.0.1:19444/niexq-html2pdf/pub/html2Pdf.do>

请求参数为:

```json
{
    "htmlTxt":"<!DOCTYPE html><html><body><h1> 聂xxxxx</h1><img src=\"https://sanzi-oss.widthsoft.com/fixdir/fix_icon/c_system/40.png\" /></body></html>",
    "landscape":false,
    "displayHeaderFooter":false,
    "printBackground":true,
    "paperWidth":8.5,
    "paperHeight":11,
    "marginTop":0,
    "marginBottom":0,
    "marginLeft":0.4,
    "marginRight":0.4,
    "pageRanges":""
}
```

响应，响应的data字段就是生成的PDF文件，使用:<http://127.0.0.1:19444/{data}>就可直接访问

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