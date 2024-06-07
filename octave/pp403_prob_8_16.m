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

function ret = groupVelocity(vvt, p, n1, n2)
  delta = (n1 * n1 - n2 * n2) / (2 * n1 * n1);
  V = vvt * p * pi / 2;
  f = @(x) fun(V, p, x);
  df = @(x) der(V, p, x);
  ll = p * pi / (2 * V);
  lr = (p + 1) * 1.56 / V;
  % uncomment below to show whether the zero xi exists for different vvt/p values.
  %step = (lr - ll) / 100;
  %ldomain = [ll : step : lr];
  %rdomain = [.1 : .01 : 1];
  %l = arrayfun(@(x) lhs(V, p, x), ldomain);
  %r = arrayfun(@(x) rhs(V, p, x), rdomain);
  %plot(ldomain, l, 'linewidth', 1.5);
  %hold on;
  %plot(rdomain, r, 'linewidth', 1.5);
  rootRange = fzero(f, infsup(sprintf("[%f, %f]", ll, lr)), df);
  xi = (rootRange.inf + rootRange.sup) / 2;
  ct2 = 1 - 2 * delta * xi * xi;
  betaa = V * sqrt(1 - xi * xi);
  ret = sqrt(ct2) * (1 + betaa);
  ret /= n1 * (ct2 + betaa);
endfunction

vvt = [1 : .01 : 2];
vvt = [vvt, 2.1 : .1 : 10];
vg = arrayfun(@(x) groupVelocity(x, 1, 1.5, 1.0), vvt);
subplot(3, 1, 1);
plot(vvt, vg, 'linewidth', 1.5);
title('n_1=1.5,n_2=1,p=1');
set(gca, 'fontsize', 20);
xlabel("V/V_t");
ylabel("v_g/c");
vg = arrayfun(@(x) groupVelocity(x, 1, 1.01, 1.0), vvt);
subplot(3, 1, 2);
plot(vvt, vg, 'linewidth', 1.5);
title('n_1=1.01,n_2=1,p=1');
set(gca, 'fontsize', 20);
xlabel("V/V_t");
ylabel("v_g/c");

ps = [1 : 1 : 20];
vvg = arrayfun(@(p) groupVelocity(3, p, 1.01, 1.0), ps);
subplot(3, 1, 3);
plot(ps, vvg, 'linewidth', 1.5);
title('V/V_t=3,n_1=1.01,n_2=1');
set(gca, 'fontsize', 20);
xlabel("p");
ylabel("v_g/c");
