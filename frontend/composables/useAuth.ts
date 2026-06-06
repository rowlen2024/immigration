function decodeJwt(token: string): Record<string, unknown> | null {
  try {
    const payload = token.split('.')[1];
    if (!payload) return null;
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
    const decoded = atob(base64);
    return JSON.parse(decoded);
  } catch {
    return null;
  }
}

export const useAuth = () => {
  const token = ref<string | null>(null);
  const user = ref<Record<string, unknown> | null>(null);

  if (import.meta.client) {
    const storedToken = localStorage.getItem('token');
    if (storedToken) {
      token.value = storedToken;
      user.value = decodeJwt(storedToken);
    }
  }

  const isAuthenticated = computed(() => !!token.value);

  const login = async (credentials: { username: string; password: string }) => {
    const api = useApi();
    const response = await api<{ access_token: string }>(
      '/auth/login',
      {
        method: 'POST',
        body: credentials,
      }
    );

    token.value = response.access_token;
    user.value = decodeJwt(response.access_token);

    if (import.meta.client) {
      localStorage.setItem('token', response.access_token);
    }
  };

  const logout = () => {
    token.value = null;
    user.value = null;

    if (import.meta.client) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }

    const router = useRouter();
    router.push('/');
  };

  const refresh = async () => {
    try {
      const api = useApi();
      const response = await api<{ access_token: string }>(
        '/auth/refresh',
        { method: 'POST' }
      );

      token.value = response.access_token;
      user.value = decodeJwt(response.access_token);

      if (import.meta.client) {
        localStorage.setItem('token', response.access_token);
      }
    } catch {
      logout();
    }
  };

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    refresh,
  };
};
