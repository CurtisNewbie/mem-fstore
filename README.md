# In-Memory File Store

It's useful when you don't care about security, and you only want to transfer files between devices for like 10 minutes :D

An index.html is embedded in source code, you can access it on `/` path.

# Example

To upload file:

```sh
curl 'http://localhost:80/file?name=myfile.jpg' --data-binary @myfile.jpg; echo;
```

To download file:

```sh
curl 'http://localhost:80/file?name=myfile.jpg' -o myfile.jpg
```

You can also type `http://localhost:80/file?name=myfile.jpg` in your browser without using cURL.