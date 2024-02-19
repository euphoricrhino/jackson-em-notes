# Misc scripts to generate animation using asymptote

## Steps to generate mp4 from asy
* `asy total-reflection.asy`
* `magick -density 640 frame-*.pdf -alpha remove  frame-%02d.png`
* `ffmpeg -framerate 20 -i ./frame-%02d.png -c:v libx264 -vf "loop=loop=5:size=30:start=0,scale=iw/2:ih/2,pad=ceil(iw/2)*2:ceil(ih/2)*2" -profile:v high -crf 10 -pix_fmt yuv420p -y em.mp4`


https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/b3d9a9d1-dec4-477e-b672-8f68cea362fd

