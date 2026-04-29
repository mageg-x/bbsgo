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
  withCredentials: true, // 自动发送和接收 Cookie
});

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      return Promise.reject(new ApiError(res.code, res.message));
    }
  },
  (error) => {
    // 401 错误时清除本地用户状态
    if (error.response?.status === 401) {
      localStorage.removeItem("user");
    }

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
  unlockPrivate: (id) => api.post(`/topics/${id}/unlock-private`),
  extendPrivate: (id, duration, unit) => api.post(`/topics/${id}/extend-private`, { duration, unit }),
}

export const commentPinApi = {
  pinComment: (topicId, commentId, pinned) => api.put(`/topics/${topicId}/comments/${commentId}/pin`, { pinned }),
  bestComment: (topicId, commentId, best) => api.put(`/topics/${topicId}/comments/${commentId}/best`, { best }),
}

export const reportApi = {
  createReport: (data) => api.post('/reports', data),
  getMyReports: () => api.get('/user/reports'),
}

export const homeApi = {
  getHomePage: (params) => api.get('/homepage', { params }),
  getHomePageWithQuery: (params) => api.get('/homepage/query', { params }),
}

export const authApi = {
  logout: () => api.post('/logout'),
}

export default api;
