import graph;
settings.outformat = "pdf";

// This script reproduces figure 8.12, but the drift curve of the graded fiber is a couple of orders of magnitude off compared to the text.

real a = sqrt(2)/10;

real launch1 = a;
real launch2 = a/2;

real step_n(real x) {
  return 1;
}

real gaussian_n(real x) {
  return exp(-x * x / 2);
}

real nbar(real launch) {
  return cos(launch);
}

real xmax_step = a;

real xmax_gaussian(real launch) {
  real ct = cos(launch);
  return sqrt(-2*log(ct));
}

void extend(pair[] zx, real sign) {
  real center = zx[zx.length-1].x;
  for (int i=zx.length-2; i>=0; --i) {
    zx.push((2*center-zx[i].x, sign*zx[i].y));
  }
}

typedef real rr(real);

void zx(real xmax, pair[] ret, rr nx, real nbar, int steps) {
  real v = 0.0;
  real dx = xmax / steps;
  for (int i=0; i<steps; ++i) {
    real x = dx * i;
    real n = nx(x);
    v += nbar/sqrt(n*n-nbar*nbar)*dx;
    ret.push((v, x));
  }

  extend(ret, 1);
  extend(ret, -1);
}

currentpicture = new picture;
size(12cm, 8cm, IgnoreAspect);

pair[] step_launch1 = new pair[];
pair[] gauss_launch1 = new pair[];

zx(xmax_step, step_launch1, step_n, nbar(launch1), 200);
zx(xmax_gaussian(launch1), gauss_launch1, gaussian_n, nbar(launch1), 200);

pair[] step_launch2 = new pair[];
pair[] gauss_launch2 = new pair[];

zx(xmax_step, step_launch2, step_n, nbar(launch2), 200);
zx(xmax_gaussian(launch2), gauss_launch2, gaussian_n, nbar(launch2), 200);

draw(graph(step_launch1), red);
draw(graph(gauss_launch1), blue);

draw(graph(step_launch2), red+dashed);
draw(graph(gauss_launch2), blue+dashed);
draw(graph(new pair[]{(0,a),(8,a)}), linewidth(2)+opacity(.2));
draw(graph(new pair[]{(0,0),(8,0)}), linewidth(2)+opacity(.2));
draw(graph(new pair[]{(0,-a),(8,-a)}), linewidth(2)+opacity(.2));
ylimits(-.2, .2);
xaxis("$z$", BottomTop, LeftTicks);
yaxis("$x$",LeftRight, RightTicks);
shipout("1");

real Z(real xmax, rr nx, real nbar, int steps) {
  real v = 0;
  real dx = xmax / steps;
  for (int i = 0; i < steps; ++i) {
    real x = dx * i;
    real n = nx(x);
    v += 2*nbar/sqrt(n*n-nbar*nbar)*dx;
  }
  return v;
}

real Lopt(real xmax, rr nx, real nbar, int steps) {
  real v = 0;
  real dx = xmax / steps;
  for (int i = 0; i < steps; ++i) {
    real x = dx * i;
    real n = nx(x);
    v += 2 * n * n / sqrt(n*n-nbar*nbar)*dx;
  }
  return v;
}

rr drift(rr xmax, rr nx, int steps) {
  return new real (real launch) {
    real xmax_val = xmax(launch);
    real nbar_val = nbar(launch);
    return Lopt(xmax_val, nx, nbar_val, steps) / Z(xmax_val, nx, nbar_val, steps) - 1;
  };
}

currentpicture = new picture;
size(12cm, 8cm, IgnoreAspect);

scale(Linear(true),Log(true));

real drifts(real launch) {
  return 1/cos(launch)-1;
}

draw(graph(drift(new real (real x) { return xmax_step; }, step_n, 200), 0.015, launch1, 200), red);
draw(graph(drift(xmax_gaussian, gaussian_n, 200), 0.025, launch1, 200), blue);
xlimits(0, a);
ylimits(1e-8, 1e-2);
xaxis("", BottomTop, LeftTicks);
yaxis("",LeftRight, RightTicks);
shipout("2");

