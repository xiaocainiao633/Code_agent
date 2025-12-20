declare module '@/stores/auth' {
  export interface UserInfo {
    username: string
    token: string
    role: 'admin' | 'user' | 'guest' | 'local'
    rememberMe?: boolean
  }

  export const useAuthStore: () => {
    userInfo: import('vue').Ref<UserInfo | null>
    isAuthenticated: import('vue').ComputedRef<boolean>
    isGuest: import('vue').ComputedRef<boolean>
    isLocal: import('vue').ComputedRef<boolean>
    isAdmin: import('vue').ComputedRef<boolean>
    login: (info: UserInfo) => void
    logout: () => void
    loadUserInfo: () => void
    updateUserInfo: (updates: Partial<UserInfo>) => void
    clearStorage: () => void
  }
}