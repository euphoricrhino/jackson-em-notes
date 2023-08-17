air=[1.0108, 1.0218, 1.0333, 1.0439, 1.0548];
airTransformed = arrayfun(@(x) (x-1)/(x+2), air);
airPressure=[20,40,60,80,100];
airFit = polyfit(airPressure, airTransformed, 1);

airX = linspace(0, 120, 40);
airY = polyval(airFit, airX);
subplot(1,3,1);
scatter(airPressure, airTransformed);
hold on;
plot(airX, airY); legend('air');


pentane=[1.82, 1.96, 2.12, 2.24, 2.33];
pentaneDensity=[.613,.701,.796, .865, .907];
pentaneTransformed1 = arrayfun(@(x) (x-1)/(x+2), pentane);
pentaneTransformed2 = arrayfun(@(x) (x-1), pentane);
pentaneFit1 = polyfit(pentaneDensity, pentaneTransformed1, 1);
pentaneFit2 = polyfit(pentaneDensity, pentaneTransformed2, 1);

pentaneX = linspace(0, 1.2, 40);
pentaneY1 = polyval(pentaneFit1, pentaneX);
pentaneY2 = polyval(pentaneFit2, pentaneX);

subplot(1,3,2);
hold on;
scatter(pentaneDensity, pentaneTransformed1);
plot(pentaneX, pentaneY1);
legend('pentane: e-1/e+2');

subplot(1,3,3);
hold on;
scatter(pentaneDensity, pentaneTransformed2);
plot(pentaneX, pentaneY2);
legend('pentane: e-1');

