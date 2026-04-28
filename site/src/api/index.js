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

export const draftApi = {
  getDrafts: () => api.get('/drafts'),
  getDraft: (id) => api.get(`/drafts/${id}`),
  createDraft: (data) => api.post('/drafts', data),
  updateDraft: (id, data) => api.put(`/drafts/${id}`, data),
  deleteDraft: (id) => api.delete(`/drafts/${id}`),
  getConfig: () => api.get('/drafts/config'),
  getVersions: (draftId) => api.get(`/drafts/${draftId}/versions`),
  createVersion: (draftId, data) => api.post(`/drafts/${draftId}/versions`, data),
  getVersion: (draftId, versionId) => api.get(`/drafts/${draftId}/versions/${versionId}`),
  previewVersion: (draftId, versionId) => api.get(`/drafts/${draftId}/versions/${versionId}/preview`),
  restoreVersion: (draftId, versionId) => api.post(`/drafts/${draftId}/versions/${versionId}/restore`),
  deleteVersion: (draftId, versionId) => api.delete(`/drafts/${draftId}/versions/${versionId}`),
}

export default api;
