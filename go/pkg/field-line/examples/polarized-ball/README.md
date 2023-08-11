# Field line rendering for polarized ball in constant field (Jackson example on page 158).
## Example - how to run
```
go run main.go --output /tmp/polball
```

This will generate a series of png files which can be stiched using ffmpeg to generate an mp4 or gif.
Example
```
ffmpeg -framerate 10 -i /tmp/polball-%03d.png -c:v libx264 -profile:v high -crf 10 -pix_fmt yuv420p polball.mp4
```


![polball-094](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/106d181d-f512-4a69-a286-ff23f343b10b)
