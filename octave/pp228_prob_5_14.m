mu = logspace(-2,6,400);

function ret = ratio(x, lambda)
  ret = 4 * x;
  ret /= (x+1)*(x+1)-(x-1)*(x-1)*lambda;
endfunction

y1 = arrayfun(@(x) ratio(x, 0.5), mu);
y2 = arrayfun(@(x) ratio(x, 0.1), mu);

semilogx(mu, y1, 'LineWidth', 1.5); hold on;
semilogx(mu, y2, 'LineWidth', 1.5); hold on;
xlabel('\mu_r');
ylabel('|B|/B_0');
set(gca, 'fontsize', 16);
h = legend('(a/b)^2=0.5','(a/b)^2=0.1');
set(h, "fontsize", 16);
