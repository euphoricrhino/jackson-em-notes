function v = add_velocity(v1, v2)
  % relativistic_velocity_add(v1, v2)
  % Computes the relativistic composition of two 3-velocities v1 and v2
  % Assumes c = 1
  %
  % Input:
  %   v1, v2 - 3x1 column vectors (velocities)
  % Output:
  %   v - resulting 3-velocity (3x1 column vector)

  v1 = v1(:);
  v2 = v2(:);

  v1_mag2 = dot(v1, v1);
  gamma1 = 1 / sqrt(1 - v1_mag2);

  dot_v1v2 = dot(v1, v2);

  denom = 1 + dot_v1v2;

  v_par = (dot_v1v2 / v1_mag2) * v1;

  v_perp = v2 - v_par;

  % Composite velocity using relativistic velocity addition formula
  v = (v1 + (1 / gamma1) * v_perp + v_par) / denom;

  % Safety check: ensure |v| < 1
  if norm(v) >= 1
    warning("Resulting speed exceeds the speed of light! |v| = %f", norm(v));
  endif
endfunction
