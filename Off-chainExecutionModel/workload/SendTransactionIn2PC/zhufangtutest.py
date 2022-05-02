import matplotlib.pyplot as plt
import matplotlib

# matplotlib.use('Agg')
list_00 = [670.2188264, 802.5360138, 1034.846851, 1193.409767, 1247.4673, 1327.77466, 1452.383796, 1476.887575]
list_01 = [648.4952883, 806.8928824, 1012.563701, 1235.244617, 1283.501392, 1352.670534, 1389.62132, 1397.284574]
list_02 = [276, 302, 313, 315, 314, 315, 315, 315]
list_03 = [673.0522499, 759.2690138, 954.8732853, 1133.391302, 1151.836324, 1153.336284, 1303.563375, 1302.219114]

# list_00 = [0.16, 0.33, 0.20, 0.31]
# list_01 = [0.05, 0.33, 0.39, 0.29]
# list_02 = [0.07, 0.5, 0.26, 0.38]
# list_03 = [0.04, 0.13, 0.19, 0.15]


name_list = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']
x = list(range(len(name_list)))
total_width, n = 0.4, 4
width = total_width / n
plt.bar(x, list_00, width=width, label='a', tick_label=name_list, fc='y')
for i in range(len(x)):
    x[i] = x[i] + width
plt.bar(x, list_01, width=width, label='b', fc='r')
for i in range(len(x)):
    x[i] = x[i] + width
plt.bar(x, list_02, width=width, label='c', fc='b')
for i in range(len(x)):
    x[i] = x[i] + width
plt.bar(x, list_03, width=width, label='d', fc='g')
plt.legend()
plt.show()
