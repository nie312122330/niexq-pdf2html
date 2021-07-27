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
