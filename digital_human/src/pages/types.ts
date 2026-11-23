/** @ticket 07-数字人Coze：数字人场景角色类型定义（登录页与主界面共用） */

// 数字人场景角色：寄件引导 sender / 取件指引 receiver / 智能客服 service
// 导出 type DhRole，供其他模块引用

export type DhRole = 'sender' | 'receiver' | 'service'
