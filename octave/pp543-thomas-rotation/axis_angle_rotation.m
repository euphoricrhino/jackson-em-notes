function R = axis_angle_rotation(axis, theta)
  % rotation_matrix_axis_angle Returns the 3x3 rotation matrix
  %  that rotates a vector by angle theta about the specified axis.
  %
  %  Inputs:
  %    axis : 3x1 vector defining the axis of rotation (can be any non-zero length)
  %    theta: rotation angle in radians
  %
  %  Output:
  %    R    : 3x3 rotation matrix

  % Ensure 'axis' is a column vector
  if size(axis,1) ~= 3
      axis = axis(:);
  end

  % Normalize the axis to have unit length
  axis = axis / norm(axis);

  % Extract normalized axis components
  ux = axis(1);
  uy = axis(2);
  uz = axis(3);

  % Precompute trigonometric values
  ct = cos(theta);
  st = sin(theta);
  vt = 1 - ct;

  % Rodrigues' rotation formula:
  % R = I + (sin(theta)) * K + (1 - cos(theta)) * K^2
  % where K is the skew-symmetric cross-product matrix of the unit axis [u].
  %
  % We can write out the expanded form of that directly:

  R = [ ct + ux^2*vt,     ux*uy*vt - uz*st, ux*uz*vt + uy*st;
        uy*ux*vt + uz*st, ct + uy^2*vt,     uy*uz*vt - ux*st;
        uz*ux*vt - uy*st, uz*uy*vt + ux*st, ct + uz^2*vt     ];
end
