export const useApi = () => {
  const token = ref<string | null>(null);

  if (import.meta.client) {
    token.value = localStorage.getItem('token');
  }

  const api = $fetch.create({
    baseURL: '/api/v1',
    onRequest({ options }) {
      const currentToken =
        token.value ||
        (import.meta.client ? localStorage.getItem('token') : null);

      if (currentToken) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${currentToken}`,
        } as any;
      }
    },
    onResponse({ response }) {
      // Unwrap backend envelope
      if (response._data && typeof response._data === 'object' && 'code' in response._data && 'data' in response._data) {
        const body = response._data as any;
        if (body.pagination) {
          // PaginatedResponse: { code, data: [...], pagination: { page, per_page, total } }
          // -> { items: [...], total: N, page: N, perPage: N }
          response._data = {
            items: body.data,
            total: body.pagination.total,
            page: body.pagination.page,
            perPage: body.pagination.per_page,
          };
        } else {
          // Response: { code, message, data } -> data
          response._data = body.data;
        }
      }
    },
    onResponseError({ response }) {
      if (response.status === 401 && import.meta.client) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        const router = useRouter();
        const route = useRoute();
        if (!route.path.startsWith('/admin/login')) {
          router.push('/admin/login');
        }
      }
    },
  });

  return api;
};
