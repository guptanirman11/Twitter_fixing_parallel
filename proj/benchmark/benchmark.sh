#!/bin/bash
#
#SBATCH --mail-user=guptanirman11@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/guptanirman11/Parallel/project-2-guptanirman11/proj2/benchmark
#SBATCH --partition=general 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=4:00:00


module load golang/1.19
# Your command here

python3 plot_script.py