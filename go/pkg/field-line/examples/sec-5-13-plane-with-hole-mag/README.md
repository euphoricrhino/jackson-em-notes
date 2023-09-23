# Field line rendering for the magnetic field of Conducting plane with hole in section 5.13
## Example - how to run
```
go run main.go --output /tmp/hole --step=0.001
```

This will generate a series of png files which can be stiched using ffmpeg to generate an mp4 or gif.
Example
```
ffmpeg -framerate 10 -i /tmp/hole-%03d.png -c:v libx264 -profile:v high -crf 10 -pix_fmt yuv420p hole.mp4
```

![hole-010](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/f044470a-e0bb-4bac-87b9-c29bf7ece770)
