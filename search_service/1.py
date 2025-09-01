def first_value(n, r):
    x = (1 + r)**n
    return x * r * n / (x - 1)

def second_value(n, r):
    return 1 + r * (n + 1) / 2

n = int(input())
r = float(input())
print(first_value(n, r), second_value(n, r))