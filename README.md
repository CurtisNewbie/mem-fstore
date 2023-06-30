# In-memory File Store

It's useful when you don't care about security, and you only want to transfer files between devices for like 10 minutes :D

An index.html is embedded in source code, you can access it on `/` path.

# Example

Lets say the IP address of this app is `192.168.31.153`.

To run the server:

```sh
./run.sh
```

## Upload File

```sh
curl 'http://192.168.31.153:8080/file/myfile.jpg' --data-binary @myfile.jpg; echo;
```

## Download File

```sh
curl 'http://192.168.31.153:8080/file/myfile.jpg' -o myfile.jpg
```

You can also type `http://192.168.31.153:8080/file/myfile.jpg` in your browser without using cURL.