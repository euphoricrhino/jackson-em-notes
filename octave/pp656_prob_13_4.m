mProton = 938e6;
mMeson = 105.7e6;

function n = N(rho, M, Z)
  n = Z * rho * 6.02e23 / M;
endfunction

function bq = Bq(Z,gamma,beta)
  bq = 2 * gamma * gamma * beta * beta * .511e6 / Z/12;
endfunction

function dedx = dEdx(T, M, Z, N)
  gamma = T/M+1;
  beta = sqrt(1-1/gamma^2);
  eta = 4*pi * N * Z*(1.44e-13)^2/.511e6;
  dedx = eta / beta / beta * (log(Bq(Z,gamma,beta))) - eta;
endfunction

function dedxrho = dEdx_rho(T, M, Z)
  gamma = T/M+1;
  beta = sqrt(1-1/gamma^2);
  eta = .15;
  dedxrho = eta / beta / beta * (log(Bq(Z,gamma,beta))) - eta;
endfunction

rho_al = 2.7;
rho_cu = 8.96;
rho_pb = 11.34;

M_al = 26.98;
M_cu = 63.55;
M_pb = 207.2;

Z_al = 13;
Z_cu = 29;
Z_pb = 82;

N_al = N(rho_al, M_al, Z_al)
N_cu = N(rho_cu, M_cu, Z_cu)
N_pb = N(rho_pb, M_pb, Z_pb)

N_air = 3.51e20;
Z_air = 14.4;


air_proton_10M = dEdx(10e6, mProton, Z_air, N_air)
al_proton_10M = dEdx(10e6, mProton, Z_al, N_al)
cu_proton_10M = dEdx(10e6, mProton, Z_cu, N_cu)
pb_proton_10M = dEdx(10e6, mProton, Z_pb, N_pb)

air_proton_100M = dEdx(100e6, mProton, Z_air, N_air)
al_proton_100M = dEdx(100e6, mProton, Z_al, N_al)
cu_proton_100M = dEdx(100e6, mProton, Z_cu, N_cu)
pb_proton_100M = dEdx(100e6, mProton, Z_pb, N_pb)

air_proton_1000M = dEdx(1000e6, mProton, Z_air, N_air)
al_proton_1000M = dEdx(1000e6, mProton, Z_al, N_al)
cu_proton_1000M = dEdx(1000e6, mProton, Z_cu, N_cu)
pb_proton_1000M = dEdx(1000e6, mProton, Z_pb, N_pb)

air_rho_proto_10M = dEdx_rho(10e6, mProton, Z_air)
al_rho_proton_10M = dEdx_rho(10e6, mProton, Z_al)
cu_rho_proton_10M = dEdx_rho(10e6, mProton, Z_cu)
pb_rho_proton_10M = dEdx_rho(10e6, mProton, Z_pb)

air_rho_proton_100M = dEdx_rho(100e6, mProton, Z_air)
al_rho_proton_100M = dEdx_rho(100e6, mProton, Z_al)
cu_rho_proton_100M = dEdx_rho(100e6, mProton, Z_cu)
pb_rho_proton_100M = dEdx_rho(100e6, mProton, Z_pb)

air_rho_proton_1000M = dEdx_rho(1000e6, mProton, Z_air)
al_rho_proton_1000M = dEdx_rho(1000e6, mProton, Z_al)
cu_rho_proton_1000M = dEdx_rho(1000e6, mProton, Z_cu)
pb_rho_proton_1000M = dEdx_rho(1000e6, mProton, Z_pb)

air_meson_10M = dEdx(10e6, mMeson, Z_air, N_air)
al_meson_10M = dEdx(10e6, mMeson, Z_al, N_al)
cu_meson_10M = dEdx(10e6, mMeson, Z_cu, N_cu)
pb_meson_10M = dEdx(10e6, mMeson, Z_pb, N_pb)

air_meson_100M = dEdx(100e6, mMeson, Z_air, N_air)
al_meson_100M = dEdx(100e6, mMeson, Z_al, N_al)
cu_meson_100M = dEdx(100e6, mMeson, Z_cu, N_cu)
pb_meson_100M = dEdx(100e6, mMeson, Z_pb, N_pb)

air_meson_1000M = dEdx(1000e6, mMeson, Z_air, N_air)
al_meson_1000M = dEdx(1000e6, mMeson, Z_al, N_al)
cu_meson_1000M = dEdx(1000e6, mMeson, Z_cu, N_cu)
pb_meson_1000M = dEdx(1000e6, mMeson, Z_pb, N_pb)

air_rho_meson_10M = dEdx_rho(10e6, mMeson, Z_air)
al_rho_meson_10M = dEdx_rho(10e6, mMeson, Z_al)
cu_rho_meson_10M = dEdx_rho(10e6, mMeson, Z_cu)
pb_rho_meson_10M = dEdx_rho(10e6, mMeson, Z_pb)

air_rho_meson_100M = dEdx_rho(100e6, mMeson, Z_air)
al_rho_meson_100M = dEdx_rho(100e6, mMeson, Z_al)
cu_rho_meson_100M = dEdx_rho(100e6, mMeson, Z_cu)
pb_rho_meson_100M = dEdx_rho(100e6, mMeson, Z_pb)

air_rho_meson_1000M = dEdx_rho(1000e6, mMeson, Z_air)
al_rho_meson_1000M = dEdx_rho(1000e6, mMeson, Z_al)
cu_rho_meson_1000M = dEdx_rho(1000e6, mMeson, Z_cu)
pb_rho_meson_1000M = dEdx_rho(1000e6, mMeson, Z_pb)
