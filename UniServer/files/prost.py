import random
import sys
#быстрая сортировка Хоара
str = sys.argv[1]
arr = []
for el in str.split(','):
    arr.append(int(el))


def quicksort(nums):
    if len(nums) <= 1:
        return nums
    else:
        q = random.choice(nums)
    l_nums = [n for n in nums if n < q]

    e_nums = [q] * nums.count(q)
    b_nums = [n for n in nums if n > q]
    return quicksort(l_nums) + e_nums + quicksort(b_nums)


print(quicksort(arr))