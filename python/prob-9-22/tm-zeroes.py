from scipy.optimize import newton
from scipy.special import spherical_jn

def find_tm_root(l, x0, tol=1e-6, max_iter=100):
    def f(x):
        return spherical_jn(l, x) + x * spherical_jn(l, x, derivative=True)

    def df(x):
        return -x * spherical_jn(l + 1, x, derivative=True) - spherical_jn(l+1, x) +(l+1)*spherical_jn(l, x, derivative=True)

    # Use scipy's newton method with both function and its derivative
    root = newton(f, x0, fprime=df, tol=tol, maxiter=max_iter)
    return root

if __name__ == "__main__":
    # These initial guesses are eyeballed from the plots plot-tm.py
    initial_guesses = [
            [2, 5, 8, 11, 14],
            [3, 6, 9.5, 12, 15.5],
            [4, 7, 11, 14, 17],
            [5.5, 8.5, 12, 15, 19],
            [6, 10, 13.5, 16.5, 20]
            ]
    # Format to LaTeX table
    for i in range(5):
        line = f"{i+1} &"
        for l in range(5):
            root = find_tm_root(l, initial_guesses[l][i])
            line += f" {root:.5f} &"
        print(line[:-1] + r"\str\\")

