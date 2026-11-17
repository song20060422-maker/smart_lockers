/** @ticket 07-数字人Coze：Redux Toolkit 状态层（对话列表 / 连接状态 / 当前动作） */
import { configureStore, createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { ChatMessage, Live2DAction } from '@/types'

interface HubState {
  messages: ChatMessage[]
  connected: boolean
  thinking: boolean
  action: Live2DAction | null
  error: string | null
}

const initialState: HubState = { messages: [], connected: false, thinking: false, action: null, error: null }

const hubSlice = createSlice({
  name: 'hub',
  initialState,
  reducers: {
    appendMessage(state, a: PayloadAction<ChatMessage>) { state.messages.push(a.payload) },
    setConnected(state, a: PayloadAction<boolean>) { state.connected = a.payload },
    setThinking(state, a: PayloadAction<boolean>) { state.thinking = a.payload },
    setAction(state, a: PayloadAction<Live2DAction | null>) { state.action = a.payload },
    setError(state, a: PayloadAction<string | null>) { state.error = a.payload },
  },
})

export const { appendMessage, setConnected, setThinking, setAction, setError } = hubSlice.actions
export const store = configureStore({ reducer: { hub: hubSlice.reducer } })
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
