omega = [0:.01:2*pi];

function r = rval(n1, n2, n3, omega)
  arg = 2 * n2 * omega;
  r = cos(arg) + i * sin(arg);
  r *= (n2 - n3) / (n2 + n3);
endfunction

function ret = reflCoeff(n1, n2, n3, omega)
  r = rval(n1, n2, n3, omega);
  ret = (1 + r) * n1 - (1 - r) * n2;
  ret /= (1 + r) * n1 + (1 - r) * n2;
  ret = abs(ret);
  ret *= ret;
endfunction

function ret = transCoeff(n1, n2, n3, omega)
  r = rval(n1, n2, n3, omega);
  t = (1 + r) * n1 + (1 - r) * n2;
  t = abs(t);
  t *= t;
  ret = 2 * n2 / (n2 + n3);
  ret *= ret;
  ret *= 4 * n1 * n3 / t;
endfunction

rval = arrayfun(@(x) reflCoeff(1, 2, 3, x), omega);
tval = arrayfun(@(x) transCoeff(1, 2, 3, x), omega);
subplot(3, 1, 1);
plot(omega, rval, 'linewidth', 1.5, 'color', '#3595f6', omega, tval, 'linewidth', 1.5, 'color', '#45b635', omega, rval+tval, 'linewidth', 1.5, 'color', '#f63595');
title("n_1=1, n_2=2, n_3=3")
set(gca, 'fontsize', 20);
h = legend('R', 'T', 'R+T');
xlabel("ω/ω_0");
set(h, "fontsize", 20);
hold on;

rval = arrayfun(@(x) reflCoeff(3, 2, 1, x), omega);
tval = arrayfun(@(x) transCoeff(3, 2, 1, x), omega);
subplot(3, 1, 2);
plot(omega, rval, 'linewidth', 1.5, 'color', '#3595f6', omega, tval, 'linewidth', 1.5, 'color', '#45b635', omega, rval+tval, 'linewidth', 1.5, 'color', '#f63595');
title("n_1=3, n_2=2, n_3=1")
set(gca, 'fontsize', 20);
h = legend('R', 'T', 'R+T');
xlabel("ω/ω_0");
set(h, "fontsize", 20);
hold on;

rval = arrayfun(@(x) reflCoeff(2, 4, 1, x), omega);
tval = arrayfun(@(x) transCoeff(2, 4, 1, x), omega);
subplot(3, 1, 3);
plot(omega, rval, 'linewidth', 1.5, 'color', '#3595f6', omega, tval, 'linewidth', 1.5, 'color', '#45b635', omega, rval+tval, 'linewidth', 1.5, 'color', '#f63595');
title("n_1=2, n_2=4, n_3=1")
set(gca, 'fontsize', 20);
h = legend('R', 'T', 'R+T');
xlabel("ω/ω_0");
set(h, "fontsize", 20);

