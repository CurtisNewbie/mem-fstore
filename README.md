# In-memory File Store

It's useful when you don't care about security, and you only want to transfer files between devices for like 10 minutes :D

# Example

Lets say the IP address of this app is `192.168.31.153`.

## Upload File

```sh
curl -X PUT 'http://192.168.31.153:8080/file' --data-binary @myfile.jpg -H 'File:myfile.jpg'
```

## Download File

```sh
curl -v http://192.168.31.153:8080/file?file=myfile.jpg -o myfile.jpg
```

You can also type `http://192.168.31.153:8080/file?file=myfile.jpg` in your browser.