a=1;
l=5;
z=[-3:.01:3];

function ret = bfield(z,a,l)
  ret = (z-l/2) / sqrt(a*a+(l/2-z)*(l/2-z)) - (z+l/2)/sqrt(a*a+(l/2+z)*(l/2+z));
  ret /= -2;
endfunction

function ret = hfield(z,a,l)
  ret = bfield(z,a,l);
  if abs(z) < l/2
    ret -= 1;
  endif
endfunction

bval = arrayfun(@(x) bfield(x*l,a,l), z);
hval = arrayfun(@(x) hfield(x*l,a,l), z);

plot(z, bval,'linewidth', 1.5, z, hval, 'linewidth', 1.5);
xlabel('z/L');
set(gca, 'fontsize', 20);
set(gca, 'linewidth', 1.5);
h = legend('B/\mu_0M_0', 'H/M_0');
set(h, "fontsize", 20);
