# gcsupload

This is a simple CLI that allows you to upload a file to your Google Cloud Storage bucket.


## Authentication

Please visit https://cloud.google.com/docs/authentication/getting-started


## Installing it

```shell
go get -u github.com/GoogleCloudPlatform/golang-samples/storage/gcsupload
```


## Using it

```shell
gcsupload -project orijtech-161805 -source ~/Desktop/birthdayPic.jpg -bucket orijtech-gcs-test
```
which will give you a result like this
```shell
URL: https://storage.googleapis.com/orijtech-gcs-test/birthdayPic.jpg
Size: 865096
MD5: 5b6c7b4aed837e8ed0f9950564a10b32
```
