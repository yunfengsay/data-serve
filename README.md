## 接受ios传图服务

快捷指令传图, 在主要工作环境是linux的情况下，ios可以传图给自己的主机，如果是win可以稍微改下代码中的文件保存路径即可
## 使用方法
```bash
sudo npm i  @yunfengsay/data-serve -g
```
会自动创建一个开机启动任务，这个开机启动任务会监听 45531 端口，并提供一个 post /uploads api，该api可以接收文件上传任务并且保存到 ~/Pictures 文件夹下

ios下的快捷指令
https://www.icloud.com/shortcuts/84ae3d6f0e6c4fffb128619fe6914e75