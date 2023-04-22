tx = ty = linspace(0, 1, 61);
[xx, yy] = meshgrid(tx, ty);

function ret = trial(x, y)
  ret = x * (1 - x) * y * (1 - y) * 5 / 4;
endfunction

trialPotential = arrayfun(@(x, y) trial(x, y), xx, yy);

function ret = seriesItem(x, y, m)
  tmp = 2 * m + 1;
  ret = sin(tmp * pi * x) / power(tmp, 3);
  v = cosh(tmp * pi * (y - 0.5)) / cosh(tmp * pi / 2);
  ret *= (1 - v);
  ret *= 4 / power(pi, 3);
endfunction

function ret = series(maxM, x, y)
  ret = 0;
  for m = 0 : maxM
    ret += seriesItem(x, y, m);
  endfor
endfunction

seriesPotential = arrayfun(@(x, y) series(100, x, y), xx, yy);

subplot(2, 2, 1);
mesh(tx, ty, trialPotential);
title('trial potential, A=5/(4\epsilon_0)');

trial1 = arrayfun(@(x) trial(x, 0.25), tx);
series1 = arrayfun(@(x) series(100, x, 0.25), tx);
subplot(2, 2, 2);
plot(tx, trial1, "--", tx, series1, "-");
legend("trial", "series");
title('trial vs series, y=0.25');

subplot(2, 2, 3);
mesh(tx, ty, seriesPotential);
title('series potential');

trial2 = arrayfun(@(x) trial(x, 0.5), tx);
series2 = arrayfun(@(x) series(100, x, 0.5), tx);
subplot(2, 2, 4);
plot(tx, trial2, "--", tx, series2, "-");
title('trial vs series, y=0.5');
legend("trial", "series");
