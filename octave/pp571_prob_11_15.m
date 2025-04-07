addpath("./pp543-thomas-rotation/");

% generate random alpha=|B|/|E| in [-3,3].
alpha = 3 * rand;

% generate random angle theta in [0, pi].
theta = pi * rand;

% E, B in K frame.
E = [1; 0; 0];
B = [alpha*cos(theta); alpha*sin(theta); 0];

% compute beta by equation (18) with minus sign.
beta = ((1 + alpha^2) - sqrt((1 + alpha^2)^2 - 4 * alpha^2 * sin(theta)^2)) / (2 * alpha * sin(theta));
gamma_v = 1 / sqrt(1 - beta^2);
v = beta * [0; 0; 1];

% transform E/B via B(v) to K' frame
Ep = gamma_v * (E + cross(v, B)) - gamma_v^2 * v * dot(v, E) / (gamma_v + 1);
Bp = gamma_v * (B - cross(v, E)) - gamma_v^2 * v * dot(v, B) / (gamma_v + 1);
printf("E'xB':\n");
cross(Ep, Bp)


% unit transverse direction of E'||B' in K' frame (equation (15), (16)).
unit_rho = [1-alpha * beta * sin(theta); alpha * beta * cos(theta); 0];
unit_rho = unit_rho / norm(unit_rho);

% generate a random transverse component of the u-velocity.
u_trans = sqrt((1-beta^2)*rand);

u = u_trans * unit_rho + v;

gamma_u = 1 / sqrt(1 - dot(u,u));

% transform E/B via B(u) to K'' frame
Epp = gamma_u * (E + cross(u, B)) - gamma_u^2*u*dot(u, E)/(gamma_u+1);
Bpp = gamma_u * (B - cross(u, E)) - gamma_u^2*u*dot(u, B)/(gamma_u+1);
printf("E''xB'':\n");
cross(Epp, Bpp)

% compute the Thomas rotation matrix r.
axis = cross(-v, u);
axis = axis / norm(axis);

gamma = gamma_v * gamma_u * (1 + dot(-v, u));
angle = acos((1 + gamma + gamma_v + gamma_u)^2 / (gamma + 1) / (gamma_v + 1) / (gamma_u + 1) - 1);

r = axis_angle_rotation(axis, angle);

% verify that rE'=E'', rB'=B''.
printf("|rE'-E''|:\n")
norm(r * Ep - Epp)
printf("|rB'-B''|:\n")
norm(r * Bp - Bpp)
