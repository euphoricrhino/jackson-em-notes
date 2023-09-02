# Field line rendering for the magnetic field of Jackson problem 5-13.
## Example - how to run
```
go run main.go --output /tmp/sphere
```

This will generate a series of png files which can be stiched using ffmpeg to generate an mp4 or gif.
Example
```
ffmpeg -framerate 10 -i /tmp/sphere-%03d.png -c:v libx264 -profile:v high -crf 10 -pix_fmt yuv420p sphere.mp4
```
![sphere](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/ce792c32-f3b4-40d4-8589-63227bafe899)
