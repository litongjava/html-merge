# html-merge
将 body,css,javascript合并成一个文件
## 构建
```shell
go build
```
## 使用
需要配合服务端:ping-base64-webapi一起使用
```
html-merge -b index.html -c index.css -j index.js -o output.html 
```
```shell
-b body部分代码文件,不包含body标签
-c css部分代码文件,
-m JavaScript部分代码文件
-o 输出的文件
```

