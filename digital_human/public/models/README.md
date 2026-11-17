<!-- @ticket 07-数字人Coze：Live2D 模型资源接入说明（依赖隔离，模型不随仓库提交） -->
# Live2D 模型资源

本项目基于 `awesome-digital-human-live2d`（https://github.com/wan-h/awesome-digital-human-live2d）接入 Live2D 渲染。

## 放置模型
1. 将 Live2D 模型目录（含 `*.model3.json`、纹理与动作文件）放入 `public/models/<角色名>/`。
2. 在 `.env.local` 设置 `VITE_LIVE2D_MODEL_PATH=/models/<角色名>/<角色名>.model3.json`。
3. Live2D Core SDK 为商业授权组件，请按官方许可自行获取；模型资源文件不纳入本仓库。

## 角色切换
在 `src/config/index.ts` 的 `ROLES` 中登记角色，前端通过 `setRole(roleId)` 热切换模型路径与语音音色。
