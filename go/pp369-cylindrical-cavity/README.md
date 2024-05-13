# Program for rendering 3D E/H fields of cylindrical resonant cavity.

## Example commands - generating images
```
go run main.go --p=0 --m=0 --xmn=2.405 --mode="TM (mnp=010)" --out-dir=./frames/tm-010
go run main.go --p=4 --m=3 --xmn=9.761 --mode="TM (mnp=324)" --out-dir=./frames/tm-324
go run main.go --p=2 --m=5 --xmn=18.9801 --mode="TM (mnp=542)" --out-dir=./frames/tm-542
go run main.go --p=1 --m=1 -xmn=1.841 --mode="TE (mnp=111)" --out-dir=./frames/te-111
go run main.go --p=3 --m=1 -xmn=5.331 --mode="TE (mnp=123)" --out-dir=./frames/te-123
go run main.go --p=6 --m=4 -xmn=19.196 --mode="TE (mnp=456)" --out-dir=./frames/te-456
```
## Example - stitching images into video
```
ffmpeg -framerate 20.5 -i ./frames/te-456/frame-%04d.png -c:v libx264  -profile:v high -crf 10 -pix_fmt yuv420p -y te-456.mp4
```
## Click the image below or full youtube video

[![frame-0599](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/8528345f-6875-4c63-b71b-b33d06f4238f)](https://youtu.be/rWC1cr8goaU)


## Example images

![frame-0584](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/6c6bcce1-8ef7-43d5-9f9b-aa8b1c359985)
![frame-0674](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/54b90c03-cd4f-428a-9a08-ef6d0107520a)
![frame-0584](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/b7817e11-f73c-4ecc-b5be-55e64c4084d7)
![frame-0543](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/597eaa97-e781-4e4a-9fd0-bccc31b91631)
![frame-0374](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/ae639019-23a8-4f51-ba69-6b6e9a9c1eeb)
![frame-0495](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/805996d8-852b-4695-97f2-a3eaaff6d136)
![frame-0554](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/6382e3fb-6a3c-4dd9-a91f-b34e04cb685e)
