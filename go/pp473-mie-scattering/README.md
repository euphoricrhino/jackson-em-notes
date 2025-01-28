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


![mie-total-293](https://github.com/user-attachments/assets/1b1abf1a-952c-489b-82df-75a0c746a005)

![mie-total-132](https://github.com/user-attachments/assets/21a019c2-9187-434c-a516-824d4677574c)

![mie-total-092](https://github.com/user-attachments/assets/95cd0f74-f542-4c2a-b3f7-f8e300a0c1e0)

![mie-total-136](https://github.com/user-attachments/assets/a1a97f74-4a4b-4e5a-80e3-5d481388bdae)
