# Program for rendering 3D E/H fields of spherical resonant cavity.

## Example commands - generating images
```
go run *.go --l=2 --m=0 --xln=9.09501 --mode="TE (lmn=202)" --out-dir=./frames/te-202
go run *.go --l=3 --m=1 --xln=20.12181 --mode="TE (lmn=315)" --out-dir=./frames/te-315
go run *.go --l=4 --m=3 --xln=15.03966 --mode="TE (lmn=433)" --out-dir=./frames/te-433
go run *.go --l=2 --m=2 --xln=17.10274 --mode="TM (lmn=225)" --out-dir=./frames/tm-225
go run *.go --l=3 --m=2 --xln=4.97342 --mode="TM (lmn=321)" --out-dir=./frames/tm-321
go run *.go --l=4 --m=1 --xln=9.96755 --mode="TM (lmn=412)" --out-dir=./frames/tm-412
```
## Example - stitching images into video
```
ffmpeg -framerate 21.38 -i ./frames/te-202/frame-%04d.png -c:v libx264  -profile:v high -crf 10 -pix_fmt yuv420p -y te-202.mp4
```

## Example images


![frame-0615](https://github.com/user-attachments/assets/a3a057f8-011c-4c7e-8f9e-3124e743c912)
![frame-0598](https://github.com/user-attachments/assets/44a95c91-fc5e-4618-aaf2-2ef58509336d)
![frame-0532](https://github.com/user-attachments/assets/9c735ea9-5b38-4d7d-8d5d-6741960f7e61)
