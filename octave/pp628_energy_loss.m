% Script to plot y = eta * (1 + 1/x^2) * ln(2 * x^2 * A / B) - eta on a log-log scale
% Parameters
eta = 0.15;
A = 0.511e6; % 0.511 * 10^6
B = 160;
epsilon = 1e4;

% Define x range for log-log plot (avoid x too small to ensure ln argument > 1)
x = logspace(-1, 4, 1000); % x from 0.1 to 10, 1000 points

% Compute the function
y = eta * (1 + 1./(x.^2)) .* log(2 * x.^2 * A / B) - eta;
y2 = .5*eta * (1+1./(x.^2)) .* log(2 * x.^2 * A * epsilon / B /B) - eta;

% Create log-log plot
figure;
loglog(x, y, 'b-', 'LineWidth', 1);
hold on;
loglog(x, y2, '--', 'LineWidth', 1);
grid on;
xlabel('x');
ylabel('y');
title('Log-Log Plot of y = \eta (1 + 1/x^2) ln(2 x^2 A / B) - \eta');
