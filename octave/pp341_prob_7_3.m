n = 1.5;
% gap distance in unit of wave lengths.
d = [0:.01:10];

function t = transCoeffPerp(n, cosalpha, coseta, d)
  xi = 2 * pi * coseta * d;
  a = n * n * cosalpha * cosalpha + coseta * coseta;
  a *= a;
  sinhxi = sinh(xi);
  a *= sinhxi * sinhxi;
  b = 4 * n * n * cosalpha * cosalpha * coseta * coseta;
  t = b / (a + b);
endfunction

function t = transCoeffPara(n, cosalpha, coseta, d)
  xi = 2 * pi * coseta * d;
  a = n * n * coseta * coseta + cosalpha * cosalpha;
  a *= a;
  sinhxi = sinh(xi);
  a *= sinhxi * sinhxi;
  b = 4 * n * n * cosalpha * cosalpha * coseta * coseta;
  t = b / (a + b);
endfunction

% critical angle is 41.8deg for n=1.5.
alpha1 = 42.0 * pi / 180;
cosalpha1 = cos(alpha1);
sinalpha1 = sin(alpha1);
coseta1 = sqrt(n * n * sinalpha1 * sinalpha1 - 1);

perp1 = arrayfun(@(x) transCoeffPerp(n, cosalpha1, coseta1, x), d);
para1 = arrayfun(@(x) transCoeffPara(n, cosalpha1, coseta1, x), d);

alpha2 = 50.0 * pi / 180;
cosalpha2 = cos(alpha2);
sinalpha2 = sin(alpha2);
coseta2 = sqrt(n * n * sinalpha2 * sinalpha2 - 1);
perp2 = arrayfun(@(x) transCoeffPerp(n, cosalpha2, coseta2, x), d);
para2 = arrayfun(@(x) transCoeffPara(n, cosalpha2, coseta2, x), d);

alpha3 = 41 * pi / 180;
cosalpha3 = cos(alpha3);
sinalpha3 = sin(alpha3);
coseta3 = sqrt(n * n * sinalpha3 * sinalpha3 - 1);
perp3 = arrayfun(@(x) transCoeffPerp(n, cosalpha3, coseta3, x), d);
para3 = arrayfun(@(x) transCoeffPara(n, cosalpha3, coseta3, x), d);

subplot(2, 1, 1);
plot(d, perp1, 'linewidth', 1.5, d, perp2, 'linewidth', 1.5, d, perp3, 'linewidth', 1.5);
title('perpendicular polarization, critical angle=41.8\deg, n=1.5');
set(gca, 'fontsize', 20);
xlabel("d/λ");
ylabel("T");
h = legend('𝛼=42\deg', '𝛼=50\deg', '𝛼=41\deg');
set(h, "fontsize", 20);
hold on;

subplot(2, 1, 2);
plot(d, para1, 'linewidth', 1.5, d, para2, 'linewidth', 1.5, d, para3, 'linewidth', 1.5);
title('parallel polarization, critical angle=41.8\deg, n=1.5');
set(gca, 'fontsize', 20);
xlabel("d/λ");
ylabel("T");
h = legend('𝛼=42\deg', '𝛼=50\deg', '𝛼=41\deg');
set(h, "fontsize", 20);
