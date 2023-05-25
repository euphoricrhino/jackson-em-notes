n = 320;
tx = [-n:1:n];
ty = [-n:1:n];
[xx, yy] = meshgrid(tx, ty);

function ret = potential(e0, e1, a, x, y)
  a2 = a * a;
  l = (y * y + x * x - a2) / a2;
  r = sqrt(l * l + 4 * y * y / a2);
  v1 = sqrt((r - l) / 2);
  v2 = abs(y) / a * atan(sqrt(2 / (r + l)));
  ret = (e0 - e1) * a / pi * (v1 - v2);
  if y > 0
    ret += e0 * y;
  else
    ret += e1 * y;
  endif
endfunction

E0 = 100;
E1 = 0;
rad = 50;
phi = @arrayfun(@(x, y) potential(E0, E1, rad, x, y), xx, yy);
lo = min(min(phi));
hi = max(max(phi));
normalized = (phi - lo) / (hi - lo);
subplot(1,2,1);
img = imagesc(tx, ty, normalized);
imadjust(img, stretchlim(img), []);
set(gca,'YDir','normal')
subplot(1,2,2);
contour(tx, ty, phi, 0:50:800);
