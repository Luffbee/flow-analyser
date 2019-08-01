#!/usr/bin/env python3

import sys
import json
import subprocess
import pandas as pd
import matplotlib.pyplot as plt


def raw2df(datas):
    sers = {}

    for data in datas:
        ser = pd.Series({int(v[0]): float(v[1]) for v in data['values']})
        name = data['metric']['store']
        if sers.get(name) is None:
            sers[name] = ser
        else:
            ser0 = sers[name]
            assert(len(ser0.index.intersection(ser.index)) == 0)
            sers[name] = ser0.add(ser, fill_value=0)

    nonempty = []
    summ = pd.Series()
    for name in sers:
        ser = sers[name]
        if ser.sum() > 10:
            ser.name = name
            summ = summ.add(ser, fill_value=0)
            nonempty.append(ser)

    summ.name = "sum"
    nonempty.append(summ)

    return pd.concat(nonempty, axis=1)


def drop_dup(ser, intvl):
    res = pd.Series()

    lastT, lastV = 0, 0
    for t, v in ser.items():
        if t >= lastT + intvl * 0.8 or v != lastV:
            res.at[t] = v
            lastT, lastV = t, v

    return res


def process(ser):
    proc = subprocess.Popen(
        ['go', 'run', 'main.go', 'stdin', 'ema'],
        stdout=subprocess.PIPE,
        stdin=subprocess.PIPE)

    proc.stdin.write(f'{len(ser)}\n'.encode())
    for t, v in ser.items():
        proc.stdin.write(f'{t} {v}\n'.encode())
    proc.stdin.close()

    out = iter(proc.stdout.read().decode().strip().split())
    data = dict()
    for t in out:
        data[int(t)] = float(next(out))

    return pd.Series(data)


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("argument not enough")
        sys.exit(1)

    data = None
    with open(sys.argv[1]) as f:
        data = raw2df(json.load(f)['data']['result'])
    data.to_html('data.html')

    for col in data.columns:
        ser = data[col].dropna()
        plt.plot(ser.index.values, ser.values, 'o-', label=col)
        #ser = drop_dup(ser, 60)
        after = process(ser)
        plt.plot(after.index.values, after.values, 'x-', label=col+'after')
        plt.legend()
        plt.show()
