<template>
  <el-form-item :label="label" :prop="prop" :rules="rules">
    <!-- 输入框 -->
    <el-input 
      v-if="type === 'input'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :maxlength="maxlength"
      :show-word-limit="showWordLimit"
      :clearable="clearable"
    />
    
    <!-- 文本域 -->
    <el-input 
      v-else-if="type === 'textarea'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :maxlength="maxlength"
      :show-word-limit="showWordLimit"
      :rows="rows"
      type="textarea"
    />
    
    <!-- 数字输入框 -->
    <el-input-number 
      v-else-if="type === 'number'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :min="min"
      :max="max"
      :step="step"
      :precision="precision"
      :controls="controls"
    />
    
    <!-- 选择器 -->
    <el-select 
      v-else-if="type === 'select'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :clearable="clearable"
      :filterable="filterable"
      :multiple="multiple"
      :collapse-tags="collapseTags"
    >
      <el-option 
        v-for="item in options" 
        :key="item.value" 
        :label="item.label" 
        :value="item.value"
        :disabled="item.disabled"
      />
    </el-select>
    
    <!-- 日期选择器 -->
    <el-date-picker 
      v-else-if="type === 'date'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :type="dateType || 'date'"
      :format="format"
      :value-format="valueFormat"
      :clearable="clearable"
    />
    
    <!-- 时间选择器 -->
    <el-time-picker 
      v-else-if="type === 'time'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :placeholder="placeholder" 
      :disabled="disabled"
      :format="format"
      :value-format="valueFormat"
      :clearable="clearable"
    />
    
    <!-- 开关 -->
    <el-switch 
      v-else-if="type === 'switch'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :disabled="disabled"
      :active-text="activeText"
      :inactive-text="inactiveText"
      :active-value="activeValue"
      :inactive-value="inactiveValue"
    />
    
    <!-- 单选框组 -->
    <el-radio-group 
      v-else-if="type === 'radio'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :disabled="disabled"
    >
      <el-radio 
        v-for="item in options" 
        :key="item.value" 
        :label="item.value"
        :disabled="item.disabled"
      >
        {{ item.label }}
      </el-radio>
    </el-radio-group>
    
    <!-- 复选框组 -->
    <el-checkbox-group 
      v-else-if="type === 'checkbox'" 
      :model-value="modelValue" 
      @update:model-value="updateValue"
      :disabled="disabled"
    >
      <el-checkbox 
        v-for="item in options" 
        :key="item.value" 
        :label="item.value"
        :disabled="item.disabled"
      >
        {{ item.label }}
      </el-checkbox>
    </el-checkbox-group>
    
    <!-- 默认插槽 -->
    <slot></slot>
    
    <!-- 帮助文本 -->
    <div v-if="help" class="form-item-help">
      <el-icon><InfoFilled /></el-icon>
      <span>{{ help }}</span>
    </div>
  </el-form-item>
</template>

<script setup>
import { computed } from 'vue';
import { InfoFilled } from '@element-plus/icons-vue';

const props = defineProps({
  // 表单项类型
  type: {
    type: String,
    default: 'input',
    validator: (value) => [
      'input', 'textarea', 'number', 'select', 'date', 
      'time', 'switch', 'radio', 'checkbox'
    ].includes(value)
  },
  // 标签
  label: {
    type: String,
    default: ''
  },
  // 属性名
  prop: {
    type: String,
    default: ''
  },
  // 验证规则
  rules: {
    type: [Object, Array],
    default: () => []
  },
  // 占位符
  placeholder: {
    type: String,
    default: '请输入'
  },
  // 是否禁用
  disabled: {
    type: Boolean,
    default: false
  },
  // 最大长度
  maxlength: {
    type: Number,
    default: undefined
  },
  // 是否显示字数统计
  showWordLimit: {
    type: Boolean,
    default: false
  },
  // 是否可清空
  clearable: {
    type: Boolean,
    default: true
  },
  // 文本域行数
  rows: {
    type: Number,
    default: 3
  },
  // 数字输入框最小值
  min: {
    type: Number,
    default: -Infinity
  },
  // 数字输入框最大值
  max: {
    type: Number,
    default: Infinity
  },
  // 数字输入框步长
  step: {
    type: Number,
    default: 1
  },
  // 数字输入框精度
  precision: {
    type: Number,
    default: undefined
  },
  // 数字输入框是否显示控制按钮
  controls: {
    type: Boolean,
    default: true
  },
  // 选择器选项
  options: {
    type: Array,
    default: () => []
  },
  // 选择器是否可过滤
  filterable: {
    type: Boolean,
    default: false
  },
  // 选择器是否多选
  multiple: {
    type: Boolean,
    default: false
  },
  // 选择器多选时是否折叠标签
  collapseTags: {
    type: Boolean,
    default: false
  },
  // 日期选择器类型
  dateType: {
    type: String,
    default: 'date'
  },
  // 日期/时间格式
  format: {
    type: String,
    default: ''
  },
  // 日期/时间值格式
  valueFormat: {
    type: String,
    default: ''
  },
  // 开关激活文本
  activeText: {
    type: String,
    default: ''
  },
  // 开关未激活文本
  inactiveText: {
    type: String,
    default: ''
  },
  // 开关激活值
  activeValue: {
    type: [String, Number, Boolean],
    default: true
  },
  // 开关未激活值
  inactiveValue: {
    type: [String, Number, Boolean],
    default: false
  },
  // 帮助文本
  help: {
    type: String,
    default: ''
  },
  // 绑定值
  modelValue: {
    type: [String, Number, Boolean, Array, Object, Date],
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'change']);

// 更新值的方法
const updateValue = (value) => {
  emit('update:modelValue', value);
  emit('change', value);
};
</script>

<style scoped>
.form-item-help {
  margin-top: 5px;
  color: #909399;
  font-size: 12px;
  display: flex;
  align-items: center;
}

.form-item-help .el-icon {
  margin-right: 5px;
  font-size: 14px;
}
</style> 