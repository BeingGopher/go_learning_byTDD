# 复习（DP）

## 遇到问题如何解决

1. 这道题目我举例推导状态转移公式了么？
2. 我打印dp数组的日志了么？
3. 打印出来了dp数组和我想的一样么？
4. 遍历顺序确定：可以看dp[i]依赖于什么（比如依赖于dp[i-1]，那说明就是从前往后遍历）

## 背包

01背包优先遍历物品，背包遍历倒序（防止重复使用）

完全背包，遍历背包顺序（可以重复加入），如果是组合，就是先物品后背包；如果是排列，就是先背包后物品（排列讲究顺序）

## 递推公式

问能否能装满背包（或者最多装多少）：dp[j] = max(dp[j], dp[j - nums[i]] + nums[i]);

问装满背包有几种方法：dp[j] += dp[j - nums[i]] 

问背包装满最大价值：dp[j] = max(dp[j], dp[j - weight[i]] + value[i]); 

问装满背包所有物品的最小个数：dp[j] = min(dp[j - coins[i]] + 1, dp[j]); 

> [!IMPORTANT]
>
> 之所以都分为两个部分，是因为选择的时候：
>
> 1.选当前物品
>
> 2.不选当前物品，这时候就是去掉当前物品的容量再加上当前物品的价值
>
> 最后，一定要注意dp的定义、种类区分、遍历顺序

## 周总结部分

### 树形dp

使用后序遍历+递归，因为后序遍历会返回结果值到上一层（递归过程中），同时对于dp的定义也会有要求，通常不是节点数量（比如是长度为2的数组，0表示不偷该节点，1表示偷该节点）

### 买卖股票

dp(i)(1)，**表示的是第i天，买入股票的状态，并不是说一定要第i天买入股票

```go
for j := 0; j < 2 * k - 1; j += 2 {
    dp[i][j + 1] = max(dp[i - 1][j + 1], dp[i - 1][j] - prices[i]);//买入或者不买入
    dp[i][j + 2] = max(dp[i - 1][j + 2], dp[i - 1][j + 1] + prices[i]);//卖出或者不卖出
}
```

> 股票问题很可能涉及状态转移，要学会定义状态

### 编辑距离

dp(i)(j)表示s[i]与t[j]之间的最近编辑距离

dp(i-1)(j-1)表示不需要操作，就是相等的，使用之前的数量即可

dp(i-1)(j)+1，表示s需要删除一个元素（t不变）

dp(i)(j-1)+1，表示t需要删除一个元素（s不变）

删除元素和增加元素之间的操作数量是相同的

dp(i-1)(j-1)+1，表示需要替换元素，（s和t都不变）
