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

[![v2-659ef91f404a6b0806771618c7e6761c_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/616cf04e-5794-40ba-aa4c-b4d21f35f4a0)](https://youtu.be/SZbFrAs2Zfc)

## Example images

![v2-96357cc51795af20b93b33be9ced583e_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/116cfc48-60f2-4a57-b3e0-d274a84b79d9)
![v2-9b363455bbe5b914639029a51c0c339b_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/28204110-ec00-436f-b91e-5f935e772df4)
![v2-d60d24ffe753eb76ac6ead3ea50f2176_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/6a4fcd10-82ab-4780-8fdd-05e3edaac2dd)
![v2-827b130c23752126e7344a092f1e4489_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/7ed29cfe-cf12-4c37-a016-f03312823024)
![v2-9b4def1e0d1daafd5e91f77147c11473_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/13451a84-1466-41b3-a576-85abbc2cefd6)
![v2-abc27c6ab6490e4dfcc79ce79e80b8c9_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/4d9c5a03-5b03-4583-aca1-d8a65c0900c7)
![v2-e62661d6f243d193445ea51a654629b0_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/0be9b9ca-819d-4524-96a1-ad44d485582c)
![v2-106ac861bf9dbb67a5a0216cff9bc3ab_1440w](https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/106cff8b-768e-4bd5-a172-dbb4970e1573)
