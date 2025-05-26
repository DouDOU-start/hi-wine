const express = require('express');
const path = require('path');
const app = express();
const port = 9002;

app.use(express.static(path.join(__dirname, 'resources')));

app.listen(port, () => {
  console.log(`静态资源服务已启动: http://localhost:${port}`);
  console.log(`可通过 http://localhost:${port}/wines.json 访问配置文件`);
}); 