from run import two_sum
import json

with open("testcases.json", "r") as f:
    content = f.read()

for idx, line in enumerate(json.loads(content)):
    arr = line["input"]["nums"]
    target = line["input"]["target"]
    result = line["expected_output"]
    calculation = two_sum(arr, target)

    assert calculation == result, f"""
        Failed at testcase {idx+1}.
        Expected {result}, got {calculation}.
    """ 

print(f"""
Passed {idx+1}/{idx+1}.
""")