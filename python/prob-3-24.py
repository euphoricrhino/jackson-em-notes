import scipy.special
import math

def calc1(x):
    v = scipy.special.jv(1, x / 2)
    w = scipy.special.jv(1, x)
    return v / (w * w * x) * math.sinh(x) / math.sinh(2 * x)

j0zeros = scipy.special.jn_zeros(0, 10)
print("formula (9)")
sum1 = 0
for n, x in enumerate(j0zeros):
    sum1 += calc1(x)
    print("  iteration {}: {:.10f}".format(n, sum1))

def calc2(n):
    t1 = (2 * n + 1) * math.pi / 4
    t2 = (2 * n + 1) * math.pi / 2
    return 0.5 * (1 if n % 2 == 1 else -1) * (
            scipy.special.kn(1, t1) +
            scipy.special.kn(0, t2) * scipy.special.iv(1, t1) / scipy.special.iv(0, t2))

sum2 = 0.5 
print("formula (17)")
for n in range(10):
    sum2 += calc2(n)
    print("  iteration {}: {:.10f}".format(n, sum2))

