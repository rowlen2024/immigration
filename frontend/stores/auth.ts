import { defineStore } from 'pinia';

interface User {
  id: string;
  username: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,

    currentUser: (state) => state.user,
  },

  actions: {
    setAuth(token: string, user: User) {
      this.token = token;
      this.user = user;

      if (import.meta.client) {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
      }
    },

    async login(credentials: { username: string; password: string }) {
      this.loading = true;

      try {
        const api = useApi();
        const response = await api<{ token: string; user: User }>(
          '/auth/login',
          {
            method: 'POST',
            body: credentials,
          }
        );

        this.setAuth(response.token, response.user);
      } finally {
        this.loading = false;
      }
    },

    logout() {
      this.token = null;
      this.user = null;

      if (import.meta.client) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }

      const router = useRouter();
      router.push('/');
    },

    loadFromStorage() {
      if (import.meta.client) {
        const storedToken = localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');

        if (storedToken) {
          this.token = storedToken;
        }

        if (storedUser) {
          try {
            this.user = JSON.parse(storedUser);
          } catch {
            localStorage.removeItem('user');
          }
        }
      }
    },
  },
});
