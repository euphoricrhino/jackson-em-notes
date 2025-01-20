# Program for rendering Mie Scattering by spherical dielectric sphere.

## Example commands - generating data
```
go run *.go --max-l 42 --n 1.33
```
## Example command - render the generated data
```
cd render
go run main.go --heatmap-file ../../heatmaps/wikipedia.png --output ./mie-scattered --data-file=../mie-scattered --count 376 --gamma=.5 --width 800 --height 800
```
## Example - stitching images into video
```
ffmpeg -framerate 20.75 -i ./mie-scattered-%03d.png -c:v libx264  -profile:v high -crf 10 -pix_fmt yuv420p -y mie-scattered.mp4
```

## Example images

R=3.18λ, n=1.33
![mie-total-293](https://github.com/user-attachments/assets/f7c4af60-bc52-4d86-b340-ef1a1fc6215e)
