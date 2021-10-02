#!/usr/bin/env node

const fs = require('fs')
const http = require('http')
const covert = require('heic-convert')
const {
  copy
} = require("copy-paste");

const {
  pipeline,
  Readable,
} = require('stream')
const log = (...args) => console.log.bind('[log] ')(...args)
const PORT = 45531;

// 生成12位的随机字符串
const generateRandomString = () => {
  return Math.random().toString(36).substr(2, 20) + new Date().getTime()
}

const streamToJson = async (req) => {
  return new Promise((resolve, reject) => {
    let body = '';
    req.on('data', buff => body += buff);
    req.on('end', () => {
      const result = JSON.parse(body);
      resolve(result)
    });
    req.on('error', reject);
  })
}

const homedir = require('os').homedir();

// 默认都是debain系统, 放在home的pictures下
const PicDir = `${homedir}/Pictures`;
const server = http.createServer((req, res) => {
  const {
    method,
    url,
    headers
  } = req;
  const contentType = headers['content-type'];
  const ctx = {
    req,
    res,
    method,
    url,
    headers,
    contentType
  }
  if (url.startsWith('/uploads')) {
    uploadServe(ctx)
  }
  if (url.startsWith('/shareclipboard')) {
    shareClipBoard(ctx)
  }
})

const shareClipBoard = async (ctx) => {
  const {
    req,
    res
  } = ctx;

  const {
    data
  } = await streamToJson(req);
  copy(data)
  res.end('ok')
}

const uploadServe = (ctx) => {
  const {
    req,
    res,
    contentType
  } = ctx;
  const fileType = contentType.split('/')[1];
  const fileName = `${PicDir}/${generateRandomString()}`;
  // 普通格式图片
  if (['png', 'jpeg', 'jpg', 'gif'].includes(fileType)) {
    const fileStream = fs.createWriteStream(`${fileName}.${fileType}`);
    pipeline(req, fileStream, (err) => {
      if (err) {
        log(err)
      }
      res.end('ok')
    })
    return
  }
  // 苹果拍摄的格式保存
  if (['heic'].includes(fileType)) {
    const buffer = [];
    req.on('data', chunk => {
      buffer.push(chunk)
    })
    req.on('end', async () => {
      const data = await covert({
        buffer: Buffer.concat(buffer),
        format: 'JPEG',
        quality: 1,
      })
      const fileStream = fs.createWriteStream(`${fileName}.png`);
      Readable.from(data).pipe(fileStream).on('finish', () => {
        res.end('ok')
      })
    })
    return
  }
  const fileStream = fs.createWriteStream(`${fileName}.${fileType}`);
  req.pipe(fileStream).on('finish', () => res.end('ok'))
}

server.listen(PORT)