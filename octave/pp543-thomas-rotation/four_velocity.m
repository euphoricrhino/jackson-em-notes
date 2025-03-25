function u = four_velocity(v)
  % four_velocity(v): Returns the 4-velocity corresponding to 3-velocity v
  % Assumes c = 1 (natural units)
  %
  % Input:
  %   v - 3x1 vector (3-velocity)
  % Output:
  %   u - 4x1 vector (4-velocity)

  v = v(:); % ensure column vector
  v2 = dot(v, v);

  if v2 >= 1
    error("Speed must be less than the speed of light (|v| < 1).");
  endif

  gamma = 1 / sqrt(1 - v2);
  u = gamma * [1; v];
endfunction
