% This code reproduces Figure 2.11.
x = [0:0.01:1];

function ret = seriesSum(x, y, terms)
  ret = 0;
  for k = 0 : terms - 1
    n = 2 * k + 1;
    ret += exp(-n * pi * y) * sin(n * pi * x) * 4 / (n * pi);
  endfor
endfunction

% exact plot by 2.65
function ret = exact(x, y)
  v = sin(pi * x) / sinh(pi * y);
  ret = 2 * atan(v) / pi;
endfunction

ylo100 = arrayfun(@(x) seriesSum(x, 0.1, 100), x);
yloFirst = arrayfun(@(x) seriesSum(x, 0.1, 1), x);
yloExact = arrayfun(@(x) exact(x, 0.1), x);

yhi100 = arrayfun(@(x) seriesSum(x, 0.5, 100), x);
yhiFirst = arrayfun(@(x) seriesSum(x, 0.5, 1), x);
yhiExact = arrayfun(@(x) exact(x, 0.5), x);

plot(x, ylo100, ':', 'linewidth', 2, x, yloFirst, '--', x, yloExact, x, yhi100, ':', 'linewidth', 2, x, yhiFirst, '--', x, yhiExact);
xlabel('x/a');
ylabel('\Phi(x,y)/V');
legend('y/a=0.1, 100terms', 'y/a=0.1, 1 term', 'y/a=0.1, exact', 'y/a=0.5, 100 terms', 'y/a=0.5, 1 term', 'y/a=0.5, exact')
