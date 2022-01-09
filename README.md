# Imaginarium
A simple golang image storage engine. Used to create and store different sizes/thumbnails of user uploaded images.

## Description
**Imaginarium** enables you to create copies (or thumbnails) of your images and stores
them along with the original image on your filesystem. The image and its
copies are stored in a file structure based on the MD5 checksum of the original image.
The first character of the checksum used as the lvl 1 directory name.

**Imaginarium** supports png, jpg formats and also images encoded with base64. The decoder for any given image is chosen by the image's mimetype.

## Config file: /etc/imaginarium/config.yml
```yaml
---
storage:
  path: /mnt/storage

exporter:
  enabled: true
  port: 9360

server:
  port: 81
  uploader:
    maxSize: 5M
    contexts:
      - context: thumbnal
        width: 320
        height: 0
      - context: article
        width: 640
        height: 0
    allow:
      - image/png
      - image/jpg
      - image/jpeg

```

## Docker
```shell
# docker run -p 81:81 -p 9360:9360 -v $(PWD)/storage:/mnt/storage -v $(PWD)/config.yml:/etc/imaginarium/config.yml:ro -it vasary/imaginarium:latest 
```

## Usage
```curl
# curl --request POST --form "image=@1.png" http://127.0.0.1:81/upload

HTTP/1.0 201 Created
Content-Type: application/json; charset=UTF-8
X-Request-Id: hGGBnNMrq0jFohotl4ksvG18AB7sgpeJ
Date: Sun, 09 Jan 2022 19:06:06 GMT
Content-Length: 357

[
  {
    "1.png": [
      {
        "Name": "0afb00f5e1babc3aaa723368c44b618d.png",
        "Context": "original"
      },
      {
        "Name": "0afb00f5e1babc3aaa723368c44b618d_thumbnal.png",
        "Context": "thumbnal"
      },
      {
        "Name": "0afb00f5e1babc3aaa723368c44b618d_article.png",
        "Context": "article"
      }
    ]
  }
]