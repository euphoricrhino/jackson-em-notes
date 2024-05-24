pkg load interval;

function ret = lhs(V, p, x)
  ret = tan(V * x - p * pi / 2);
endfunction

function ret = rhs(V, p, x)
  ret = sqrt(1 / (x * x) - 1);
endfunction

function ret = fun(V, p, x)
  ret = lhs(V, p, x) - rhs(V, p, x);
endfunction

function ret = der(V, p, x)
  arg = V * x - p * pi / 2;
  c = cos(arg);
  ret = V / (c * c);
  ret += 1 / (x * x) / sqrt(1 - x * x);
endfunction

function ret = approx(V, p)
  ret = (p + 1) * pi / (2*(V + 1));
  ret *= (1 - power((p + 1) * pi, 2) / (24 * power(V + 1, 3)));
endfunction

function run(V, p, i, j, k)
  f = @(x) fun(V, p, x);
  df = @(x) der(V, p, x);
  ll = p * pi / (2 * V);
  lr = (p + 1) * 1.54 / V;
  step = (lr - ll) / 100;
  ldomain = [ll : step : lr];
  rdomain = [.1 : .01 : 1];
  l = arrayfun(@(x) lhs(V, p, x), ldomain);
  r =arrayfun(@(x) rhs(V, p, x), rdomain);
  subplot(i, j, k)
  plot(ldomain, l, 'linewidth', 1.5);
  hold on;
  plot(rdomain, r, 'linewidth', 1.5);
  title(sprintf("V=%d,p=%d", V, p));
  set(gca, 'fontsize', 20);
  set(gca, 'linewidth', 1.5);
  rootRange = fzero(f, infsup(sprintf("[%f, %f]", ll, lr)), df);
  numerical = (rootRange.inf + rootRange.sup) / 2;
  approximate = approx(V, p);
  printf("V=%d, p=%d, numerical=%f, approximated root=%f\n", V, p, numerical, approximate);
endfunction


run(1, 0, 2, 2, 1);
run(2, 0, 2, 2, 2);
run(3, 1, 2, 2, 3);
run(10, 3, 2, 2, 4);

