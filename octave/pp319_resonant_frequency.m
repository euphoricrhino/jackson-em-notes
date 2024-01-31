om = [.5, 1.3, 2.6, 5.7, 6.2, 7.8, 9.7];
gamma = [om(1) * .05, om(2) * .07, om(3) * .12, om(4) * .03, om(5) * .08, om(6) * .15, om(7) * .11];
f = [132, 498, 123, 699, 359, 439, 392];
omega = [0:.01:12];

function [re, im] = epsilon(om, gamma, f, omega)
  re = [];
  im = [];
  for k = 1 : length(omega)
    x = omega(k);
    v = 0;
    for j = 1 : length(om)
      z = om(j) * om(j) - x * x - i * x * gamma(j);
      v += f(j) / z;
    endfor
    re = [re, real(v)];
    im = [im, imag(v)];
  endfor
endfunction

[res, ims] = epsilon(om, gamma, f, omega);

subplot(2, 1, 1);
reh = plot(omega, res, 'linewidth', 1.5);
title('Re(\epsilon) ~ \omega');
set(gca, 'xtick', om);
hold on;
subplot(2, 1, 2);
imh = plot(omega, ims, 'linewidth', 1.5);
title('Im(\epsilon) ~ \omega');
set(gca, 'xtick', om);
