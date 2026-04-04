import axios from "axios";

// 自定义 API 错误类，包含错误码
export class ApiError extends Error {
  constructor(code, message) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      // 不弹错误提示，由组件自行处理
      return Promise.reject(new ApiError(res.code, res.message));
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/login";
      return Promise.reject(error);
    }

    // 不在这里自动显示错误消息，让各组件自己处理
    // 这样可以避免重复弹窗，并且组件可以提供更有针对性的错误信息
    return Promise.reject(error);
  },
);

export const pollApi = {
  createPoll: (data) => api.post('/polls', data),
  getPoll: (id) => api.get(`/polls/${id}`),
  getPollByTopic: (topicId) => api.get(`/topics/${topicId}/poll`),
  submitVote: (data) => api.post('/polls/vote', data),
}

export const commentApi = {
  deleteComment: (id) => api.delete(`/comments/${id}`),
}

export const topicApi = {
  deleteTopic: (id) => api.delete(`/topics/${id}`),
  pinTopic: (id, pinned) => api.put(`/topics/${id}/pin`, { pinned }),
}

export const commentPinApi = {
  pinComment: (topicId, commentId, pinned) => api.put(`/topics/${topicId}/comments/${commentId}/pin`, { pinned }),
  bestComment: (topicId, commentId, best) => api.put(`/topics/${topicId}/comments/${commentId}/best`, { best }),
}

export const reportApi = {
  createReport: (data) => api.post('/reports', data),
  getMyReports: () => api.get('/user/reports'),
}

export default api;
