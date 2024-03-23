import graph;
settings.outformat = "pdf";

// Commands to generate gif.
// 1. asy packet-dispersion.asy
// 2. magick -density 480 wide-*.pdf -alpha remove  wide-%04d.png
// 3. ffmpeg -i wide-%04d.png -vf fps=20 wide.gif

// Wide initial width
// string prefix = "wide";
// real k0 = 103.79;
// real L = .1;
// real nu = 2.0;
// real a = 1.0 / k0;
// real vscale = 1.0;

// Narrow initial width
//string prefix = "narrow";
//real k0 = 8;
//real L = .1;
//real nu = .25;
//real a = .5 / k0;
//real vscale = 2.0;

// Very narrow initial width
string prefix = "very-narrow";
real k0 = 2;
real L = .1;
real nu = .25;
real a = .2 / k0;
real vscale = 4.0;

// Square root of a complex number.
pair sqrtz(pair z) {
  real r = sqrt(z.x * z.x + z.y * z.y);
  real theta = acos(z.x / r);
  real sqrtr = sqrt(r);
  return (sqrtr * cos(theta/2), sqrtr * sin(theta/2));
}

// exp(z) for z complex.
pair expz(pair z) {
  real r = exp(z.x);
  return (r * cos(z.y), r * sin(z.y));
}

// Real divided by complex.
pair divz(real a, pair z) {
  real r = z.x * z.x + z.y * z.y;
  return (a * z.x / r, -a * z.y / r);
}

// Complex multiplied by complex.
pair zmulz(pair z1, pair z2) {
  return (z1.x * z2.x - z1.y * z2.y, z1.x * z2.y + z1.y * z2.x);
}

pair zdivz(pair z1, pair z2) {
  real r = z2.x * z2.x + z2.y * z2.y;
  z1 = (z1.x / r, z1.y / r);
  return zmulz(z1, (z2.x, -z2.y));
}

// Evaluation of (7.98).
real uk0(real x, real t, real kk0) {
  real v = x - nu * a * a * kk0 * t;
  v *= -v;
  pair z1 = (2 * L * L, 2.0 * a * a * nu * t);
  pair env = expz(divz(v, z1));
  pair wave = expz((0, kk0 * x - nu * t * (1.0 + a * a * kk0 * kk0 / 2)));
  pair den = sqrtz((1, a * a * nu * t / (L * L)));
  pair w = zdivz(zmulz(env, wave), den);
  return w.x / 2.0;
}

real Lt(real t) {
  real v = a * a * nu * t / L;
  return sqrt(L * L + v * v);
}

typedef real rr(real);

rr func(real t) {
  return new real (real x) {
    return vscale * (uk0(x, t, k0) + uk0(x, t, -k0));
  };
}

real range(real t) {
  real center = nu * a * a * k0 * t;
  real lt = 5 * Lt(t);
  return center + lt;
}

int frames = 150;
real boxRange = range(frames-1);
for (int f = 0; f < frames; ++f) {
  picture pic;
  size(pic, 8cm, 8cm);
  unitsize(pic, 4cm);

  draw(pic, box((-boxRange,-1.2*vscale), (boxRange, 1.2*vscale)), opacity(0));
  real ranget = range(f);
  draw(pic, graph(func(f), -ranget, ranget, 1000), red);
  shipout(format(prefix+"-%04d", f), scale(4) * pic);
}
