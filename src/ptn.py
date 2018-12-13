#!/usr/bin/python3

import copy
import itertools
import os
import json

def _iterate_cards_(card_arr):
    mark = {}
    for flag in range(0, 2):
        ke_arr = []
        shun_arr = []
        index = 0
        is_hu = True
        dup_arr = copy.deepcopy(card_arr)
        for num_arr in dup_arr:
            for j in range(0, len(num_arr)):
                if flag == 0:
                    if num_arr[j] >= 3:
                        num_arr[j] -= 3
                        ke_arr.append(index)                    
                    while (len(num_arr) - j >= 3 and
                    num_arr[j] > 0 and 
                    num_arr[j+1] > 0 and
                    num_arr[j+2] > 0):
                        num_arr[j] -= 1
                        num_arr[j+1] -= 1
                        num_arr[j+2] -= 1 
                        shun_arr.append(index)
                else:
                    while (len(num_arr) - j >= 3 and
                    num_arr[j] > 0 and 
                    num_arr[j+1] > 0 and
                    num_arr[j+2] > 0):
                        num_arr[j] -= 1
                        num_arr[j+1] -= 1
                        num_arr[j+2] -= 1 
                        shun_arr.append(index)
                    if num_arr[j] >= 3:
                        num_arr[j] -= 3
                        ke_arr.append(index)
                index += 1
                if num_arr[j] > 0:
                    is_hu = False
                    break
            if not is_hu:
                break
        if is_hu:
            ret = [len(ke_arr)] + ke_arr + [len(shun_arr)] + shun_arr
            # print("_iterate_cards_", ke_arr, shun_arr, ret, card_arr)
            key = str(ret)
            if not mark.get(key):
                mark[key] = True
                yield ret
            

# card_arr : [[1,1,1], [1,1,1], [1,1,1], [1,1,1], [2]]
def find_comb(card_arr):
    # print("find_comb", card_arr)
    ret_arr = []
    eye_arr = []
    index = 0
    for i, nums in enumerate(card_arr):
        for j, card_num in enumerate(nums):
            if card_num >= 2:  # 找到将
                eye_arr.append([index, i, j])
            index += 1
    for eye in eye_arr:
        dup_arr = copy.deepcopy(card_arr)
        dup_arr[eye[1]][eye[2]] -= 2
        for v in _iterate_cards_(dup_arr):
            # print("find_comb_ret", eye, v)
            t = []
            t += [eye[0]]
            t += v
            ret_arr.append(t)

    return ret_arr

# card_arr : [[1,1,1], [1,1,1], [1,1,1], [1,1,1], [2]]
def ptn(card_arr):
    if len(card_arr) == 1:
        return [card_arr]
    h1 = {}
    ret_arr = []
    ret_arr += itertools.permutations(card_arr)
    for i in range(0, len(card_arr)):
        for j in range(i+1, len(card_arr)):
            arr1, arr2 = card_arr[i], card_arr[j]
            key = str([arr1, 0, arr2])
            if not h1.get(key):
                h1[key] = True
                h2 = {}
                for k in range(0, len(arr1) + len(arr2) + 1):
                    t = [0] * len(arr2) + arr1 + [0] * len(arr2)
                    for m in range(0, len(arr2)):
                        t[k+m] += arr2[m]
                    t = list(filter(lambda x: x > 0, t))
                    if len(t) > 9 or any(x > 4 for x in t):
                        continue
                    key = str(t)
                    if not h2.get(key):
                        h2[key] = True
                        t2 = copy.deepcopy(card_arr)
                        del(t2[i])
                        del(t2[j-1])
                        ret_arr += ptn([t]+t2)
    return ret_arr

# card_arr : [[1,1,1], [1,1,1], [1,1,1], [1,1,1], [2]]
def calc_key(card_arr):
    key = 0
    for arr in card_arr:
        isContinue = False
        for num in arr:
            if num == 1:
                if isContinue:
                    key = key << 1
                else:
                    key = (key << 2) + (1 << 2) - 2
            else:
                bn = 2 * num
                if isContinue:
                    bn = bn - 1
                key = (key << bn) + (1 << bn) - 2
            isContinue = True
    return key


combs2 = [
    [[2]],
]
combs5 = [
    [[1, 1, 1], [2]],
    [[3], [2]],
]
combs8 = [
    [[1, 1, 1], [1, 1, 1], [2]],
    [[1, 1, 1], [3], [2]],
    [[3], [3], [2]],
]
combs11 = [
    [[1, 1, 1], [1, 1, 1], [1, 1, 1], [2]],
    [[1, 1, 1], [1, 1, 1], [3], [2]],
    [[1, 1, 1], [3], [3], [2]],
    [[3], [3], [3], [2]],
]
combs14 = [
    [[1, 1, 1], [1, 1, 1], [1, 1, 1], [1, 1, 1], [2]],
    [[1, 1, 1], [1, 1, 1], [1, 1, 1], [3], [2]],
    [[1, 1, 1], [1, 1, 1], [3], [3], [2]],
    [[1, 1, 1], [3], [3], [3], [2]],
    [[3], [3], [3], [3], [2]],
]


def record(*args):
    max = []
    keyMap = {}
    for combs in args:
        arr = []
        for v in combs:
            arr += ptn(v)
        for v in arr:
            k = calc_key(v)
            if not keyMap.get(k):
                keyMap[k] = find_comb(v)
                if len(keyMap[k]) > len(max):
                    max = keyMap[k]

    dir_path = os.path.dirname(os.path.realpath(__file__))
    output_path = os.path.join(dir_path, "output.json")
    with open(output_path, "w") as f:
        json.dump(keyMap, f)
        f.close()
    print("write finished, total", len(keyMap), max)

# record(combs14, combs11, combs8, combs5, combs2)

# arr = find_comb([[4,4,4,2]])
# print(arr)

# key = calc_key([[3,1,1,1],[1,1,1],[3],[2]])
# print(format(key, "b"))