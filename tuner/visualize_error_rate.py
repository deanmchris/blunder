import matplotlib.pyplot as plt
import os

error_file_path = os.path.join(os.path.dirname(os.getcwd()), 'blunder', 'errors.txt')
error_rates = []

with open(error_file_path, 'r') as infile:
    for line in infile:
        error_rates.append(float(line.strip('\n')))

steps = list(range(1, len(error_rates)+1))
plt.plot(steps, error_rates)
plt.show()
