alphas = [0.001:.001:5];
alphas2 = [2:.001:5];
function ret = potential(alpha)
  ret = power(alpha*pi, -1.5) * exp(-1/alpha);
endfunction

a = arrayfun(@(x) potential(x), alphas);
asym = arrayfun(@(x) power(x*pi,-1.5), alphas2);

plot(alphas, a, 'LineWidth', 2, alphas2, asym, 'LineWidth', 2);
xlabel('\alpha/|x|^2');
h = legend('|A|/|A_0|', '\alpha^{-3/2}');
set(gca, 'fontsize', 20);
set(h, "fontsize", 20);
