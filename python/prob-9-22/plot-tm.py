import numpy as np
import matplotlib.pyplot as plt
from scipy.special import spherical_jn

def plot_tm(l_max, x_max):
    x = np.linspace(0, x_max, 1000)  # Define x range, avoid zero to prevent division by zero

    plt.figure(figsize=(10, 6))

    for l in range(l_max + 1):
        j_l = spherical_jn(l, x)
        j_l_prime = spherical_jn(l, x, derivative=True)
        plt.plot(x, j_l + x * j_l_prime, label=f'$j_{l}(x)+xj_{l}\'(x)$')

    # Configure plot
    plt.xlabel('x')
    plt.ylabel(r'$j_l(x)+xj_l\'(x)$')
    plt.title('Spherical Bessel Functions $j_l(x)+xj_l\'(x)$')
    plt.legend()
    plt.grid(True)
    plt.minorticks_on()
    plt.show()

# Example usage:
if __name__ == "__main__":
    plot_tm(4, 25)
