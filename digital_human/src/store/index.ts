/** @ticket 07-数字人Coze：Redux Toolkit 状态层（对话列表 / 连接状态 / 当前动作） */
import { configureStore, createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { ChatMessage, Live2DAction } from '@/types'

// 定义接口，约束对象的结构：interface HubState {

interface HubState {
  // 定义对象属性 messages，取值为 ChatMessage[]

  messages: ChatMessage[]
  // 定义对象属性 connected，取值为 boolean

  connected: boolean
  // 定义对象属性 thinking，取值为 boolean

  thinking: boolean
  // 定义对象属性 action，取值为 Live2DAction | null

  action: Live2DAction | null
  // 定义对象属性 speaking，取值为 boolean（是否正在语音播报，驱动数字人说话口型）

  speaking: boolean

  // 定义对象属性 speakPulse，取值为 number（语音边界计数，每次边界递增以触发一次张口脉冲）

  speakPulse: number

  // 定义对象属性 error，取值为 string | null

  error: string | null
// 当前代码块结束

}

// 声明变量 initialState: HubState = { messages: [], connected: false, thinking: false, action: null, error: null }，用于保存数据或函数引用

const initialState: HubState = { messages: [], connected: false, thinking: false, action: null, speaking: false, speakPulse: 0, error: null }

// 声明变量 hubSlice = createSlice({，用于保存数据或函数引用

const hubSlice = createSlice({
  // 定义对象属性 name，取值为 'hub'

  name: 'hub',
  // 执行 initialState,，实现对应的前端逻辑

  initialState,
  // 定义对象属性 reducers，取值为 {

  reducers: {
    // 调用 appendMessage(state, a: PayloadAction<ChatMessage>) { state.messages.push(a.payload) },，触发对应功能

    appendMessage(state, a: PayloadAction<ChatMessage>) { state.messages.push(a.payload) },
    // 赋值更新：setConnected(state, a: PayloadAction<boolean>) { state.connected = a.payload },

    setConnected(state, a: PayloadAction<boolean>) { state.connected = a.payload },
    // 赋值更新：setThinking(state, a: PayloadAction<boolean>) { state.thinking = a.payload },

    setThinking(state, a: PayloadAction<boolean>) { state.thinking = a.payload },
    // 赋值更新：setAction(state, a: PayloadAction<Live2DAction | null>) { state.action = a.payload },

    setAction(state, a: PayloadAction<Live2DAction | null>) { state.action = a.payload },
    // 赋值更新：setError(state, a: PayloadAction<string | null>) { state.error = a.payload },

    // 赋值更新：setSpeaking(state, a: PayloadAction<boolean>) { state.speaking = a.payload }

    setSpeaking(state, a: PayloadAction<boolean>) { state.speaking = a.payload },

    // 赋值更新：bumpSpeakPulse(state) { state.speakPulse += 1 }，语音边界触发时递增计数

    bumpSpeakPulse(state) { state.speakPulse += 1 },

    setError(state, a: PayloadAction<string | null>) { state.error = a.payload },

    // 清空会话：账号切换/重连时重置对话与全部瞬时状态，避免上一个账号的欢迎语与消息残留 @ticket 07
    // 返回全新初始对象（不复用 initialState 引用）防止后续 push 污染初始态

    clearConversation() { return { messages: [], connected: false, thinking: false, action: null, speaking: false, speakPulse: 0, error: null } },
  // 定义结束

  },
// 执行 })，实现对应的前端逻辑

})

// 导出 const { appendMessage, setConnected, setThinking, setAction, setError } = hubSlice.actions，供其他模块引用

export const { appendMessage, setConnected, setThinking, setAction, setSpeaking, bumpSpeakPulse, setError, clearConversation } = hubSlice.actions
// 导出 const store = configureStore({ reducer: { hub: hubSlice.reducer } })，供其他模块引用

export const store = configureStore({ reducer: { hub: hubSlice.reducer } })
// 导出 type RootState = ReturnType<typeof store.getState>，供其他模块引用

export type RootState = ReturnType<typeof store.getState>
// 导出 type AppDispatch = typeof store.dispatch，供其他模块引用

export type AppDispatch = typeof store.dispatch
