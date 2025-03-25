function v = random_3velocity()
  % random_3velocity(): Generates a random 3-velocity vector with |v| < 1
  %
  % Output:
  %   v - 3x1 column vector with norm(v) < 1

  while true
    v = 2 * rand(3, 1) - 1;  % Random vector in [-1, 1]^3
    if norm(v) < 1
      break;
    endif
  endwhile
endfunction