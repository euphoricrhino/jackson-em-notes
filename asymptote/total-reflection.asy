settings.outformat = "pdf";
settings.render = 0;
import three;

currentprojection = perspective((-.5,-.4,.3), up=Z);

real n = 1.50;
real nn = 1;
real iangle = radians(90);

int xsamples = 120;
int xstart = 20;
real DE = 12;
real xstep = DE / xsamples;
real lambda = 6;
real kr = 2 * pi / lambda;
real si = sin(iangle);
real ezamp = n / nn * si;
real examp = sqrt(n * n * si * si / nn / nn - 1);
real bamp = .6;
int frames = 30;
int tailNodes = 18;

triple[] ende = new triple[];
triple[] endb = new triple[];
for (int f = 0; f < frames; ++f) {
  real omegat = f * 2  pi / frames;
  triple epos = (examp * sin(-omegat), 0, ezamp * cos(-omegat));
  ende.push(epos);
  triple bpos = (0, -bamp * cos(-omegat), 0);
  endb.push(bpos);
}

for (int f = 0; f < frames; ++f) {
  real omegat = f * 2 * pi / frames;
  picture pic;
  size(pic, 8cm, 8cm);
  unitsize(pic, 2cm);

  draw(pic, O--2Z, mediumgray, arrow=Arrow3(emissive(mediumgray)));
  draw(pic, O--2Y, mediumgray, arrow=Arrow3(emissive(mediumgray)));
  draw(pic, O--15X, mediumgray, arrow=Arrow3(emissive(mediumgray)));

  path3 ge, gb;
  for (int i = 0; i <= xsamples; ++i) {
    real x = xstep * i;
    real ez = ezamp * cos(kr * x - omegat);
    real ex = examp * sin(kr * x - omegat);
    triple epos = (x + ex, 0, ez);
    triple bpos = (x, -bamp * cos(kr * x - omegat), 0);
    ge = ge -- epos;
    gb = gb -- bpos;
    if (i >= xstart) {
      draw(pic, (x, 0, 0) -- epos, red);
      draw(pic, (x, 0, 0) -- (bpos), blue);
    }
    if (i == 0) {
      for (int t = 0; t <= tailNodes; ++t) {
        real linewidth = 1-.5*t/tailNodes;
        pen pene = rgb(1, t/tailNodes, t/tailNodes);
        pen penb = rgb(t/tailNodes, t/tailNodes, 1);
        draw(pic, ende[(f-t+frames)%frames] -- ende[(f-t-1+frames)%frames], pene+linewidth);
        draw(pic, endb[(f-t+frames)%frames] -- endb[(f-t-1+frames)%frames], penb+linewidth);
      }
      draw(pic, (x, 0, 0) -- (epos), red, arrow=Arrow3(DefaultHead2(normal=Y)));
      draw(pic, (x, 0, 0) -- (bpos), blue, arrow=Arrow3(DefaultHead2(normal=Y)));
    }
  }
  draw(pic, ge, palered+linewidth(1pt));
  draw(pic, gb, paleblue+linewidth(1pt));
  draw(pic, box((0,-2,-2), (15.5,3.5,2.5)), opacity(0));
  label(pic, "$x$", (15.6, 0, 0));
  label(pic, "$y$", (0, 2.3, 0));
  label(pic, "$z$", (0, 0, 2.3));

  shipout(format("frame-%02d", f), scale(4)*pic);
}
