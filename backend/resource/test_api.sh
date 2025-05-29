#!/bin/bash

# API测试脚本，用于验证后端API和响应格式是否正确

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 从配置文件读取API基础URL
API_BASE_URL=$(grep -A 1 "BASE_URL" ../config.js 2>/dev/null | grep -o "'http://[^']*'" | tr -d "'")

if [ -z "$API_BASE_URL" ]; then
  # 默认使用localhost:8000
  API_BASE_URL="http://localhost:8000"
  echo "未找到配置文件，使用默认API地址: $API_BASE_URL"
else
  echo "使用配置的API地址: $API_BASE_URL"
fi

# 测试函数
test_api() {
  local endpoint=$1
  local method=$2
  local data=$3
  local auth_header=$4
  local description=$5

  echo -e "\n测试: $description ($method $endpoint)"
  
  # 构建curl命令
  cmd="curl -s -X $method"
  
  # 添加认证头（如果有）
  if [ ! -z "$auth_header" ]; then
    cmd="$cmd -H \"Authorization: Bearer $auth_header\""
  fi
  
  # 添加数据（如果有）
  if [ ! -z "$data" ]; then
    cmd="$cmd -H \"Content-Type: application/json\" -d '$data'"
  fi
  
  # 完成命令
  cmd="$cmd \"$API_BASE_URL$endpoint\""
  
  # 输出要执行的命令
  echo "执行: $cmd"
  
  # 执行命令并解析响应
  response=$(eval $cmd)
  
  # 打印响应并检查格式
  echo "响应: $response"
  
  # 检查响应格式是否正确
  if echo $response | grep -q '"code":0'; then
    echo -e "${GREEN}✓ 成功: 响应格式正确${NC}"
    return 0
  else
    echo -e "${RED}✗ 失败: 响应格式不正确${NC}"
    return 1
  fi
}

# 记录测试结果
total=0
passed=0

# 执行测试

# 1. 测试分类列表接口
echo -e "\n=== 测试分类列表接口 ==="
total=$((total+1))
test_api "/api/category/list" "GET" "" "your_test_token" "获取分类列表"
[ $? -eq 0 ] && passed=$((passed+1))

# 2. 测试商品列表接口
echo -e "\n=== 测试商品列表接口 ==="
total=$((total+1))
test_api "/api/product/list" "GET" "" "your_test_token" "获取商品列表"
[ $? -eq 0 ] && passed=$((passed+1))

# 3. 测试按分类过滤商品
echo -e "\n=== 测试按分类过滤商品 ==="
total=$((total+1))
test_api "/api/product/list?categoryId=1" "GET" "" "your_test_token" "按分类获取商品列表"
[ $? -eq 0 ] && passed=$((passed+1))

# 4. 测试商品详情接口
echo -e "\n=== 测试商品详情接口 ==="
total=$((total+1))
test_api "/api/product/detail?id=1" "GET" "" "your_test_token" "获取商品详情"
[ $? -eq 0 ] && passed=$((passed+1))

# 显示测试结果汇总
echo -e "\n=== 测试结果汇总 ==="
echo -e "总测试数: $total"
echo -e "通过测试: $passed"
echo -e "失败测试: $((total-passed))"

if [ $passed -eq $total ]; then
  echo -e "${GREEN}所有测试通过!${NC}"
else
  echo -e "${RED}部分测试失败，请检查API实现和响应格式。${NC}"
fi 