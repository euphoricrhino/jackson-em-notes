from scipy.optimize import newton
from scipy.special import spherical_jn

def find_te_root(l, x0, tol=1e-6, max_iter=100):
    def f(x):
        return spherical_jn(l, x)

    def df(x):
        return spherical_jn(l, x, derivative=True)

    # Use scipy's newton method with both function and its derivative
    root = newton(f, x0, fprime=df, tol=tol, maxiter=max_iter)
    return root

if __name__ == "__main__":
    # These initial guesses are eyeballed from the plots plot-te.py
    initial_guesses = [
            [3, 6, 9, 12, 16],
            [4, 8, 11, 14, 17],
            [6, 9, 12.5, 15.5, 19],
            [7.5, 10.5, 13.5, 17, 21],
            [8.5, 11.5, 15, 18, 22]
            ]
    # Format to LaTeX table
    for i in range(5):
        line = f"{i+1} &"
        for l in range(5):
            root = find_te_root(l, initial_guesses[l][i])
            line += f" {root:.5f} &"
        print(line[:-1] + r"\str\\")

