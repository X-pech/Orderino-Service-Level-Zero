vegeta attack -name=stress -duration=10s -rate=4000/s -targets names.txt -output results.bin && cat results.bin | vegeta plot > plot.html
