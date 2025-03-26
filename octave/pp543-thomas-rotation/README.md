# Misc functions for relativistic calculation for the purpose of numerically verify Thomas rotation.

## Verification of Thomas Rotation

* Generate random velocities **u**,**v**

```
octave:1> u=random_3velocity
u =

  -0.410839
   0.039569
  -0.416401

octave:2> v=random_3velocity
v =

   0.4523
  -0.2336
  -0.7617
```

* Their corresponding boost matrix
```
octave:3> Bu=boost_matrix(u)
Bu =

   1.234421   0.507148  -0.048844   0.514015
   0.507148   1.115108  -0.011086   0.116666
  -0.048844  -0.011086   1.001068  -0.011236
   0.514015   0.116666  -0.011236   1.118246

octave:4> Bv=boost_matrix(v)
Bv =

   2.4950  -1.1286   0.5827   1.9005
  -1.1286   1.3644  -0.1882  -0.6137
   0.5827  -0.1882   1.0972   0.3169
   1.9005  -0.6137   0.3169   2.0334
```

* Composite velocities **u**⊕**v** and **v**⊕**u** and verify they have the same norm
```
octave:5> u_oplus_v=add_velocity(u,v)
u_oplus_v =

  -0.064271
  -0.130976
  -0.946039

octave:6> v_oplus_u=add_velocity(v,u)
v_oplus_u =

   0.2915
  -0.2122
  -0.8867

octave:7> norm(u_oplus_v)-norm(v_oplus_u)
ans = 0
```

* boost matrices of the composite velocities
```
octave:8> Bu_oplus_v=boost_matrix(u_oplus_v)
Bu_oplus_v =

   3.455991   0.222121   0.452653   3.269502
   0.222121   1.011072   0.022564   0.162977
   0.452653   0.022564   1.045982   0.332126
   3.269502   0.162977   0.332126   3.398937

octave:9> Bv_oplus_u=boost_matrix(v_oplus_u)
Bv_oplus_u =

   3.4560  -1.0074   0.7332   3.0645
  -1.0074   1.2278  -0.1658  -0.6928
   0.7332  -0.1658   1.1206   0.5043
   3.0645  -0.6928   0.5043   3.1076
```

* Calculate 4x4 rotation matrix via B^-1(**v**⊕**u**)B(**v**)B(**u**)
```
octave:10> R1=inverse(Bv_oplus_u)*Bv*Bu
R1 =

   1.0000e+00   7.9682e-16  -3.4614e-16  -2.3914e-15
   1.2269e-15   9.2984e-01  -3.9142e-02  -3.6588e-01
  -5.4802e-16   7.3959e-02   9.9391e-01   8.1628e-02
  -3.6247e-15   3.6046e-01  -1.0296e-01   9.2707e-01
```

* Verify it rotates **u**⊕**v** into **v**⊕**u**
```
octave:11> R1*four_velocity(u_oplus_v)
ans =

   3.4560
   1.0074
  -0.7332
  -3.0645

octave:12> four_velocity(v_oplus_u)
ans =

   3.4560
   1.0074
  -0.7332
  -3.0645
```

* Calculate 4x4 rotation matrix via B(**v**)B(**u**)B^-1(**u**⊕**v**) and verify its rotation
```
octave:13> R2=Bv*Bu*inverse(Bu_oplus_v)
R2 =

   1.0000e+00  -1.5351e-16   2.3722e-16  -4.9702e-16
  -1.2135e-16   9.2984e-01  -3.9142e-02  -3.6588e-01
   3.9334e-18   7.3959e-02   9.9391e-01   8.1628e-02
  -1.3672e-15   3.6046e-01  -1.0296e-01   9.2707e-01

octave:14> R2*four_velocity(u_oplus_v)
ans =

   3.4560
   1.0074
  -0.7332
  -3.0645

octave:15> four_velocity(v_oplus_u)
ans =

   3.4560
   1.0074
  -0.7332
  -3.0645
```

* Verify the 3x3 block is a rotation in SO(3)
```
octave:16> R1_3x3=R1(2:4,2:4)
R1_3x3 =

   0.929837  -0.039142  -0.365884
   0.073959   0.993915   0.081628
   0.360463  -0.102961   0.927074

octave:17> R1_3x3*R1_3x3'
ans =

   1.0000e+00  -1.0817e-16   2.7519e-17
  -1.0817e-16   1.0000e+00  -1.7501e-16
   2.7519e-17  -1.7501e-16   1.0000e+00

octave:18> det(R1_3x3)
ans = 1.0000
```

* Verify the rotation angle
```
octave:19> gamma_u=(1-dot(u,u))^(-.5)
gamma_u = 1.2344
octave:20> gamma_v=(1-dot(v,v))^(-.5)
gamma_v = 2.4950
octave:21> gamma=gamma_u*gamma_v*(1+dot(u,v))
gamma = 3.4560
octave:22> cos_theta=(1+gamma+gamma_u+gamma_v)^2/(gamma+1)/(gamma_u+1)/(gamma_v+1)-1
cos_theta = 0.9254
octave:23> dot(u_oplus_v,v_oplus_u)/dot(u_oplus_v,u_oplus_v)
ans = 0.9254
```

* Verify the 3x3 rotation is identical as the one with expected axis and angle
```
octave:24> axis_angle_rotation(cross(u,v),acos(cos_theta))
ans =

   0.929837  -0.039142  -0.365884
   0.073959   0.993915   0.081628
   0.360463  -0.102961   0.927074

octave:25> R1_3x3
R1_3x3 =

   0.929837  -0.039142  -0.365884
   0.073959   0.993915   0.081628
   0.360463  -0.102961   0.927074
```
