#Simple TensorFlow Inception client 

Send image file and get back labels only jpeg/jpg supported.

## Run on single image 

```
./tfclient -addr 104.155.124.88:9000 -image file.jpeg

```


## Run on group of images
```
for f in ~/Desktop/*.jpg; do ./tfclient -addr 104.155.124.88:9000 -image $f;done

```