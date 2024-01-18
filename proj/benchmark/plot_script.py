import subprocess
import statistics
import matplotlib
matplotlib.use('Agg')  
import matplotlib.pyplot as plt

# List of test cases
test_cases = ["xsmall", "small", "medium", "large", "xlarge"]

# List of thread counts to test
thread_counts = [2, 4, 6, 8, 12]

# Dictionary to store the results
speedup_data = {test_case: [] for test_case in test_cases}
print("starting")

for test_case in test_cases:
    elapsed_times_sequential = []
    for _ in range(5):
        cmd = ["go", "run", "benchmark.go", "s", test_case]
        output = subprocess.check_output(cmd, text=True)
        elapsed_time = float(output.strip())
        elapsed_times_sequential.append(elapsed_time)
    
    avg_elapsed_time_sequential = statistics.mean(elapsed_times_sequential)

    # benchmark for parallel version and calculate the speedup
    for thread_count in thread_counts:
        elapsed_times_parallel = []
        for _ in range(5):
            cmd = ["go", "run", "benchmark.go", "p", test_case, str(thread_count)]
            output = subprocess.check_output(cmd, text=True)
            elapsed_time = float(output.strip())
            elapsed_times_parallel.append(elapsed_time)

       
        avg_elapsed_time_parallel = statistics.mean(elapsed_times_parallel)

        speedup = avg_elapsed_time_sequential / avg_elapsed_time_parallel
        speedup_data[test_case].append(speedup)
    print(speedup_data)


for test_case in test_cases:
    plt.plot(thread_counts, speedup_data[test_case], marker='o', label=test_case)

plt.xlabel("Number of Threads")
plt.ylabel("Speedup")
plt.title("Speedup vs. Number of Threads")
plt.legend()
plt.grid()
plt.show()

output_file = "speedup_graph_local.png"
plt.savefig(output_file)

print(f"Speedup graph saved to {output_file}")
print("End")