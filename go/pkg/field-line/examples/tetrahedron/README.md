# Field line rendering for two positive and two negative charges arranged on the tetrahedron vertices.
## Example - how to run
```
go run main.go --output /tmp/tetra
```

This will generate a series of png files which can be stiched using ffmpeg to generate an mp4 or gif.
Example
```
ffmpeg -framerate 10 -i /tmp/tetra-%03d.png -c:v libx264 -profile:v high -crf 10 -pix_fmt yuv420p tetra.mp4
```


![tetra-angle-083](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/dcd92a50-a00a-436e-acf0-e19bdb1be61f)
