0;

% integral of int_0^1 rho^k(1-\rho)^l drho.
function ret = intkl(k, l)
  ret = 0;
  sgn = 1;
  for m = 0 : l
    ret += sgn * nchoosek(l, m) / (k + m + 1);
    sgn *= -1;
  endfor
endfunction 

% g(rho) = -5(1-rho)+10^4 rho^5(1-rho)^5.
function ret = g(rho)
  ret = power(rho, 5);
  ret -= 5 * power(rho, 6);
  ret += 10 * power(rho, 7);
  ret -= 10 * power(rho, 8);
  ret += 5 * power(rho, 9);
  ret -= power(rho, 10);
  ret *= 1e4;
  ret += 5 * rho - 5;
endfunction

% solve for psi_1 using variational method.
vec1 = [1e4 * intkl(6, 6) - 5 * intkl(1, 2); 1e4 * intkl(6, 7) - 5 * intkl(1, 3); 1e4 * intkl(6, 8) - 5 * intkl(1, 4)];
A1 = [
  intkl(1, 0), 2 * intkl(1, 1), 3 * intkl(1, 2)
  2 * intkl(1, 1), 4 * intkl(1, 2), 6 * intkl(1, 3)
  3 * intkl(1, 2), 6 * intkl(1, 3), 9 * intkl(1, 4)
];
invA1 = inv(A1);
% alpha_1, beta_1, gamma_1.
psi1Params = invA1 * vec1;

% psi_1(rho) = alpha_1 (1-rho) + beta_1(1-rho)^2 + gamma_1(1-rho)^3.
function ret = psi1(rho, psi1Params)
  ret = psi1Params(1) * (1 - rho);
  ret += psi1Params(2) * power(1 - rho, 2);
  ret += psi1Params(3) * power(1 - rho, 3);
endfunction

% e(n) = int_0^1 g(rho^n-1)rho drho
function ret = e(n)
  ret = -5 * intkl(n + 1, 1);
  ret += 5 * intkl(1, 1);
  ret += 1e4 * intkl(n + 6, 5);
  ret -= 1e4 * intkl(6, 5);
endfunction

A2 = [
  1 6/5 4/3
  6/5 3/2 12/7
  4/3 12/7 2
];

invA2 = inv(A2);

% alpha, beta, gamma.
vec2 = [e(2); e(3); e(4)];
psi2Params = invA2 * vec2;

% psi_2(rho) = alpha rho^2+beta rho^3+gamma rho^4-(alpha+beta+gamma)
function ret = psi2(rho, psi2Params)
  ret = psi2Params(1) * rho * rho;
  ret += psi2Params(2) * power(rho, 3);
  ret += psi2Params(3) * power(rho, 4);
  ret -= psi2Params(1) + psi2Params(2) + psi2Params(3);
endfunction

% indefinite integral of exact solution psi_E.
function ret = indefExact(rho)
  ret = power(rho, 7) / 49;
  ret -= power(rho, 8) * 5 / 64;
  ret += power(rho, 9) * 10 / 81;
  ret -= power(rho, 10) / 10;
  ret += power(rho, 11) * 5 / 121;
  ret -= power(rho, 12) / 144;
  ret *= 1e4;
  ret -= power(rho, 2) * 5 / 4;
  ret += power(rho, 3) * 5 / 9;
endfunction

% psi_E with translation so psi_E(rho=1)=0.
function ret = exact(rho)
  rhoAt1 = indefExact(1);
  ret = indefExact(rho) - rhoAt1;
  ret = -ret;
endfunction

rho = [0:0.01:1];
% arbitrary scale for g.
gVals = arrayfun(@(rho) 0.05 * g(rho), rho);
psi2Vals = arrayfun(@(rho) psi2(rho, psi2Params), rho);
psi1Vals = arrayfun(@(rho) psi1(rho, psi1Params), rho);
exactVals = arrayfun(@(rho) exact(rho), rho);

% plot everything in one figure.
plot(rho, psi1Vals, ':', 'LineWidth', 1.5, rho, psi2Vals, '--', 'LineWidth', 1.5, rho, gVals, '-.', 'LineWidth', 1.5, rho, exactVals, '-', 'LineWidth', 1.5);
set(gca, 'fontsize', 16);
h = legend('\Psi_1(\rho)', '\Psi_2(\rho)', 'g(\rho)', '\Psi_E(\rho)');
set(h, "fontsize", 16);
