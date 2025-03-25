function Lambda = boost_matrix(v)
  % boost_matrix(v): Computes the Lorentz boost matrix for a 3-velocity v (column vector)
  %                  Assumes c = 1 (natural units).
  %
  % Input:
  %   v - 3x1 velocity vector [vx; vy; vz], with norm(v) < 1
  %
  % Output:
  %   Lambda - 4x4 Lorentz boost matrix

  % Ensure v is a column vector
  v = v(:);
  v2 = dot(v, v);

  if v2 >= 1
    error("Speed must be less than the speed of light (|v| < 1).");
  endif

  gamma = 1 / sqrt(1 - v2);
  beta = v;
  beta_mag = sqrt(v2);

  % Initialize 4x4 matrix
  Lambda = eye(4);

  % Time-time component
  Lambda(1,1) = gamma;

  % Time-space components
  Lambda(1,2:4) = -gamma * beta';
  Lambda(2:4,1) = -gamma * beta;

  % Space-space components
  if beta_mag == 0
    Lambda(2:4,2:4) = eye(3);
  else
    beta_outer = beta * beta'; % Outer product
    Lambda(2:4,2:4) = eye(3) + (gamma - 1) / v2 * beta_outer;
  endif
endfunction
