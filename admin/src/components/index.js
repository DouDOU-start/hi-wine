// 布局组件
import Sidebar from './layout/Sidebar.vue';
import Navbar from './layout/Navbar.vue';

// 通用组件
import PageHeader from './common/PageHeader.vue';
import LoadingState from './common/LoadingState.vue';
import EmptyState from './common/EmptyState.vue';
import SearchForm from './common/SearchForm.vue';
import TableToolbar from './common/TableToolbar.vue';
import ThemeSwitch from './common/ThemeSwitch.vue';

// 表单组件
import FormItem from './form/FormItem.vue';
import UploadImage from './form/UploadImage.vue';

// 导出所有组件
export {
  // 布局组件
  Sidebar,
  Navbar,
  
  // 通用组件
  PageHeader,
  LoadingState,
  EmptyState,
  SearchForm,
  TableToolbar,
  ThemeSwitch,
  
  // 表单组件
  FormItem,
  UploadImage
};

// 注册全局组件
export function registerGlobalComponents(app) {
  app.component('PageHeader', PageHeader);
  app.component('LoadingState', LoadingState);
  app.component('EmptyState', EmptyState);
  app.component('SearchForm', SearchForm);
  app.component('TableToolbar', TableToolbar);
  app.component('FormItem', FormItem);
  app.component('UploadImage', UploadImage);
  app.component('ThemeSwitch', ThemeSwitch);
} 